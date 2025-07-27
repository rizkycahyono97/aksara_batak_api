package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type TranslateController struct {
	TranslateService services.TranslateService
	Log              *slog.Logger
}

func NewTranslateController(translateService services.TranslateService, log *slog.Logger) *TranslateController {
	return &TranslateController{
		TranslateService: translateService,
		Log:              log,
	}
}

func (c *TranslateController) Translate(f *fiber.Ctx) error {
	var req web.TranslateRequest

	if err := f.BodyParser(&req); err != nil {
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "Invalid request body",
		})
	}

	if req.Text == "" {
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "Text cannot be empty",
		})
	}

	if !req.Direction.IsValid() {
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "Invalid Direction",
		})
	}

	response, err := c.TranslateService.Translate(f.Context(), req)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_sERVER_ERROR",
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "Success",
		Data:    response,
	})
}
