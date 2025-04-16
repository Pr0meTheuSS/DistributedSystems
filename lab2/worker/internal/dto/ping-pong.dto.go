package dto

import "worker/internal/model"

type PingResponse struct {
	Pong string `json:"pong"`
}

func MapPongToPingResponse(pong model.Pong) PingResponse {
	return PingResponse{
		Pong: pong.Message,
	}
}
