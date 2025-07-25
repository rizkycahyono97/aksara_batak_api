package domain

type QuestionOptions struct {
	ID         uint   `json:"id" gorm:"primaryKey;unsigned;not null;autoIncrement"`
	QuestionID uint   `json:"question_id" gorm:"unsigned;not null"`
	OptionText string `json:"option_text" gorm:"type:text;not null"`
	AksaraText string `json:"aksara_text" gorm:"type:text;null"`
	ImageURL   string `json:"image_url" gorm:"type:varchar(255);null"`
	AudioURL   string `json:"audio_url" gorm:"type:varchar(255);null"`
	IsCorrect  bool   `json:"is_correct" gorm:"type:boolean;not null;default:false"`

	//1:M BelongsTo Question
	Questions Questions `json:"questions" gorm:"foreignKey:question_id"`
}
