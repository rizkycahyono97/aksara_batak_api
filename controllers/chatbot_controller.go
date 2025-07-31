package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type ChatbotController struct {
	ChatbotService services.ChatbotService
	Log            *slog.Logger
	Validate       *validator.Validate
}

func NewChatbotController(chatbotService services.ChatbotService, log *slog.Logger, validate *validator.Validate) *ChatbotController {
	return &ChatbotController{
		ChatbotService: chatbotService,
		Log:            log,
		Validate:       validate,
	}
}

func (c *ChatbotController) HandlePublicChat(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "HandlePublicChat started.......")

	//parse request ke json
	var request web.ChatbotRequest
	if err := f.BodyParser(&request); err != nil {
		c.Log.InfoContext(f.Context(), "request body parse error", "request", request.Message)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD_REQUEST",
			Data:    "request body parse error",
		})
	}

	// jika pesan kosong
	if request.Message == "" {
		c.Log.InfoContext(f.Context(), "pesan kosong")
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

func (c *ChatbotController) HandlePrivateChat(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "HandlePrivateChat started.......")

	//bind request
	var req web.ChatbotRequest
	if err := f.BodyParser(&req); err != nil {
		c.Log.InfoContext(f.Context(), "request body parse error", "request", req.Message)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD_REQUEST",
		})
	}

	//ambil userID
	userToken := f.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uuid"].(string)
	c.Log.InfoContext(f.Context(), "user identified for personalized quiz list", "userID", userID)

	req.Userid = userID
	response, err := c.ChatbotService.GeneratePrivateResponse(f.Context(), req)
	if err != nil {
		c.Log.Error("Failed to process chat", slog.String("user_id", req.Userid), slog.Any("error", err))
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "OK",
		Data:    response,
	})
}

func (c *ChatbotController) GetChatPrivateHistory(f *fiber.Ctx) error {
	userTokenClaim := f.Locals("user")
	if userTokenClaim == nil {
		c.Log.Warn("Gagal mengambil history: middleware tidak menyediakan data token")
		return f.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Code:    "401",
			Message: "UNAUTHORIZED",
		})
	}

	userToken := userTokenClaim.(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uuid"].(string)
	c.Log.Info("Mengambil riwayat chat untuk pengguna", "userID", userID)

	history, err := c.ChatbotService.GetChatPrivateHistory(f.Context(), userID)
	if err != nil {
		c.Log.Error("Service gagal mengambil riwayat chat", "userID", userID, "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_SERVER_ERROR",
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "OK",
		Data:    history,
	})
}
