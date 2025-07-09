package domain

import "time"

type UserProfiles struct {
	UserID        string    `json:"user_id" gorm:"primaryKey;type:varchar(36);not null"`
	TotalXP       uint      `json:"total_xp" gorm:"type:int;not null;default:0"`
	CurrentStreak uint      `json:"current_streak" gorm:"type:int;not null;default:0"`
	LastActiveAt  time.Time `json:"last_active_at" gorm:"type:date;null"`
}
