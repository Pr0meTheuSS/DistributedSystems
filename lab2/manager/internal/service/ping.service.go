package service

import (
	"manager/internal/model"

	"go.uber.org/zap"
)

type PingServiceInterface interface {
	GetPing() (model.Pong, error)
}

type PingService struct {
	logger *zap.Logger
}

func NewPingService(logger *zap.Logger) PingServiceInterface {
	return &PingService{
		logger: logger,
	}
}

func (s *PingService) GetPing() (model.Pong, error) {
	s.logger.Info("GetPing method called.")
	return model.Pong{
		Message: "Pong",
	}, nil
}
