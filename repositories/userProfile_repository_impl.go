package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
	"time"
)

type UserProfileRepositoryImpl struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &UserProfileRepositoryImpl{db: db}
}

// method untuk membuat userProfile baru
func (u UserProfileRepositoryImpl) CreateUserProfile(ctx context.Context, profile *domain.UserProfiles) error {
	return u.db.WithContext(ctx).Create(profile).Error
}

func (u UserProfileRepositoryImpl) FindUserProfileByID(ctx context.Context, userID string) (domain.UserProfiles, error) {
	var profile domain.UserProfiles
	if err := u.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return domain.UserProfiles{}, err
	}
	return profile, nil
}

func (u UserProfileRepositoryImpl) UserProfileUpdate(ctx context.Context, profile *domain.UserProfiles) error {
	return u.db.WithContext(ctx).Save(&profile).Error
}

// method untuk memperbarui total_xp, current_streak dan las_active_at
// gorm Updates -> untuk memperbarui banyak kolom sekaligus
func (u UserProfileRepositoryImpl) UpdateXPAndStreak(ctx context.Context, userID string, xpToAdd int, newStreak uint, lastActive time.Time) error {
	err := u.db.WithContext(ctx).
		Model(&domain.UserProfiles{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_xp":       gorm.Expr("total_xp + ?", xpToAdd),
			"current_streak": newStreak,
			"last_active_at": lastActive,
		}).Error
	return err
}

// method untuk mencari xp tertinggi untuk setiap user
func (u UserProfileRepositoryImpl) GetTopUsers(ctx context.Context, limit int) ([]domain.UserWithProfile, error) {
	var results []domain.UserWithProfile

	// - Model: Menentukan tabel utama adalah user_profiles.
	// - Select: Memilih kolom yang kita butuhkan dari kedua tabel.
	// - Joins: Menggabungkan dengan tabel users berdasarkan user_id dan uuid.
	// - Order: Mengurutkan berdasarkan total_xp secara menurun.
	// - Limit: Membatasi jumlah hasil.
	// - Scan: Memasukkan hasil query ke dalam slice 'resultss'.
	err := u.db.WithContext(ctx).
		Model(&domain.UserProfiles{}).
		Select("users.uuid, users.name, users.avatar_url, user_profiles.total_xp").
		Joins("JOIN users ON users.uuid = user_profiles.user_id").
		Where("users.role != ?", "admin"). //jika role admin, jangan tampilkan
		Order("user_profiles.total_xp DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
