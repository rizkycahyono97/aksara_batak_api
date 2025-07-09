package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	Repo     repositories.AuthRepository
	Validate *validator.Validate
}

func NewAuthService(repo repositories.AuthRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		Repo:     repo,
		Validate: validate,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req web.RegisterUserRequest) (domain.Users, error) {
	//validation
	if err := validator.New().Struct(req); err != nil {
		return domain.Users{}, err
	}

	//cek email
	_, err := s.Repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		return domain.Users{}, errors.New("user already exists")
	}

	//hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.Users{}, errors.New("failed to hash password")
	}

	//buat object users baru
	userNew := domain.Users{
		UUID:         uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         "user", // default role
	}

	//simpan ke database
	err = s.Repo.CreateUser(ctx, &userNew)
	if err != nil {
		return domain.Users{}, errors.New("failed to create user")
	}

	return userNew, nil
}
