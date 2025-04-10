package handler

import (
	"encoding/json"
	"fmt"
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

func buildInternalServerErrorString(err error) string {
	return fmt.Sprintf("Internal server error in service layout. Error message: %s", err.Error())
}

func (p *PingHandler) GetPing(writer http.ResponseWriter, request *http.Request) {
	pong, err := p.service.GetPing()
	if err != nil {
		http.Error(writer, buildInternalServerErrorString(err), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dto.MapPongToPingResponse(pong))
}
