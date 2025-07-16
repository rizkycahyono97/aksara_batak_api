package domain

type UserWithProfile struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	TotalXP   uint   `json:"total_xp"`
}
