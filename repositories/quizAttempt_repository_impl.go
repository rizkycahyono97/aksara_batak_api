package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type QuizAttemptRepositoryImpl struct {
	DB *gorm.DB
}

func NewQuizAttemptRepository(db *gorm.DB) QuizAttemptRepository {
	return QuizAttemptRepositoryImpl{DB: db}
}

func (q QuizAttemptRepositoryImpl) FindAllQuizAttemptByUserID(ctx context.Context, userID string) ([]domain.QuizAttempts, error) {
	var attempts []domain.QuizAttempts

	//preload di table Quizzes
	err := q.DB.WithContext(ctx).
		Preload("Quizzes").
		Where("user_id = ?", userID).
		Find(&attempts).Error
	if err != nil {
		return []domain.QuizAttempts{}, err
	}
	return attempts, nil
}
