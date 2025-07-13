package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type UserProfileRepository interface {
	FindUserProfileByID(ctx context.Context, userID string) (domain.UserProfiles, error)
	UserProfileUpdate(ctx context.Context, profile *domain.UserProfiles) error
}
