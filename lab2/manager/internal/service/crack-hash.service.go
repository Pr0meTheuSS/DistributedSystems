package service

import (
	"context"
	"manager/internal/model"

	"go.uber.org/zap"
)

type CrackHashServiceInterface interface {
	CrackHash(context.Context, model.CrackHashRequest) (*model.CrackHashRecord, error)
	GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error)
	GetCrackHashResult(context context.Context, requestID string) (*model.CrackHashResult, error)
}

type CrackHashService struct {
	logger *zap.Logger
}

func NewCrackHashService(logger *zap.Logger) CrackHashServiceInterface {
	return &CrackHashService{
		logger: logger,
	}
}

func (s *CrackHashService) CrackHash(context.Context, model.CrackHashRequest) (*model.CrackHashRecord, error) {
	return &model.CrackHashRecord{}, nil
}

func (s *CrackHashService) GetCrackHashProgress(context context.Context, requestID string) (*model.CrackHashProgress, error) {
	return &model.CrackHashProgress{}, nil
}

func (s *CrackHashService) GetCrackHashResult(context context.Context, requestID string) (*model.CrackHashResult, error) {
	return &model.CrackHashResult{}, nil
}
