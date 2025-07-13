package repositories

import (
	"context"
	"errors"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*domain.Users, error) {
	var users domain.Users

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&users).Error; err != nil {
		return nil, errors.New("User not found")
	}
	return &users, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, users *domain.Users) error {
	return r.db.WithContext(ctx).Create(users).Error
}

func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, userID string) (*domain.Users, error) {
	var users domain.Users
	if err := r.db.WithContext(ctx).Where("uuid = ?", userID).First(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
