package dto

import "test-mnc/model"

type WalletRegisterRequest struct {
	UserId      string `json:"user_id"`
	Type        string `json:"type"`
	PIN         int    `json:"pin"`
	PIN_confirm int    `json:"pin_confirm"`
}

type WalletRespose struct {
	WalletId string `json:"wallet_id"`
	UserId   string `json:"user_id"`
	Type     string `json:"type"`
	Balance  int    `json:"balance"`
}

type WalletTopUpRequest struct {
	Amount int `json:"amount"`
}

type WalletTransferRequest struct {
	ToWalletId string `json:"to_wallet_id"`
	Amount     int    `json:"amount"`
	PIN        int    `json:"pin"`
}

type WalletTransferConfirmRequest struct {
	ToWalletId string `json:"to_wallet_id"`
	Amount     int    `json:"amount"`
}

type WalletTransferResponse struct {
	FromWallet *model.Wallet `json:"from_wallet"`
	ToWallet   *model.Wallet `json:"to_wallet"`
	Amount     int           `json:"amount"`
}

type WalletWithdrawRequest struct {
	Amount int `json:"amount"`
}
