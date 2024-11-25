package model

type LogActivity struct {
	LogId     string `json:"log_id"`
	UserId    string `json:"user_id"`
	Status    int    `json:"status"`
	Endpoint  string `json:"endpoint"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
