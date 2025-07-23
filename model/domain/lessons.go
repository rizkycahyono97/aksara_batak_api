package domain

import "time"

type Lessons struct {
	ID          uint      `json:"id" gorm:"primaryKey;unsigned;not null;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;null"`
	IconURL     string    `json:"icon_url" gorm:"type:varchar(255);null"`
	OrderIndex  string    `json:"order_index" gorm:"type:integer;"`
	CreatedAT   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAT   time.Time `json:"updated_at" gorm:"type:timestamp;not null"`

	//1:M hasMany quizzes
	Quizzes []Quizzes `json:"quizzes" gorm:"foreignKey:LessonID;references:ID"`
}
