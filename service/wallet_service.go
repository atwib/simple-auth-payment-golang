package service

import (
	"errors"

	"net/http"
	"strconv"
	"strings"
	"test-mnc/constant"
	"test-mnc/dto"
	"test-mnc/model"
	"test-mnc/repository"
	"test-mnc/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletService struct {
	repo *repository.WalletRepository
}
type Confirmation struct {
	Amount     int           `json:"amount"`
	AdminFee   int           `json:"admin_fee"`
	FromWallet *model.Wallet `json:"from_wallet"`
	ToWallet   *model.Wallet `json:"to_wallet"`
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (w *WalletService) AddWallet(req dto.WalletRegisterRequest, c *gin.Context) (int, string, interface{}) {
	if req.Type != "Customer" && req.Type != "Merchant" {
		return http.StatusBadRequest, constant.InvalidRequest, constant.ErrWalletType
	}

	if req.PIN != req.PIN_confirm {
		return http.StatusBadRequest, constant.InvalidRequest, constant.ErrWalletPinConfirmNotMatch
	}

	if req.PIN < 1000 || req.PIN > 9999 {
		return http.StatusBadRequest, constant.InvalidRequest, nil
	}

	userId := c.GetString("sub")
	wallet, _ := w.repo.GetWalletByUserId(userId)
	if wallet != nil {
		return http.StatusConflict, constant.ErrWalletAlreadyExists, nil
	}

	hashedPin, err := utils.HashPassword(strconv.Itoa(req.PIN))
	if err != nil {
		return http.StatusInternalServerError, constant.ErrWalletFailedHashPin, nil
	}

	NewWallet := model.Wallet{
		WalletId: uuid.New().String(),
		PIN:      hashedPin,
		UserId:   userId,
		Type:     req.Type,
		Balance:  0,
	}

	err = w.repo.AddWallet(NewWallet)
	if err != nil {
		return http.StatusInternalServerError, constant.ErrWalletFailedRegister, nil
	}

	return http.StatusCreated, constant.SucWalletRegister, dto.WalletRespose{
		WalletId: NewWallet.WalletId,
		UserId:   NewWallet.UserId,
		Type:     NewWallet.Type,
		Balance:  NewWallet.Balance,
	}
}

func (w *WalletService) GetWalletByUserId(userId string) (*model.Wallet, error) {
	return w.repo.GetWalletByUserId(userId)
}

func (w *WalletService) GetWalletByWalletId(walletId string) (*model.Wallet, error) {
	return w.repo.GetWalletByWalletId(walletId)
}

func (w *WalletService) TopUpWallet(c *gin.Context, amount int) (int, string, interface{}) {
	userId := c.GetString("sub")

	if amount < constant.MinTopUp {
		return http.StatusBadRequest, strings.Replace(constant.ErrWalletMinTopUp, "_MINTOPUP_", strconv.Itoa(constant.MinTopUp), 1), nil
	}

	wallet, err := w.repo.GetWalletByUserId(userId)
	if err != nil {
		return http.StatusNotFound, err.Error(), nil
	}

	wallet.Balance += amount
	UpdatedWallet, err := w.repo.UpdateWallet(wallet)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}

	return http.StatusOK, constant.SucWalletTopUp, dto.WalletRespose{
		WalletId: UpdatedWallet.WalletId,
		UserId:   UpdatedWallet.UserId,
		Type:     UpdatedWallet.Type,
		Balance:  UpdatedWallet.Balance,
	}
}

func (w *WalletService) TransferConfirm(c *gin.Context, toWalletId string, amount int) (int, string, interface{}) {
	userId := c.GetString("sub")
	// Get FROM WALLET and check error
	fromWallet, errFromWallet := w.repo.GetWalletByUserId(userId)
	if errFromWallet != nil {
		return http.StatusNotFound, errFromWallet.Error(), nil
	}

	// Get TO WALLET and check error
	toWallet, errToWallet := w.repo.GetWalletByWalletId(toWalletId)
	if errToWallet != nil {
		return http.StatusNotFound, errToWallet.Error(), nil
	}

	// Calculate Amout + Admin Fee
	amountFee := amount + constant.AdminFeeTransfer

	// check if FROM WALLET Balance less than amount + admin fee return error
	// because insufficient balance
	if fromWallet.Balance < amountFee {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, gin.H{
			"balance":   fromWallet.Balance,
			"amount":    amount,
			"admin_fee": constant.AdminFeeTransfer,
		}
	}

	// check if FROM WALLET type is Customer or Merchant and Balance - (amount + admin fee) less than min saldo customer or min saldo merchant
	// return error because insufficient balance
	if fromWallet.Type == "Customer" && fromWallet.Balance-amountFee < constant.MinSaldoCustomer {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, gin.H{
			"balance":   fromWallet.Balance,
			"amount":    amount,
			"admin_fee": constant.AdminFeeTransfer,
		}
	} else if fromWallet.Type == "Merchant" && fromWallet.Balance-amountFee < constant.MinSaldoMerchant {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, gin.H{
			"balance":   fromWallet.Balance,
			"amount":    amount,
			"admin_fee": constant.AdminFeeTransfer,
		}
	}

	confirmation := Confirmation{
		Amount:     amount,
		AdminFee:   amountFee,
		FromWallet: fromWallet,
		ToWallet:   toWallet,
	}

	return http.StatusOK, constant.SucWalletTransferConfirm, confirmation

}

