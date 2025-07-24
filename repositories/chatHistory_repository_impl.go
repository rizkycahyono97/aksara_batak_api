package repositories

import (
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type ChatHistoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewChatHistoryRepository(db *gorm.DB) ChatHistoryRepository {
	return &ChatHistoryRepositoryImpl{
		DB: db,
	}
}

func (c ChatHistoryRepositoryImpl) Save(history domain.ChatHistory) (domain.ChatHistory, error) {
	err := c.DB.Create(&history).Error

	return history, err
}

func (c ChatHistoryRepositoryImpl) FindLastByUserID(userID string, limit int) ([]domain.ChatHistory, error) {
	var histories []domain.ChatHistory

	err := c.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&histories).Error
	if err != nil {
		return histories, err
	}

	return histories, nil
}
