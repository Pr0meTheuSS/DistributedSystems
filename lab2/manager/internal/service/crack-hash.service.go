package service

import (
	"context"
	"fmt"
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
	UpdateProgressOrResult(ctx context.Context, response dto.WorkerResponseDto) error
}

type CrackHashService struct {
	logger              *zap.Logger
	repository          *repository.RequestsMongoRepository
	subTaskQueueService SubTaskQueueServiceInterface
	workersAmount       int64
}

func (s *CrackHashService) UpdateProgressOrResult(ctx context.Context, response dto.WorkerResponseDto) error {
	err := s.repository.UpdateSubTaskProgress(ctx, response.ID, response.Progress)
	if err != nil {
		s.logger.Error("Failed to update subtask progress", zap.Error(err))
		return err
	}

	s.logger.Info("Updated progress for subtask", zap.Any("subtask", response))

	// if response.IsFinished {
	s.logger.Info("Subtask finished, saving answers...", zap.String("subtask_id", response.ID))

	if err := s.repository.MarkSubTaskAsFinished(ctx, response.ID, response.Answers); err != nil {
		s.logger.Error("Failed to mark subtask as finished", zap.Error(err))
		return err
	}

	// Получаем TaskID, чтобы проверить статус всех сабтасок
	subTask, err := s.repository.GetSubTaskByID(ctx, response.ID)
	if err != nil {
		s.logger.Error("Failed to fetch subtask for completion check", zap.Error(err))
		return err
	}
	if subTask == nil {
		s.logger.Warn("Subtask not found by ID", zap.String("subtask_id", response.ID))
		return nil
	}

	allFinished, err := s.repository.AreAllSubTasksFinished(ctx, subTask.TaskID)
	if err != nil {
		s.logger.Error("Failed to check subtask completion status", zap.Error(err))
		return err
	}

	if allFinished {
		s.logger.Info("All subtasks complete. Marking main task as finished...", zap.String("task_id", subTask.TaskID))
		answers, err := s.repository.CollectAllAnswersByTaskID(ctx, subTask.TaskID)
		fmt.Println("ANSWERS: ", answers)

		if err != nil {
			s.logger.Error("Failed to collect all answers", zap.Error(err))
			return err
		}

		if err := s.repository.MarkMainTaskAsFinished(ctx, subTask.TaskID, answers); err != nil {
			s.logger.Error("Failed to MARK main task as finished", zap.Error(err))
			return err
		}

		// mainRecord, err := s.repository.GetByID(ctx, subTask.TaskID)
		// if err != nil {
		// 	s.logger.Error("Failed to load main task for final update", zap.Error(err))
		// 	return err
		// }
		// if mainRecord == nil {
		// 	s.logger.Warn("Main task not found", zap.String("task_id", subTask.TaskID))
		// 	return nil
		// }

		// mainRecord.Answers = answers
		// mainRecord.Status = model.READY
		// now := time.Now().UTC()
		// mainRecord.FinishedAt = &now

		// if _, err := s.repository.Update(ctx, mainRecord); err != nil {
		// 	s.logger.Error("Failed to update main task with final data", zap.Error(err))
		// 	return err
		// }
	}
	// }

	return nil
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
		subTask := model.SubTask{
			ID:         uuid.NewString(),
			TaskID:     record.ID,
			Hash:       request.Hash,
			Length:     request.Length,
			Alphabet:   request.Alphabet,
			PartNumber: i,
			PartCount:  s.workersAmount,
		}

		_, err := s.repository.SaveSubTask(ctx, &subTask)
		if err != nil {
			s.logger.Error("Error on saving sub task", zap.Error(err))
			return nil, err
		}

		if err := s.subTaskQueueService.SendSubTask(ctx, &subTask); err != nil {
			s.logger.Error("Error on sending sub task", zap.Error(err))
			return nil, err
		}

		s.logger.Info("Sub task sent to message queue", zap.Any("SubTask", subTask))
	}

	return record, nil
}

func (s *CrackHashService) GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error) {
	return &model.CrackHashProgress{}, nil
}

func (s *CrackHashService) GetCrackHashResult(ctx context.Context, requestID string) (*model.CrackHashResult, error) {
	s.logger.Info("call service.GetCrackHashResult()", zap.String("requestID", requestID))
	record, err := s.repository.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	return &model.CrackHashResult{
		RecordID: record.ID,
		Answers:  record.Answers,
	}, nil
}
