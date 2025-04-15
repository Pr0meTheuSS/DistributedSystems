package handler

import (
	"encoding/json"
	"errors"
	"manager/internal/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetCrackHashProgress godoc
// @Summary     Получить текущий прогресс задачи.
// @Description Принимает ID задачи, возвращает текущий прогресс выполнения задачи.
// @Tags        crack-hash
// @Accept      json
// @Produce     json
// @Param       id query string true "Уникальный идентификатор задачи (ID)"
// @Success     200 {object} dto.CrackHashProgressDto
// @Failure     400 {string} string "Невалидный запрос"
// @Failure     404 {string} string "Задача по данному идентификатору не найдена"
// @Failure     500 {string} string "Ошибка на сервере"
// @Router      /api/hash-cracks/{id}/progress [get]
func (h *CrackHashHandler) GetCrackHashProgress(writer http.ResponseWriter, request *http.Request) {
	requestID := chi.URLParam(request, "id")
	if requestID == "" {
		http.Error(writer, buildErrorMessage(errors.New("unknown request ID"), http.StatusNotFound), http.StatusNotFound)
		return
	}

	crackHashProgress, err := h.service.GetCrackHashProgress(request.Context(), requestID)
	if err != nil {
		http.Error(writer, buildErrorMessage(err, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	json.NewEncoder(writer).Encode(dto.MapCrackHashProgressToDto(*crackHashProgress))
}
