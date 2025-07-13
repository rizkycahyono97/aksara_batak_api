package web

import "time"

type UserProfileResponse struct {
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	TotalXP       int       `json:"total_xp"`
	CurrentStreak int       `json:"current_streak"`
	LastActivaAt  time.Time `json:"last_activated_at"`
	JoinedAt      time.Time `json:"joined_at"`
}
