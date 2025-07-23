package domain

import "time"

type Quizzes struct {
	ID          uint      `json:"id" gorm:"primaryKey;unsigned;not null;autoIncrement"`
	LessonID    uint      `json:"lesson_id" gorm:"type:int;null;unsigned"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Level       uint      `json:"level" gorm:"type:int;not null;default:0"`
	Dialect     string    `json:"dialect" gorm:"type:enum('toba', 'karo');not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:DATETIME;column:updated_at; autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"column:deleted_at;type:DATETIME;index"`

	//1:M HasMany questions
	Questions []Questions `json:"questions" gorm:"foreignKey:quiz_id;references:id"`
	//1:M HasMany quiz_attempts
	QuizAttempts []QuizAttempts `json:"quiz_attempts" gorm:"foreignKey:quiz_id;references:id"`
	//1:M BelongsTo lessons
	Lessons Lessons `json:"lessons" gorm:"foreignKey:LessonID"`
}
