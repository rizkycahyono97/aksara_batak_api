package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type ContactUsController struct {
	ContactUsService services.ContactUsService
	Log              *slog.Logger
	Validate         *validator.Validate
}

func NewContactUsController(
	contactUsService services.ContactUsService,
	log *slog.Logger,
	validate *validator.Validate,
) *ContactUsController {
	return &ContactUsController{
		ContactUsService: contactUsService,
		Log:              log,
		Validate:         validate,
	}
}

func (c *ContactUsController) SubmitContact(f *fiber.Ctx) error {
	var request web.ContactUsRequest
	if err := f.BodyParser(&request); err != nil {
		c.Log.Warn("Gagal mem-parsing body request 'Contact Us'", "error", err)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD_REQUEST",
			Data:    "format tidak valid",
		})
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warn("Request 'Contact Us' gagal validasi", "error", err.Error(), "request_data", request)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD_REQUEST",
			Data:    "format tidak valid",
		})
	}

	if err := c.ContactUsService.ProcessSubmission(f.Context(), request); err != nil {
		c.Log.Error("Service 'Contact Us' mengembalikan error", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL_SERVER_ERROR",
			Data:    err.Error(),
		})
	}

	return f.Status(fiber.StatusAccepted).JSON(web.ApiResponse{
		Code:    "202",
		Message: "ACCEPTED",
		Data:    "Pesan Anda telah diterima dan akan segera kami proses. Terima kasih!",
	})
}
