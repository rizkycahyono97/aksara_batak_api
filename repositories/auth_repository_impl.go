package repositories

import (
	"errors"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{db}
}

func (r *AuthRepositoryImpl) FindUserByEmail(email string) (*domain.Users, error) {
	var user domain.Users

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}

func (r *AuthRepositoryImpl) CreateUser(user *domain.Users) error {
	return r.db.Create(user).Error
}
