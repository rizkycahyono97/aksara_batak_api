package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"gorm.io/gorm"
	"log/slog"
)

type UserProfileServiceImpl struct {
	UserRepo        repositories.UserRepository
	UserProfileRepo repositories.UserProfileRepository
	Validate        *validator.Validate
	Log             *slog.Logger
}

func NewUserProfileService(
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	validate *validator.Validate,
	log *slog.Logger,
) UserProfileService {
	return &UserProfileServiceImpl{
		UserRepo:        userRepo,
		UserProfileRepo: userProfileRepo,
		Validate:        validate,
		Log:             log,
	}
}

func (s *UserProfileServiceImpl) FindUserProfileByID(ctx context.Context, userID string) (web.UserProfileResponse, error) {
	s.Log.InfoContext(ctx, "find user profile process started", "userID", userID)

	//ambil data dari user
	user, err := s.UserRepo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.InfoContext(ctx, "find user record not found", "userID", userID)
			return web.UserProfileResponse{}, errors.New("user not found")
		}
		s.Log.InfoContext(ctx, "failed to find user by id", "userID", userID)
		return web.UserProfileResponse{}, err
	}

	//ambil data dari user_profiles
	profile, err := s.UserProfileRepo.FindUserProfileByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.InfoContext(ctx, "find user profile record not found", "userID", userID)
			return web.UserProfileResponse{}, errors.New("user not found")
		}
		s.Log.InfoContext(ctx, "failed to find user profile by id", "userID", userID)
		return web.UserProfileResponse{}, err
	}

	//DTO response
	response := web.UserProfileResponse{
		UUID:          user.UUID,
		Name:          user.Name,
		Email:         user.Email,
		AvatarURL:     user.AvatarURL,
		Role:          user.Role,
		TotalXP:       int(profile.TotalXP),
		CurrentStreak: int(profile.CurrentStreak),
		LastActivaAt:  profile.LastActiveAt,
		JoinedAt:      user.CreatedAt,
	}

	s.Log.InfoContext(ctx, "successfully retrieved user profile", "userID", userID)
	return response, nil
}

func (s *UserProfileServiceImpl) UpdateUserProfile(ctx context.Context, userID string, request web.UserProfileUpdateRequest) (*web.UserProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}
