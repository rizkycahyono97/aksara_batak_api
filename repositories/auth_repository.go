package repositories

import "github.com/rizkycahyono97/aksara_batak_api/model/domain"

type AuthRepository interface {
	FindUserByEmail(email string) (*domain.Users, error)
	CreateUser(user *domain.Users) error
}
