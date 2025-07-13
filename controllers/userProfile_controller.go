package controllers

import (
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
