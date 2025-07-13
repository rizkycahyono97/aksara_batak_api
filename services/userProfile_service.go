package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type UserProfileService interface {
	FindUserProfileByID(ctx context.Context, userID string) (web.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, userID string, request web.UserProfileUpdateRequest) (*web.UserProfileResponse, error)
}
