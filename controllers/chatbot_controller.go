package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type ChatbotController struct {
	ChatbotService services.ChatbotService
	Log            *slog.Logger
}

func NewChatbotController(chatbotService services.ChatbotService, log *slog.Logger) *ChatbotController {
	return &ChatbotController{
		ChatbotService: chatbotService,
		Log:            log,
	}
}

func (c *ChatbotController) HandlePublicChat(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "HandlePublicChat started.......")

	//parse request ke json
	var request web.ChatbotRequest
	if err := f.BodyParser(&request); err != nil {
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD_REQUEST",
			Data:    "request body parse error",
		})
	}

	// jika pesan kosong
	if request.Message == "" {
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_SERVER_ERROR",
			Data:    "pesan tidak boleh kosong",
		})
	}

	//service
	response, err := c.ChatbotService.GeneratePublicResponse(f.Context(), request)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_SERVER_ERROR",
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "OK",
		Data:    response,
	})
}
