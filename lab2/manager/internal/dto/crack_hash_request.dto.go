package dto

import (
	"manager/internal/model"
	"time"
)

type CrackHashRequestDto struct {
	Hash     string `json:"hash"`
	Length   int64  `json:"length"`
	Alphabet string `json:"alphabet"`
}

type CrackHashResponseDto struct {
	ID         string     `json:"id"`
	Hash       string     `json:"hash"`
	Length     int64      `json:"length"`
	Alphabet   string     `json:"alphabet"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"create_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type CrackHashProgressDto struct {
	ID                      string  `json:"id"`
	ProgressInPercents      float64 `json:"progress_in_percents"`
	CurrentHandlingDuration string  `json:"current_handling_duration"`
	Hash                    string  `json:"hash"`
}

type CrackHashResultDto struct {
	ID      string   `json:"id"`
	Answers []string `json:"answers"`
}

func MapCrackHashRequestToModel(request CrackHashRequestDto) model.CrackHashRequest {
	return model.CrackHashRequest{
		Hash:     request.Hash,
		Alphabet: request.Alphabet,
		Length:   request.Length,
	}
}

func MapCrackHashRecordToDto(record model.CrackHashRecord) CrackHashResponseDto {
	return CrackHashResponseDto{
		ID:         record.ID,
		Hash:       record.Hash,
		Length:     record.Length,
		Alphabet:   record.Alphabet,
		Status:     string(record.Status),
		CreatedAt:  record.CreatedAt,
		FinishedAt: record.FinishedAt,
	}
}

func MapCrackHashProgressToDto(progress model.CrackHashProgress) CrackHashProgressDto {
	return CrackHashProgressDto{
		ID:                      progress.RecordID,
		ProgressInPercents:      progress.ProgressInPercents,
		CurrentHandlingDuration: progress.CurrentHandlingDuration.String(),
		Hash:                    progress.Hash,
	}
}

func MapCrackHashResultToDto(result model.CrackHashResult) CrackHashResultDto {
	return CrackHashResultDto{
		ID:      result.RecordID,
		Answers: result.Answers,
	}
}
