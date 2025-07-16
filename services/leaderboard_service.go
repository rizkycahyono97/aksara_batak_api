package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type LeaderboardService interface {
	GetLeaderboards(ctx context.Context, limit int) ([]web.LeaderboardResponse, error)
}
