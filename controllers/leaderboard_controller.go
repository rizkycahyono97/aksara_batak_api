package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
	"strconv"
)

type LeaderboardController struct {
	LeaderboardService services.LeaderboardService
	Log                *slog.Logger
}

func NewLeaderboardController(leaderboardService services.LeaderboardService, log *slog.Logger) *LeaderboardController {
	return &LeaderboardController{
		LeaderboardService: leaderboardService,
		Log:                log,
	}
}

func (c *LeaderboardController) GetLeaderboards(f *fiber.Ctx) error {
	c.Log.InfoContext(f.Context(), "get leaderboard process started")

	//query params
	limitStr := f.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	//service
	leaderboard, err := c.LeaderboardService.GetLeaderboards(f.Context(), limit)
	if err != nil {
		c.Log.ErrorContext(f.Context(), "internal server error on get leaderboard", "error", err)
		return f.Status(fiber.StatusInternalServerError).JSON(web.ApiResponse{
			Code:    "500",
			Message: "Internal Server Error",
			Data:    nil,
		})
	}

	return f.Status(fiber.StatusOK).JSON(web.ApiResponse{
		Code:    "200",
		Message: "OK",
		Data:    leaderboard,
	})
}
