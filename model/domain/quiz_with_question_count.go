package domain

type QuizWithQuestionCount struct {
	ID            uint   `gorm:"column:id"`
	Title         string `gorm:"column:title"`
	Description   string `gorm:"column:description"`
	Level         uint   `gorm:"column:level"`
	Dialect       string `gorm:"column:dialect"`
	QuestionCount int    `gorm:"column:question_count"`
}
