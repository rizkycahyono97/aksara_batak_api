package web

type AuthRequest struct {
	Name string `json:"name" binding:"required"`
}
