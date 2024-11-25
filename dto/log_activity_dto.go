package dto

type LogActivityRequest struct {
	LogId       string `json:"log_id"`
	UserId      string `json:"user_id"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	CreatedAt   string `json:"created_at"`
}
