package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
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
	if err := s.Validate.Struct(req); err != nil {
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

func (s *AuthServiceImpl) Login(ctx context.Context, req web.LoginUserRequest) (string, error) {
	//validation
	if err := s.Validate.Struct(req); err != nil {
		return "", err
	}

	//cari pengguna
	user, err := s.Repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	//membadingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.PasswordHash)); err != nil {
		return "", errors.New("invalid password")
	}

	//buat token jwt jika cocok
	expirationDate := time.Now().Add(72 * time.Hour) // expiration jwt token

	//jwt claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "lomba-batak-app",
		Subject:   user.UUID,
		ExpiresAt: jwt.NewNumericDate(expirationDate),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	//jwt sceret key
	secret := config.GetEnv("JWT_SECRET", "")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	//buat jwt token final
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//tandatangan jwt token dengan jwt secret key kita
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", errors.New("failed to sign token")
	}

	return signedToken, nil
}
