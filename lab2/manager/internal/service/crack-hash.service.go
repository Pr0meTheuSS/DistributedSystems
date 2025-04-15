package service

import (
	"context"
	"manager/internal/dto"
	"manager/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CrackHashServiceInterface interface {
	CrackHash(context.Context, model.CrackHashRequest) (*model.CrackHashRecord, error)
	GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error)
	GetCrackHashResult(context context.Context, requestID string) (*model.CrackHashResult, error)
}

type CrackHashService struct {
	logger              *zap.Logger
	subTaskQueueService SubTaskQueueServiceInterface
	workersAmount       int64
}

func NewCrackHashService(logger *zap.Logger, subTaskQueueservice SubTaskQueueServiceInterface, workersAmount int64) CrackHashServiceInterface {
	return &CrackHashService{
		logger:              logger,
		subTaskQueueService: subTaskQueueservice,
		workersAmount:       workersAmount,
	}
}

func (s *CrackHashService) CrackHash(ctx context.Context, request model.CrackHashRequest) (*model.CrackHashRecord, error) {
	for i := int64(0); i < s.workersAmount; i++ {
		subTask := dto.SubTask{
			TaskID:     uuid.NewString(),
			Hash:       request.Hash,
			Length:     request.Length,
			Alphabet:   request.Alphabet,
			PartNumber: i,
			PartCount:  s.workersAmount,
		}
		s.subTaskQueueService.SendSubTask(ctx, &subTask)
	}

	return &model.CrackHashRecord{}, nil
}

func (s *CrackHashService) GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error) {
	return &model.CrackHashProgress{}, nil
}

func (s *CrackHashService) GetCrackHashResult(context context.Context, requestID string) (*model.CrackHashResult, error) {
	return &model.CrackHashResult{}, nil
}
