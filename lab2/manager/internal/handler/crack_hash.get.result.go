package handler

import (
	"encoding/json"
	"errors"
	"manager/internal/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetCrackHashResult godoc
// @Summary     Получить результат задачи.
// @Description Принимает ID задачи, возвращает результат выполнения.
// @Tags        crack-hash
// @Accept      json
// @Produce     json
// @Param       id query string true "Уникальный идентификатор задачи (Id)"
// @Success     200 {object} dto.CrackHashResultDto
// @Failure     400 {string} string "Невалидный запрос"
// @Failure     404 {string} string "Задача по данному идентификатору не найдена"
// @Failure     500 {string} string "Ошибка на сервере"
// @Router      /api/hash-cracks/{id}/result [get]
func (h *CrackHashHandler) GetCrackHashResult(writer http.ResponseWriter, request *http.Request) {
	requestID := chi.URLParam(request, "id")
	if requestID == "" {
		http.Error(writer, buildErrorMessage(errors.New("unknown request ID"), http.StatusNotFound), http.StatusNotFound)
		return
	}

	crackHashResult, err := h.service.GetCrackHashResult(request.Context(), requestID)
	if err != nil {
		http.Error(writer, buildErrorMessage(err, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dto.MapCrackHashResultToDto(*crackHashResult))
}
