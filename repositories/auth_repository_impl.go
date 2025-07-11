package repositories

import (
	"context"
	"errors"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{db}
}

func (r *AuthRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*domain.Users, error) {
	var users domain.Users

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&users).Error; err != nil {
		return nil, errors.New("User not found")
	}
	return &users, nil
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, users *domain.Users) error {
	return r.db.WithContext(ctx).Create(users).Error
}
