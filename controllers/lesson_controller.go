package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type LessonController struct {
	LessonService services.LessonService
	Log           *slog.Logger
}

func NewLessonController(lessonService services.LessonService, log *slog.Logger) *LessonController {
	return &LessonController{
		Log:           log,
		LessonService: lessonService,
	}
}

func (c *LessonController) GetAllLessons(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "get all lessons process started")

	//service
	lessons, err := c.LessonService.GetAllLessons(f.Context())
	if err != nil {
		c.Log.ErrorContext(f.Context(), "internal server error on get all lessons", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "internal server error on get all lessons",
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "success",
		Data:    lessons,
	})
}
