package domain

type Questions struct {
	ID           uint   `json:"id" gorm:"primaryKey;unsigned; not null;autoIncrement"`
	QuizID       uint   `json:"quiz_id" gorm:"not null"`
	QuestionText string `json:"question_text" gorm:"type:text;not null"`

	//1:M BelongsTo quizzes
	Quizzes Quizzes `json:"quizzes" gorm:"foreignKey:quiz_id"`
	//1:M HasMany question_options
	QuestionOptions []QuestionOptions `json:"question_options" gorm:"foreignKey:question_id;references:id"`
}
