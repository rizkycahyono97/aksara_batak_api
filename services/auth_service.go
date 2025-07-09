package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type AuthService interface {
	Register(ctx context.Context, req web.RegisterUserRequest) (domain.Users, error)
}
