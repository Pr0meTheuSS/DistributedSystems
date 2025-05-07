package handler

import (
	"context"
	"encoding/json"
	"errors"
	"manager/internal/dto"
	"manager/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)


func TestGetCrackHashProgress_Success(t *testing.T) {
	duration := 15 * time.Second
	progress := &model.CrackHashProgress{
		RecordID:                "progress-id-123",
		ProgressInPercents:      42.5,
		CurrentHandlingDuration: duration,
	}

	service := &mockCrackHashService{progress: progress}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks/progress-id-123/progress", nil)

	// Устанавливаем параметр {id} в контексте маршрута
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "progress-id-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	recorder := httptest.NewRecorder()
	handler.GetCrackHashProgress(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)

	var response dto.CrackHashProgressDto
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "progress-id-123", response.ID)
	assert.Equal(t, 42.5, response.ProgressInPercents)
	assert.Equal(t, duration.String(), response.CurrentHandlingDuration)
}

func TestGetCrackHashProgress_ServiceError(t *testing.T) {
	service := &mockCrackHashService{err: errors.New("service failure")}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks/fail/progress", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "fail")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	recorder := httptest.NewRecorder()
	handler.GetCrackHashProgress(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetCrackHashProgress_MissingID(t *testing.T) {
	service := &mockCrackHashService{}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks//progress", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext())) // Без id

	recorder := httptest.NewRecorder()
	handler.GetCrackHashProgress(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