func (w *WalletService) Transfer(c *gin.Context, req dto.WalletTransferRequest) (int, string, interface{}) {
	userId := c.GetString("sub")
	toWalletId := req.ToWalletId
	amount := req.Amount
	pin := req.PIN
	// Get FROM WALLET and check error
	fromWallet, errFromWallet := w.repo.GetWalletByUserId(userId)
	if errFromWallet != nil {
		return http.StatusNotFound, errFromWallet.Error(), nil
	}

	// Get TO WALLET and check error
	toWallet, errToWallet := w.repo.GetWalletByWalletId(toWalletId)
	if errToWallet != nil {
		return http.StatusNotFound, errToWallet.Error(), nil
	}

	verify := utils.VerifyHashed(fromWallet.PIN, strconv.Itoa(pin))
	if !verify {
		return http.StatusBadRequest, constant.ErrWalletPinWrong, nil
	}

	// Calculate Amout + Admin Fee
	amountFee := amount + constant.AdminFeeTransfer

	// check if FROM WALLET Balance less than amount + admin fee return error
	// because insufficient balance
	if fromWallet.Balance < amountFee {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, nil
	}

	// check if FROM WALLET type is Customer or Merchant and Balance - (amount + admin fee) less than min saldo customer or min saldo merchant
	// return error because insufficient balance
	if fromWallet.Type == "Customer" && fromWallet.Balance-amountFee < constant.MinSaldoCustomer {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, nil
	} else if fromWallet.Type == "Merchant" && fromWallet.Balance-amountFee < constant.MinSaldoMerchant {
		return http.StatusBadRequest, constant.ErrWalletInsufficientBalance, nil
	}

	// update FROM WALLET and TO WALLET Balance
	fromWallet.Balance -= amountFee
	toWallet.Balance += amount
	w.repo.UpdateWallet(fromWallet)
	w.repo.UpdateWallet(toWallet)

	return http.StatusOK, constant.SucWalletTransfer, gin.H{
		"from_wallet": dto.WalletRespose{
			WalletId: fromWallet.WalletId,
			UserId:   fromWallet.UserId,
			Type:     fromWallet.Type,
			Balance:  fromWallet.Balance,
		},
		"to_wallet": toWallet,
		"amount":    amount,
		"admin_fee": amountFee,
	}
}

func (w *WalletService) WithdrawWallet(walletId string, amount int) (*model.Wallet, error) {
	// Get wallet by walletId
	wallet, err := w.repo.GetWalletByWalletId(walletId)
	if err != nil {
		return nil, err
	}

	// check if wallet type is Customer or Merchant and Balance - (amount + admin fee) less than min saldo customer or min saldo merchant
	if wallet.Type == "Customer" && wallet.Balance-amount+constant.AdminFeeCustomerWidraw < constant.MinSaldoCustomer {
		return nil, errors.New(constant.ErrWalletInsufficientBalance)
	}

	if wallet.Type == "Merchant" && wallet.Balance-amount+constant.AdminFeeMerchantWidraw < constant.MinSaldoMerchant {
		return nil, errors.New(constant.ErrWalletInsufficientBalance)
	}

	// check if wallet type is Customer then Balance - (amount + admin fee customer widraw)
	if wallet.Type == "Customer" {
		wallet.Balance -= amount + constant.AdminFeeCustomerWidraw
	}

	if wallet.Type == "Merchant" {
		wallet.Balance -= amount + constant.AdminFeeMerchantWidraw
	}

	return w.repo.UpdateWallet(wallet)
}
