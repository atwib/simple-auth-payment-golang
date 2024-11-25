package dto

type AuthRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserResponse struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}
