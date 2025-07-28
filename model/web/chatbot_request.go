package web

type ChatbotRequest struct {
	Userid  string `json:"user_id"`
	Message string `json:"message"`
}
