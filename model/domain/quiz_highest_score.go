package domain

type QuizHighestScore struct {
	QuizID   uint `gorm:"column:quiz_id"`
	MaxScore uint `gorm:"column:max_score"`
}
