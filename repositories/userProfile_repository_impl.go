package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
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
