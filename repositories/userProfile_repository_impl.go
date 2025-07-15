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
