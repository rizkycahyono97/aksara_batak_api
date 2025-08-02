package repositories

import (
	"context"
	"fmt"
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

// simpan chat baru
func (c ChatHistoryRepositoryImpl) Create(ctx context.Context, chat *domain.ChatHistories) error {
	err := c.DB.WithContext(ctx).Create(&chat).Error
	if err != nil {
		return err
	}

	//auto trim jika melebihi batas
	go func() {
		_ = c.DeleteExcess(context.Background(), chat.UserID, 15)
	}()

	return nil
}

// mengambil 5 history terakhir
func (c ChatHistoryRepositoryImpl) GetLastFifteenByUserID(ctx context.Context, userID string) ([]domain.ChatHistories, error) {
	var histories []domain.ChatHistories
	err := c.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(15).
		Find(&histories).Error
	if err != nil {
		return []domain.ChatHistories{}, err
	}

	//reverse urutan
	for i, j := 0, len(histories)-1; i < j; i, j = i+1, j-1 {
		histories[i], histories[j] = histories[j], histories[i]
	}

	return histories, nil
}

// menghapus semua history user
func (c ChatHistoryRepositoryImpl) DeleteByUserID(ctx context.Context, userID string) error {
	return c.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.ChatHistories{}).Error
}

// menghitunh total history user
func (c ChatHistoryRepositoryImpl) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int64
	err := c.DB.WithContext(ctx).
		Model(&domain.ChatHistories{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}

// menghapus hisotory jika melebihi limit
func (c ChatHistoryRepositoryImpl) DeleteExcess(ctx context.Context, userID string, limit int) error {
	var idsToKeep []string

	err := c.DB.WithContext(ctx).
		Model(&domain.ChatHistories{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Pluck("id", &idsToKeep).Error
	if err != nil {
		return fmt.Errorf("gagal mengambil ID riwayat chat untuk disimpan: %w", err)
	}

	// Pengamanan: Jika karena suatu alasan tidak ada ID yang ditemukan untuk disimpan,
	// jangan lanjutkan proses hapus untuk menghindari penghapusan semua data.
	if len(idsToKeep) == 0 {
		return nil
	}

	return c.DB.WithContext(ctx).
		Where("user_id = ? AND id NOT IN ?", userID, idsToKeep).
		Delete(&domain.ChatHistories{}).Error
}
