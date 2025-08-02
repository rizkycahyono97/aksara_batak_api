package domain

import "time"

type ContactSubmissions struct {
	ID        uint `json:"id" gorm:"type:bigint;unsigned;not null;autoIncrement:true;unique;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Name      string `json:"name" gorm:"type:varchar(255);not null"`
	Email
	Message string `json:"message" gorm:"type:text;not null"`
	Status  string `json:"status" gorm:"type:varchar(50);not null;default 'baru'"`
}
