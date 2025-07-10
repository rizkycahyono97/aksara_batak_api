package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

// dependencies
type AuthController struct {
	AuthService services.AuthService
	Log         *slog.Logger
}

// dependncy injection
func NewAuthController(authService services.AuthService, log *slog.Logger) *AuthController {
	return &AuthController{
		AuthService: authService,
		Log:         log,
	}
}

func (c *AuthController) Register(f *fiber.Ctx) error {
	c.Log.InfoContext(f.UserContext(), "Register new user started")

	//parsing request body
	var request web.RegisterUserRequest
	if err := f.BodyParser(&request); err != nil {
		c.Log.ErrorContext(f.UserContext(), "Failed to parse body")
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "BAD REQUEST",
			Data:    nil,
		})
	}

	//panggil service layer
	newUser, err := c.AuthService.Register(f.UserContext(), request)
	if err != nil {
		var validationError *validator.ValidationErrors
		if errors.As(err, &validationError) {
			c.Log.ErrorContext(f.UserContext(), "Validation Error")
			return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
				Code:    "400",
				Message: "BAD REQUEST",
				Data:    nil,
			})
		}

		if err.Error() == "email already exists" {
			c.Log.ErrorContext(f.UserContext(), "Email already exists")
			return f.Status(fiber.StatusConflict).JSON(web.ApiResponse{
				Code:    "409",
				Message: "Email already exists",
				Data:    nil,
			})
		}

		c.Log.ErrorContext(f.UserContext(), "Internal Server Error")
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL SERVER ERROR",
			Data:    nil,
		})
	}
	c.Log.InfoContext(f.UserContext(), "Register new user succeeded")

	//assing DTO response
	response := web.RegisterUserResponse{
		UUID:  newUser.UUID,
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}

	return f.Status(fiber.StatusCreated).JSON(web.ApiResponse{
		Code:    "201",
		Message: "User Created Successfully",
		Data:    response,
	})
}
