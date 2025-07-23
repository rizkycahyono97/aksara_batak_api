package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type QuizAttemptRepository interface {
	FindAllQuizAttemptByUserID(ctx context.Context, userID string) ([]domain.QuizAttempts, error)
	FindCompletedQuizIDsByUserID(ctx context.Context, userID string) ([]uint, error)
}
