package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"time"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, profile *domain.UserProfiles) error
	FindUserProfileByID(ctx context.Context, userID string) (domain.UserProfiles, error)
	UserProfileUpdate(ctx context.Context, profile *domain.UserProfiles) error
	UpdateXPAndStreak(ctx context.Context, userID string, xpToAdd int, newStreak uint, lastActive time.Time) error
	GetTopUsers(ctx context.Context, limit int) ([]domain.UserWithProfile, error)
}
