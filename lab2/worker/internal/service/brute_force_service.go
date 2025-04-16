package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"worker/internal/dto"

	"go.uber.org/zap"
)

type BruteForceService struct {
	logger       *zap.Logger
	progressData sync.Map // map[string]int64

}

func NewBruteForceService(logger *zap.Logger) *BruteForceService {
	return &BruteForceService{
		logger:       logger,
		progressData: sync.Map{},
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func (s *BruteForceService) Crack(ctx context.Context, subTask *dto.SubTask) (*dto.CrackHashResultDto, error) {
	s.logger.Info("Crack called with args", zap.Any("sub_task", subTask))

	totalWords := TotalWords(subTask.Alphabet, subTask.Length)

	if subTask.PartNumber >= subTask.PartCount {
		return nil, fmt.Errorf("invalid partNumber: must be less than partCount")
	}

	begin := (totalWords / subTask.PartCount) * subTask.PartNumber
	end := (totalWords / subTask.PartCount) * min(subTask.PartNumber+1, subTask.PartCount)

	if begin >= end {
		return nil, fmt.Errorf("invalid range: begin (%d) >= end (%d)", begin, end)
	}

	expectedHash := strings.ToLower(subTask.Hash)
	results := make([]string, 0)

	var progress int64
	total := end - begin

	for i := begin; i < end; i++ {
		word := IndexToWord(i, subTask.Alphabet, subTask.Length)
		hash := fmt.Sprintf("%x", md5.Sum([]byte(word)))

		if hash == expectedHash {
			results = append(results, word)
			s.logger.Info("Match found", zap.String("word", word))
		}

		progress++
		percent := float64(progress*100) / float64(total)
		s.progressData.Store(subTask.TaskID, percent)
	}

	defer s.progressData.Delete(subTask.TaskID)

	return &dto.CrackHashResultDto{
		ID:      subTask.TaskID,
		Answers: results,
	}, nil
}

func (s *BruteForceService) GetProgress(taskID string) float64 {
	val, ok := s.progressData.Load(taskID)
	if !ok {
		return -1
	}
	return val.(float64)
}

func TotalWords(alphabet string, maxLength int64) int64 {
	n := int64(len(alphabet))
	var total int64 = 0
	var power int64 = 1

	for l := int64(1); l <= maxLength; l++ {
		power *= n
		total += power
	}
	return total
}

func IndexToWord(index int64, alphabet string, maxLength int64) string {
	n := int64(len(alphabet))
	var lengths []int64
	var total int64 = 0
	power := int64(1)

	for l := int64(1); l <= maxLength; l++ {
		power *= n
		total += power
		lengths = append(lengths, total)
	}

	var length int64
	var prevTotal int64
	for i, t := range lengths {
		if index < t {
			length = int64(i) + 1
			if i > 0 {
				prevTotal = lengths[i-1]
			}
			break
		}
	}

	offset := index - prevTotal
	word := make([]byte, length)

	for i := length - 1; i >= 0; i-- {
		charIndex := offset % n
		word[i] = alphabet[charIndex]
		offset /= n
	}

	return string(word)
}
