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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetCrackHashResult_Success(t *testing.T) {
	mockResult := &model.CrackHashResult{
		RecordID: "result-id-123",
		Answers:  []string{"hello", "world"},
	}

	service := &mockCrackHashService{result: mockResult}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks/result-id-123/result", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "result-id-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	recorder := httptest.NewRecorder()
	handler.GetCrackHashResult(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)

	var response dto.CrackHashResultDto
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "result-id-123", response.ID)
	assert.Equal(t, []string{"hello", "world"}, response.Answers)
}

func TestGetCrackHashResult_ServiceError(t *testing.T) {
	service := &mockCrackHashService{err: errors.New("some internal error")}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks/fail/result", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "fail")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	recorder := httptest.NewRecorder()
	handler.GetCrackHashResult(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetCrackHashResult_MissingID(t *testing.T) {
	service := &mockCrackHashService{}
	handler := NewCrackHashHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/hash-cracks//result", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext())) // без ID

	recorder := httptest.NewRecorder()
	handler.GetCrackHashResult(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
