package web

type LeaderboardResponse struct {
	Rank      int    `json:"rank"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	TotalXP   int    `json:"total_xp"`
}
