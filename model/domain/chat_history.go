package domain

import "time"

type ChatHistories struct {
	ID          uint      `json:"id" gorm:"type:integer;autoIncrement;primaryKey"`
	UserID      string    `json:"user_id" gorm:"type:varchar(36);not null"`
	Message     string    `json:"message" gorm:"type:text;not null"`
	Reply       string    `json:"reply" gorm:"type:text;not null"`
	MessageType string    `json:"message_type" gorm:"type:varchar(30)"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;not null"`

	//1:M BelongsTo users
	Users Users `json:"users" gorm:"foreignKey:user_id;references:uuid;OnDelete:CASCADE"`
}
