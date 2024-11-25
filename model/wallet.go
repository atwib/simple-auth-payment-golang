package model

type Wallet struct {
	WalletId string `json:"wallet_id"`
	PIN      string `json:"pin"`
	UserId   string `json:"user_id"`
	Type     string `json:"type"`
	Balance  int    `json:"balance"`
}
