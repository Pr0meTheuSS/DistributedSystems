package handler

import (
	"encoding/json"
	"errors"
	"manager/internal/dto"
	"manager/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPingService struct {
	response model.Pong
	err      error
}

func (m *mockPingService) GetPing() (model.Pong, error) {
	return m.response, m.err
}

func TestPingHandler_GetPing_Success(t *testing.T) {
	mockService := &mockPingService{
		response: model.Pong{Message: "Pong"},
		err:      nil,
	}
	handler := NewPingHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	handler.GetPing(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var body dto.PingResponse
	json.NewDecoder(res.Body).Decode(&body)

	assert.Equal(t, "Pong", body.Pong)
}

func TestPingHandler_GetPing_Error(t *testing.T) {
	mockService := &mockPingService{
		response: model.Pong{},
		err:      errors.New("some error"),
	}
	handler := NewPingHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	handler.GetPing(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
