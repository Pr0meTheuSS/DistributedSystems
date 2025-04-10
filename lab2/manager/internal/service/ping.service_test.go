package service

import (
	"manager/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestPingService_GetPing(t *testing.T) {
	logger := zaptest.NewLogger(t)
	service := NewPingService(logger)

	pong, err := service.GetPing()

	assert.NoError(t, err)
	assert.Equal(t, model.Pong{Message: "Pong"}, pong)
}
