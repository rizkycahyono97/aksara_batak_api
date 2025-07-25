package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type QuizAttemptRepository interface {
	FindAllQuizAttemptByUserID(ctx context.Context, userID string) ([]domain.QuizAttempts, error)
	FindCompletedQuizIDsByUserID(ctx context.Context, userID string) ([]uint, error)
	CountByUserIDAndQuizID(ctx context.Context, userID string, quizID uint) (int64, error)
	FindHighestScoresByUserID(ctx context.Context, userID string) (map[uint]uint, error)
	HasUserPassedQuizBefore(ctx context.Context, userID string, quizID uint, passingScore uint) (bool, error)
}
