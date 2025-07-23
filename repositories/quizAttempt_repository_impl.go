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

// method untuk pengambilan quizzes id yang sudah selesai / is_completed
func (q QuizAttemptRepositoryImpl) FindCompletedQuizIDsByUserID(ctx context.Context, userID string) ([]uint, error) {
	var completedQuizIDs []uint

	// Distinct(): untuk memastikan setiap ID kuis hanya muncul sekali,
	// query untuk mencari user ini sudah menyelesaikan quiz apa aja,
	// jika ada quiz_id yang sama maka tampilkan sekali
	err := q.DB.WithContext(ctx).
		Distinct("quiz_id").
		Model(&domain.QuizAttempts{}).
		Where("user_id = ?", userID).
		Pluck("quiz_id", &completedQuizIDs).Error
	if err != nil {
		return nil, err
	}

	return completedQuizIDs, nil
}
