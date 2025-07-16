package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"log/slog"
)

type LeaderboardServiceImpl struct {
	UserProfileRepo repositories.UserProfileRepository
	Log             *slog.Logger
}

func NewLeaderboardService(userProfileRepo repositories.UserProfileRepository, log *slog.Logger) LeaderboardService {
	return &LeaderboardServiceImpl{
		UserProfileRepo: userProfileRepo,
		Log:             log,
	}
}

func (s LeaderboardServiceImpl) GetLeaderboards(ctx context.Context, limit int) ([]web.LeaderboardResponse, error) {
	s.Log.InfoContext(ctx, "get leaderboard process started", "limit", limit)

	//repository
	topUsers, err := s.UserProfileRepo.GetTopUsers(ctx, limit)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to get top users from repository", "error", err)
		return nil, err
	}

	//logika transformasi data dan penambahan peringkat
	var leaderboardsResponse []web.LeaderboardResponse
	for i, user := range topUsers {
		response := web.LeaderboardResponse{
			Rank:      i + 1,
			UserID:    user.UUID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			TotalXP:   int(user.TotalXP),
		}
		leaderboardsResponse = append(leaderboardsResponse, response)
	}

	s.Log.InfoContext(ctx, "successfully retrieved leaderboard data", "count", len(leaderboardsResponse))
	return leaderboardsResponse, nil
}
