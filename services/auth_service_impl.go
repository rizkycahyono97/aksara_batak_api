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

// dependency
type AuthServiceImpl struct {
	Repo     repositories.AuthRepository
	Validate *validator.Validate
	Log      *slog.Logger
}

// dependncy injection
func NewAuthService(repo repositories.AuthRepository, validate *validator.Validate, log *slog.Logger) AuthService {
	return &AuthServiceImpl{
		Repo:     repo,
		Validate: validate,
		Log:      log,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req web.RegisterUserRequest) (domain.Users, error) {
	s.Log.InfoContext(ctx, "register process started", "name", req.Name, "email", req.Email)
	//validation
	if err := s.Validate.Struct(req); err != nil {
		s.Log.ErrorContext(ctx, "validation failed for register request", "error", err)
		return domain.Users{}, err
	}

	//cek email
	_, err := s.Repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		s.Log.WarnContext(ctx, "register attempt failed: email has been taken", "email", req.Email)
		return domain.Users{}, errors.New("email already exists")
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.Log.WarnContext(ctx, "login attempt failed: password mismatch", "email", req.Email, "userID", user.UUID)
		return "", errors.New("invalid email or password")
	}

	//===============
	//jwt token
	//===============
	//create new token with the specified signing method and claims
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &web.JwtCustomClaims{
		UUID: user.UUID,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "lomba-batak-app",
			Subject:   user.Name,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims custom
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//secret key
	secret := config.GetEnv("JWT_SECRET_KEY", "")
	if secret == "" {
		s.Log.ErrorContext(ctx, "JWT_SECRET_KEY not configured", "error", err)
		return "", errors.New("JWT_SECRET_KEY not configured")
	}

	//signed token with secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to sign token", "error", err)
		return "", errors.New("failed to sign token")
	}

	return tokenString, nil
}
