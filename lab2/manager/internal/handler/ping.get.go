package handler

import (
	"encoding/json"
	"manager/internal/dto"
	"manager/internal/service"
	"net/http"
)

type PingHandler struct {
	service service.PingServiceInterface
}

func NewPingHandler(service service.PingServiceInterface) PingHandler {
	return PingHandler{
		service: service,
	}
}

func (p *PingHandler) GetPing(writer http.ResponseWriter, request *http.Request) {
	pong, err := p.service.GetPing()
	if err != nil {
		http.Error(writer, buildErrorMessage(err, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dto.MapPongToPingResponse(pong))
}
