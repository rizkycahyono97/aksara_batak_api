package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.Users, error)
	CreateUser(ctx context.Context, user *domain.Users) error
	FindUserByID(ctx context.Context, userID string) (*domain.Users, error)
	UserUpdate(ctx context.Context, user *domain.Users) error
}
