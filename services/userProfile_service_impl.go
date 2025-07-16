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
	UserRepo         repositories.UserRepository
	UserProfileRepo  repositories.UserProfileRepository
	QuizAttemptsRepo repositories.QuizAttemptRepository
	Validate         *validator.Validate
	Log              *slog.Logger
}

func NewUserProfileService(
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	quizAttemptsRepo repositories.QuizAttemptRepository,
	validate *validator.Validate,
	log *slog.Logger,
) UserProfileService {
	return &UserProfileServiceImpl{
		UserRepo:         userRepo,
		UserProfileRepo:  userProfileRepo,
		QuizAttemptsRepo: quizAttemptsRepo,
		Validate:         validate,
		Log:              log,
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
	s.Log.InfoContext(ctx, "update user profile process started", "userID", userID)

	//validation
	if err := s.Validate.Struct(request); err != nil {
		s.Log.WarnContext(ctx, "update profile request validation failed", "error", err)
		return nil, err
	}

	//ambil data user yang di update
	user, err := s.UserRepo.FindUserByID(ctx, userID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find user before update", "error", err)
		return nil, errors.New("user not found")
	}

	user.Name = request.Name
	user.AvatarURL = request.AvatarURL

	if err := s.UserRepo.UserUpdate(ctx, user); err != nil {
		s.Log.ErrorContext(ctx, "failed to update user in repository", "error", err)
		return nil, err
	}
	s.Log.InfoContext(ctx, "user profile updated successfully", "userID", userID)

	return nil, err
}

func (s *UserProfileServiceImpl) GetMyAttempts(ctx context.Context, userID string) ([]web.QuizAttemptsResponse, error) {
	s.Log.InfoContext(ctx, "get user quiz attempts process started", "userID", userID)

	//repository
	attempts, err := s.QuizAttemptsRepo.FindAllQuizAttemptByUserID(ctx, userID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find quiz attempts from repository", "error", err, "userID", userID)
		return nil, err
	}

	// assign ke DTO
	var attemptsResponse []web.QuizAttemptsResponse
	for _, quizAttempt := range attempts {
		response := web.QuizAttemptsResponse{
			AttemptID:   quizAttempt.ID,
			QuizID:      quizAttempt.QuizID,
			QuizTitle:   quizAttempt.Quizzes.Title,
			Score:       quizAttempt.Score,
			CompletedAt: quizAttempt.CompletedAt,
		}
		attemptsResponse = append(attemptsResponse, response)
	}

	s.Log.InfoContext(ctx, "successfully retrieved user quiz attempts", "userID", userID, "count", len(attemptsResponse))
	return attemptsResponse, nil
}
