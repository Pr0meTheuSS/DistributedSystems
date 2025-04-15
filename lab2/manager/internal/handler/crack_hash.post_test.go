package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"manager/internal/dto"
	"manager/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// мок-сервис
type mockCrackHashService struct {
	record   *model.CrackHashRecord
	progress *model.CrackHashProgress
	result   *model.CrackHashResult
	err      error
}

func (m *mockCrackHashService) CrackHash(ctx context.Context, req model.CrackHashRequest) (*model.CrackHashRecord, error) {
	return m.record, m.err
}

func (m *mockCrackHashService) GetCrackHashProgress(ctx context.Context, id string) (*model.CrackHashProgress, error) {
	return m.progress, m.err
}

func (m *mockCrackHashService) GetCrackHashResult(ctx context.Context, id string) (*model.CrackHashResult, error) {
	return m.result, m.err
}

func TestPostCrackHash_Success(t *testing.T) {
	mockResult := &model.CrackHashRecord{
		ID:     "abc-123",
		Status: "processing",
	}

	service := &mockCrackHashService{record: mockResult}
	handler := NewCrackHashHandler(service)

	requestBody := dto.CrackHashRequestDto{
		Hash:   "5d41402abc4b2a76b9719d911017c592",
		Length: 5,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/crack-hash", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.PostCrackHash(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)

	var response dto.CrackHashResponseDto
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "abc-123", response.ID)
	assert.Equal(t, "processing", response.Status)
}

func TestPostCrackHash_BadRequest(t *testing.T) {
	service := &mockCrackHashService{}
	handler := NewCrackHashHandler(service)

	// некорректный JSON
	req := httptest.NewRequest(http.MethodPost, "/crack-hash", bytes.NewReader([]byte("{not-json")))
	recorder := httptest.NewRecorder()

	handler.PostCrackHash(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPostCrackHash_InternalError(t *testing.T) {
	service := &mockCrackHashService{err: errors.New("some error")}
	handler := NewCrackHashHandler(service)

	requestBody := dto.CrackHashRequestDto{
		Hash:   "5d41402abc4b2a76b9719d911017c592",
		Length: 5,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/crack-hash", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.PostCrackHash(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
