package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
	"strconv"
)

type LessonController struct {
	LessonService services.LessonService
	QuizService   services.QuizService
	Log           *slog.Logger
}

func NewLessonController(lessonService services.LessonService, quizService services.QuizService, log *slog.Logger) *LessonController {
	return &LessonController{
		Log:           log,
		LessonService: lessonService,
		QuizService:   quizService,
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

// Get quizzes By lesson ID
func (c *LessonController) GetQuizzesByLessonID(f *fiber.Ctx) error {
	lessonIDStr := f.Params("lessonID")
	c.Log.InfoContext(f.Context(), "get quizzes process started", "lessonID", lessonIDStr)

	//parsing and validasi
	id, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		c.Log.InfoContext(f.Context(), "invalid lesson ID parameter", "lessonID", lessonIDStr)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "invalid lesson ID parameter",
		})
	}
	lessonID := uint(id)

	//service
	quizzes, err := c.QuizService.GetQuizzesByLessonID(f.Context(), lessonID)
	if err != nil {
		c.Log.ErrorContext(f.Context(), "internal server error on get quizzes by lesson", "error", err, "lessonID", lessonID)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "internal server error on get quizzes by lesson",
		})
	}

	c.Log.InfoContext(f.Context(), "successfully retrieved quizzes by lesson", "lessonID", lessonID, "count", len(quizzes))
	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "Quizzes for the lesson retrieved successfully",
		Data:    quizzes,
	})
}
