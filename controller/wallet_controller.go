package controller

import (
	"net/http"
	"test-mnc/constant"
	"test-mnc/dto"
	"test-mnc/repository"
	"test-mnc/service"
	"test-mnc/utils"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service *service.WalletService
}

func NewWalletController(dbFile string) *WalletController {
	repo := repository.NewWalletRepository(dbFile)
	walletService := service.NewWalletService(repo)
	return &WalletController{service: walletService}
}

func (wc *WalletController) RegisterHandler(c *gin.Context) {
	var req dto.WalletRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, nil)
		return
	}
	statusCode, message, data := wc.service.AddWallet(req, c)
	utils.WriteResponse(c, statusCode, message, data)
}

func (wc *WalletController) TopUpHandler(c *gin.Context) {
	var req dto.WalletTopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, nil)
		return
	}
	status, message, data := wc.service.TopUpWallet(c, req.Amount)
	utils.WriteResponse(c, status, message, data)
}

func (wc *WalletController) TransferConfirmHandler(c *gin.Context) {
	var req dto.WalletTransferConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}
	status, message, data := wc.service.TransferConfirm(c, req.ToWalletId, req.Amount)
	utils.WriteResponse(c, status, message, data)
}
func (wc *WalletController) TransferHandler(c *gin.Context) {
	var req dto.WalletTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}

	status, message, err := wc.service.Transfer(c, req)
	utils.WriteResponse(c, status, message, err)
}
func (wc *WalletController) WithdrawHandler(c *gin.Context) {
	var req dto.WalletWithdrawRequest
}