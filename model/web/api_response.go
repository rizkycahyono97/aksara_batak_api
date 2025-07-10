package web

type ApiResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
}
