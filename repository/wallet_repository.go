package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"test-mnc/model"
)

type WalletRepository struct {
	mu     sync.Mutex
	dbFile string
	wallet []model.Wallet
}

func NewWalletRepository(dbFile string) *WalletRepository {
	repo := &WalletRepository{dbFile: dbFile}
	repo.loadWallet()
	return repo
}

func (r *WalletRepository) loadWallet() error {
	if _, err := os.Stat(r.dbFile); os.IsNotExist(err) {
		if err := os.WriteFile(r.dbFile, []byte("[]"), 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(r.dbFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.wallet)
}

func (r *WalletRepository) saveWallet() error {
	data, err := json.MarshalIndent(r.wallet, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dbFile, data, 0644)
}

func (r *WalletRepository) AddWallet(wallet model.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.wallet = append(r.wallet, wallet)
	return r.saveWallet()
}

func (r *WalletRepository) GetWalletByUserId(userId string) (*model.Wallet, error) {
	for i := range r.wallet {
		if r.wallet[i].UserId == userId {
			return &r.wallet[i], nil
		}
	}
	return nil, errors.New("wallet not found")
}

func (r *WalletRepository) GetWalletByWalletId(walletId string) (*model.Wallet, error) {
	for i := range r.wallet {
		if r.wallet[i].WalletId == walletId {
			return &r.wallet[i], nil
		}
	}
	return nil, errors.New("wallet not found")
}

func (r *WalletRepository) UpdateWallet(wallet *model.Wallet) (*model.Wallet, error) {
	for i, w := range r.wallet {
		if w.WalletId == wallet.WalletId {
			r.wallet[i] = *wallet
			r.saveWallet()
			return wallet, nil
		}
	}
	return nil, nil
}
