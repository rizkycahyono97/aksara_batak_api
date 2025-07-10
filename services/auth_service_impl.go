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
	"log/slog"
	"time"
)

type AuthServiceImpl struct {
	Repo     repositories.AuthRepository
	Validate *validator.Validate
	Log      *slog.Logger
}

func NewAuthService(repo repositories.AuthRepository, validate *validator.Validate, log *slog.Logger) AuthService {
	return &AuthServiceImpl{
		Repo:     repo,
		Validate: validate,
		Log:      log,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req web.RegisterUserRequest) (domain.Users, error) {
	s.Log.InfoContext(ctx, "register process started", "name", req.Password, "email", req.Email)
	//validation
	if err := s.Validate.Struct(req); err != nil {
		s.Log.ErrorContext(ctx, "validation failed for register request", "error", err)
		return domain.Users{}, err
	}

	//cek email
	_, err := s.Repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		s.Log.WarnContext(ctx, "register attempt failed: email has been taken", "email", req.Email)
		return domain.Users{}, errors.New("user already exists")
	}

	//hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.WarnContext(ctx, "failed to hash password", "password", req.Password)
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
	s.Log.DebugContext(ctx, "assign to new object user", "user", userNew)

	//simpan ke database
	err = s.Repo.CreateUser(ctx, &userNew)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to create new user", "error", err)
		return domain.Users{}, errors.New("failed to create user")
	}

	s.Log.DebugContext(ctx, "user registered successfully", "user", userNew)
	return userNew, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req web.LoginUserRequest) (string, error) {
	s.Log.InfoContext(ctx, "login process started", "email", req.Email)
	//validation
	if err := s.Validate.Struct(req); err != nil {
		s.Log.ErrorContext(ctx, "validation failed for login request", "error", err)
		return "", err
	}

	//cari pengguna
	user, err := s.Repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		s.Log.WarnContext(ctx, "login attempt failed: user not found", "email", req.Email)
		return "", errors.New("user not found")
	}

	//membadingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.PasswordHash)); err != nil {
		s.Log.WarnContext(ctx, "login attempt failed: password mismatch", "email", req.Email, "userID", user.UUID)
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
	s.Log.InfoContext(ctx, "user logged in successfully", "userID", user.UUID)

	return signedToken, nil
}
