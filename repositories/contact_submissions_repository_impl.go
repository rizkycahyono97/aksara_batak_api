package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type contactSubmissionsRepositoryImpl struct {
	DB *gorm.DB
}

func NewContactSubmissionsRepository(db *gorm.DB) ContactSubmissionsRepository {
	return &contactSubmissionsRepositoryImpl{
		DB: db,
	}
}

func (c contactSubmissionsRepositoryImpl) Create(ctx context.Context, submission domain.ContactSubmissions) (domain.ContactSubmissions, error) {
	err := c.DB.WithContext(ctx).Create(&submission).Error
	return submission, err
}
