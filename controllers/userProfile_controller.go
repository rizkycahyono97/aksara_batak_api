package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
)

type UserProfileController struct {
	UserProfileService services.UserProfileService
	Log                *slog.Logger
}

func NewUserProfileController(userProfileService services.UserProfileService, log *slog.Logger) *UserProfileController {
	return &UserProfileController{
		UserProfileService: userProfileService,
		Log:                log,
	}
}

// handler untuk mengambil data pengguna yang login
func (c *UserProfileController) GetMyProfile(f *fiber.Ctx) error {
	//ambil token jwt dari .Locals
	userToken := f.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uuid"].(string)
	c.Log.InfoContext(f.Context(), "get user profile process started", "userID", userID)

	//service
	profile, err := c.UserProfileService.FindUserProfileByID(f.UserContext(), userID)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "user profile not found" {
			return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
				Code:    "200",
				Message: "user profile not found",
				Data:    nil,
			})
		}

		c.Log.ErrorContext(f.Context(), "internal server error on get profile", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "INTERNAL SERVER ERROR",
			Data:    nil,
		})
	}

	c.Log.InfoContext(f.Context(), "successfully retrieved user profile", "userID", userID)
	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "successfully retrieved user profile",
		Data:    profile,
	})
}

func (c *UserProfileController) UpdateMyProfile(f *fiber.Ctx) error {
	//ambil token dari .Locals
	userToken := f.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uuid"].(string)
	c.Log.InfoContext(f.Context(), "update user profile process started", "userID", userID)

	//parsing req ke DTO
	var request web.UserProfileUpdateRequest
	if err := f.BodyParser(&request); err != nil {
		c.Log.ErrorContext(f.Context(), "failed to parse request body for update profile", "error", err)
		return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
			Code:    "400",
			Message: "failed to parse request body",
			Data:    nil,
		})
	}

	//service
	updatedProfile, err := c.UserProfileService.UpdateUserProfile(f.UserContext(), userID, request)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.Log.WarnContext(f.Context(), "update profile request validation failed", "error", err.Error())
			return f.Status(fiber.StatusBadRequest).JSON(web.ApiResponse{
				Code:    "400",
				Message: "update profile request validation failed",
				Data:    nil,
			})
		}

		if err.Error() == "user not found" {
			c.Log.WarnContext(f.Context(), "update attempt for non-existent user", "userID", userID)
			return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
				Code:    "200",
				Message: "user profile not found",
				Data:    nil,
			})
		}

		c.Log.ErrorContext(f.Context(), "internal server error on update profile", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "internal server error on update profile",
			Data:    nil,
		})
	}

	c.Log.InfoContext(f.Context(), "user profile updated successfully", "userID", userID)
	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "User Updated Succesfully",
		Data:    updatedProfile,
	})
}
