package domain

import "time"

type ChatHistory struct {
	UserID    string    `json:"user_id" gorm:"type:varchar(36);not null"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"type:timestamp;not null"`

	//1:1 HasOne
	Users Users `json:"users" gorm:"foreignKey:user_id;references:uuid;OnDelete:CASCADE"`
}
