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

// method untuk menghitung riwayat pengerjaan sebuah quiz
func (q QuizAttemptRepositoryImpl) CountByUserIDAndQuizID(ctx context.Context, userID string, quizID uint) (int64, error) {
	var count int64

	//Count = GORM untuk menjalankan query "SELECT COUNT(*)"
	// untuk menghitung jumlah quiz yang dikerjakan user
	//dan masukan ke dalam variable count
	err := q.DB.WithContext(ctx).
		Model(&domain.QuizAttempts{}).
		Where("user_id = ? and quiz_id = ?", userID, quizID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// query untuk mencari skor tertinggi di semua kuis
func (q QuizAttemptRepositoryImpl) FindHighestScoresByUserID(ctx context.Context, userID string) (map[uint]uint, error) {
	var results []domain.QuizHighestScore

	// SELECT quiz_id, MAX(score) as max_score FROM quiz_attempts WHERE user_id = ? GROUP BY quiz_id
	err := q.DB.WithContext(ctx).
		Model(&domain.QuizAttempts{}).
		Select("quiz_id, MAX(score) as max_score").
		Where("user_id = ?", userID).
		Group("quiz_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	scoreMap := make(map[uint]uint)
	for _, result := range results {
		scoreMap[result.QuizID] = result.MaxScore
	}

	return scoreMap, nil
}

// nyari quiz_attempt berdasarkan score / passingScore
func (q QuizAttemptRepositoryImpl) HasUserPassedQuizBefore(ctx context.Context, userID string, quizID uint, passingScore uint) (bool, error) {
	var count int64

	err := q.DB.WithContext(ctx).
		Model(&domain.QuizAttempts{}).
		Where("user_id = ? AND quiz_id = ? AND score >= ?", userID, quizID, passingScore).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
