package domain

type Questions struct {
	ID           uint   `json:"id" gorm:"primaryKey;unsigned; not null;autoIncrement"`
	QuizID       uint   `json:"quiz_id" gorm:"not null"`
	QuestionType string `json:"question_type" gorm:"type:enum('pilihan_ganda_aksara', 'pilihan_ganda_batak', 'nulis_aksara') default 'pilihan_ganda_batak';not null"`
	QuestionText string `json:"question_text" gorm:"type:text;not null"`
	ImageURL     string `json:"image_url" gorm:"type:varchar(255);null"`
	AudioURL     string `json:"audio_url" gorm:"type:varchar(255);null"`
	LottieURL    string `json:"lottie_url" gorm:"type:varchar(255);null"`

	//1:M BelongsTo quizzes
	Quizzes Quizzes `json:"quizzes" gorm:"foreignKey:quiz_id"`
	//1:M HasMany question_options
	QuestionOptions []QuestionOptions `json:"question_options" gorm:"foreignKey:question_id;references:id"`
}
