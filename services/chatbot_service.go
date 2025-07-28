package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type ChatbotService interface {
	GeneratePublicResponse(ctx context.Context, request web.ChatbotRequest) (web.ChatbotResponse, error)
	GeneratePrivateResponse(ctx context.Context, request web.ChatbotRequest) (web.ChatbotResponse, error)
}
