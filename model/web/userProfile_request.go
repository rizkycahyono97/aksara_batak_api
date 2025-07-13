package web

// field yang bisa diupdate oleh user
type UserProfileUpdateRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=255"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}
