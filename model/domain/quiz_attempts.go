package domain

import "time"

type QuizAttempts struct {
	ID          uint      `json:"id" gorm:"primaryKey;unsigned;not null;autoIncrement"`
	UserID      string    `json:"user_id" gorm:"type:varchar(36);not null"`
	QuizID      uint      `json:"quiz_id" gorm:"type:int;unsigned;not null"`
	Score       uint      `json:"score" gorm:"type:int;not null;default:0"`
	CompletedAt time.Time `json:"completed_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	//1:M BelongsTo Users
	Users Users `json:"users" gorm:"foreignKey:user_id;references:uuid"`
	//1:M BelongsTo quizzes
	Quizzes Quizzes `json:"quizzes" gorm:"foreignKey:quiz_id;references:id"`
}
