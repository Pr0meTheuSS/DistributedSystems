package handler

import (
	"encoding/json"
	"manager/internal/dto"
	"manager/internal/service"
	"net/http"
)

type CrackHashHandler struct {
	service service.CrackHashServiceInterface
}

func NewCrackHashHandler(service service.CrackHashServiceInterface) *CrackHashHandler {
	return &CrackHashHandler{
		service: service,
	}
}

// PostCrackHash godoc
// @Summary     Отправить задачу на перебор хэша
// @Description Принимает хэш и максимальную длину пароля, возвращает ID задачи и статус
// @Tags        crack-hash
// @Accept      json
// @Produce     json
// @Param       request body dto.CrackHashRequestDto true "Параметры задачи"
// @Success     202 {object} dto.CrackHashResponseDto
// @Failure     400 {string} string "Невалидный запрос"
// @Failure     500 {string} string "Ошибка на сервере"
// @Router      /api/hash-cracks [post]
func (h *CrackHashHandler) PostCrackHash(writer http.ResponseWriter, request *http.Request) {
	var crackHashRequestDto dto.CrackHashRequestDto

	if err := json.NewDecoder(request.Body).Decode(&crackHashRequestDto); err != nil {
		http.Error(writer, buildErrorMessage(err, http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	crackHashRecord, err := h.service.CrackHash(request.Context(), dto.MapCrackHashRequestToModel(crackHashRequestDto))
	if err != nil {
		http.Error(writer, buildErrorMessage(err, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	json.NewEncoder(writer).Encode(dto.MapCrackHashRecordToDto(*crackHashRecord))
}
