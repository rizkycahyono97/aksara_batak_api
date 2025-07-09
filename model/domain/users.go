package domain

import "time"

type Users struct {
	UUID         string    `json:"uuid" gorm:"primaryKey;type:varchar(36);not null"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Email        string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	PasswordHash string    `json:"password_hash" gorm:"type:varchar(255);not null"`
	AvatarURL    string    `json:"avatar_url" gorm:"type:varchar(255);"`
	Role         string    `json:"role" gorm:"type:enum('user', 'admin') default 'user';not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:DATETIME;column:updated_at; autoUpdateTime"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"column:deleted_at;type:DATETIME;index"`

	//1:1 BelongsTo user_profiles
	UserProfiles UserProfiles `json:"user_profiles" gorm:"foreignKey:user_id;references:uuid;OnDelete:CASCADE"`
}
