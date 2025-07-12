package controllers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
	"strconv"
)

type QuizController struct {
	QuizService services.QuizService
	Log         *slog.Logger
}

func NewQuizController(quizService services.QuizService, log *slog.Logger) *QuizController {
	return &QuizController{
		QuizService: quizService,
		Log:         log,
	}
}

func (c *QuizController) GetAllQuizzes(f *fiber.Ctx) error {
	//validasi query params

	allowedParams := map[string]bool{

		"dialect": true,

		"level": true,

		"title": true,
	}

	sentParams := f.Queries()

	for paramName := range sentParams {

		if !allowedParams[paramName] {

			c.Log.InfoContext(f.Context(), fmt.Sprintf("param %s not allowed", paramName))

			return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{

				Code: "400",

				Message: "param not allowed",

				Data: nil,
			})

		}

	}

	//validasi untuk field query paramsnya

	filters := web.FilterQuizRequest{}

	filters.Title = f.Query("title")

	filters.Dialect = f.Query("dialect")

	if level := f.Query("level"); level != "" {

		level, err := strconv.Atoi(level)

		if err == nil {

			levelUint := uint(level)

			filters.Level = &levelUint

		}

		c.Log.InfoContext(f.Context(), "level parse error")

	}

	//service

	quizzess, err := c.QuizService.GetAllQuizzes(f.UserContext(), filters)

	if err != nil {

		//jika error dari validation

		var validationError *validator.ValidationErrors

		if errors.As(err, &validationError) {

			c.Log.ErrorContext(f.UserContext(), "Validation Error")

			return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{

				Code: "400",

				Message: "BAD REQUEST",

				Data: nil,
			})

		}

		c.Log.ErrorContext(f.UserContext(), "Internal Server Error")

		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{

			Code: "500",

			Message: "INTERNAL SERVER ERROR",

			Data: nil,
		})

	}

	if len(quizzess) == 0 {

		c.Log.ErrorContext(f.UserContext(), fmt.Sprintf("quizzess %s is empty", quizzess))

		return f.Status(fiber.StatusOK).JSON(web.ApiResponse{

			Code: "200",

			Message: "not found",

			Data: nil,
		})

	}

	c.Log.InfoContext(f.UserContext(), "get quizzes succeeded")

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{

		Code: "200",

		Message: "STATUS_OK",

		Data: quizzess,
	})
}

func (c *QuizController) StartQuiz(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "start quiz")

	//mengambil query params quizID
	quizIDStr := f.Query("quizID")
	id, err := strconv.Atoi(quizIDStr)
	if err != nil {
		c.Log.InfoContext(f.Context(), "quizID parse error", "quizID", quizIDStr)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD REQUEST",
			Data:    nil,
		})
	}
	quizID := uint(id)

	//mengambil userID dari jwt
	userToken := f.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uuid"].(string)

	if userID == "" {
		c.Log.ErrorContext(f.Context(), "failed to get user ID from JWT claims")
		return f.Status(fiber.StatusUnauthorized).JSON(web.ApiResponse{
			Code:    "401",
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	//service layer
	response, err := c.QuizService.StartQuiz(f.Context(), quizID, userID)
	if err != nil {
		if err.Error() == "quiz has no questions" {
			c.Log.ErrorContext(f.Context(), "quiz has no questions", "quizID", quizID)
			return f.Status(fiber.StatusNotFound).JSON(web.ApiResponse{
				Code:    "404",
				Message: "quiz not found",
				Data:    nil,
			})
		}

		c.Log.ErrorContext(f.Context(), "Internal Server Error")
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL SERVER ERROR",
			Data:    nil,
		})
	}

	c.Log.InfoContext(f.Context(), "quiz started successfully", "sessionID", response.SessionID, "userID", userID)
	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "Quiz Started Succesfully",
		Data:    response,
	})
}
