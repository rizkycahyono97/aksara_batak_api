package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type ChatHistoryRepository interface {
	Create(ctx context.Context, chat *domain.ChatHistories) error
	GetLastFiveByUserID(ctx context.Context, userID string) ([]domain.ChatHistories, error)
	DeleteByUserID(ctx context.Context, userID string) error
	CountByUserID(ctx context.Context, userID string) (int, error)
	DeleteExcess(ctx context.Context, userID string, limit int) error
}
