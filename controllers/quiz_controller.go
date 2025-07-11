package controllers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type QuizControler struct {
	QuizService services.QuizService
	Log         *slog.Logger
}

func NewQuizController(quizService services.QuizService, log *slog.Logger) QuizControler {
	return QuizControler{
		QuizService: quizService,
		Log:         log,
	}
}

func (c *QuizControler) GetAllQuizzes(f *fiber.Ctx) error {
	//parsing query params ke DTO
	var filters web.FilterQuizRequest
	if err := f.BodyParser(&filters); err != nil {
		c.Log.ErrorContext(f.UserContext(), "failed to parse query params", "error", err)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "failed to parse query params",
			Data:    nil,
		})
	}

	//service
	quizzess, err := c.QuizService.GetAllQuizzes(f.UserContext(), filters)
	if err != nil {
		//jika error dari validation
		var validationError *validator.ValidationErrors
		if errors.As(err, &validationError) {
			c.Log.ErrorContext(f.UserContext(), "Validation Error")
			return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
				Code:    "400",
				Message: "BAD REQUEST: Validation Failed",
				Data:    nil,
			})
		}

		c.Log.ErrorContext(f.UserContext(), "Internal Server Error", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL SERVER ERROR",
			Data:    nil,
		})
	}

	if len(quizzess) == 0 {
		c.Log.ErrorContext(f.UserContext(), fmt.Sprintf("quizzess %s is empty", quizzess))
		return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
			Code:    "200",
			Message: "Quizzes Not Found",
			Data:    nil,
		})
	}

	c.Log.InfoContext(f.UserContext(), "get quizzes succeeded")
	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "STATUS_OK",
		Data:    quizzess,
	})
}
