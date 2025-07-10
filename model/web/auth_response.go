package web

type RegisterUserResponse struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

//func RegisterUserResponseFromModel(user *domain.Users) RegisterUserResponse {
//	return RegisterUserResponse{
//		UUID:      user.UUID,
//		Name:      user.Name,
//		Email:     user.Email,
//		AvatarURL: user.AvatarURL,
//		Role:      user.Role,
//	}
//}
