package web

import (
	"time"
)

//type ChatbotResponse struct {
//	Reply   string                 `json:"reply"`
//	History []domain.ChatHistories `json:"history"`
//}

type ChatbotResponse struct {
	Reply   string          `json:"reply"`
	History []ChatHistories `json:"history"`
}

type ChatHistories struct {
	ID      int64  `json:"id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Reply   string `json:"reply"`
	//MessageType string    `json:"message_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Users object disembunyikan, atau bisa juga ditambahkan json:"-"
	// Users        UserResponse   `json:"users,omitempty"` // jika ingin ditampilkan nanti
}

type ChatHistoriesItemResponse struct {
	//Role      string    `json:"role"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
