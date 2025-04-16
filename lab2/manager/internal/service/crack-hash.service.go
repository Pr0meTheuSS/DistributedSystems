package service

import (
	"context"
	"manager/internal/dto"
	"manager/internal/model"
	"manager/internal/repository"
	"time"

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
	repository          *repository.RequestsMongoRepository
	subTaskQueueService SubTaskQueueServiceInterface
	workersAmount       int64
}

func NewCrackHashService(logger *zap.Logger, subTaskQueueservice SubTaskQueueServiceInterface, repository *repository.RequestsMongoRepository, workersAmount int64) CrackHashServiceInterface {
	return &CrackHashService{
		logger:              logger,
		subTaskQueueService: subTaskQueueservice,
		workersAmount:       workersAmount,
		repository:          repository,
	}
}

func (s *CrackHashService) CrackHash(ctx context.Context, request model.CrackHashRequest) (*model.CrackHashRecord, error) {
	record := &model.CrackHashRecord{
		ID:         uuid.NewString(),
		Hash:       request.Hash,
		Alphabet:   request.Alphabet,
		Length:     request.Length,
		Status:     model.PENDING,
		CreatedAt:  time.Now(),
		FinishedAt: nil,
	}
	s.logger.Info("Prepare to save crack hash record: ", zap.Any("Record", record))

	record, err := s.repository.Save(ctx, record)
	if err != nil {
		s.logger.Error("Error on saving record", zap.Error(err))
		return nil, err
	}

	for i := int64(0); i < s.workersAmount; i++ {
		subTask := dto.SubTask{
			TaskID:     record.ID,
			Hash:       request.Hash,
			Length:     request.Length,
			Alphabet:   request.Alphabet,
			PartNumber: i,
			PartCount:  s.workersAmount,
		}

		s.subTaskQueueService.SendSubTask(ctx, &subTask)
		s.logger.Info("Sub task sent to message queue", zap.Any("SubTask", subTask))
	}

	return record, nil
}

func (s *CrackHashService) GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error) {
	return &model.CrackHashProgress{}, nil
}

func (s *CrackHashService) GetCrackHashResult(context context.Context, requestID string) (*model.CrackHashResult, error) {
	return &model.CrackHashResult{}, nil
}
