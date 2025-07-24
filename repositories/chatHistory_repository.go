package repositories

import "github.com/rizkycahyono97/aksara_batak_api/model/domain"

type ChatHistoryRepository interface {
	Save(history domain.ChatHistory) (domain.ChatHistory, error)
	FindLastByUserID(userID string, limit int) ([]domain.ChatHistory, error)
}
