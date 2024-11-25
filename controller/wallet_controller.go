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
	service     *service.WalletService
	logActivity *repository.LogActivityRepository
}

func NewWalletController(dbFile string) *WalletController {
	repo := repository.NewWalletRepository(dbFile)
	walletService := service.NewWalletService(repo)
	logActivityRepo := repository.NewLogActivityRepository("data/log_activity.json")
	return &WalletController{service: walletService, logActivity: logActivityRepo}
}

func (wc *WalletController) RegisterHandler(c *gin.Context) {
	var req dto.WalletRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wc.logActivity.AddLog(utils.AddLog(c, http.StatusBadRequest, constant.InvalidRequest, c.Request.RequestURI))
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, nil)
		return
	}
	statusCode, message, data := wc.service.AddWallet(req, c)

	wc.logActivity.AddLog(utils.AddLog(c, statusCode, message, c.Request.RequestURI))
	utils.WriteResponse(c, statusCode, message, data)
}

func (wc *WalletController) TopUpHandler(c *gin.Context) {
	var req dto.WalletTopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wc.logActivity.AddLog(utils.AddLog(c, http.StatusBadRequest, constant.InvalidRequest, c.Request.RequestURI))
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, nil)
		return
	}
	status, message, data := wc.service.TopUpWallet(c, req.Amount)

	wc.logActivity.AddLog(utils.AddLog(c, status, message, c.Request.RequestURI))
	utils.WriteResponse(c, status, message, data)
}

func (wc *WalletController) TransferConfirmHandler(c *gin.Context) {
	var req dto.WalletTransferConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wc.logActivity.AddLog(utils.AddLog(c, http.StatusBadRequest, constant.InvalidRequest, c.Request.RequestURI))
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}
	status, message, data := wc.service.TransferConfirm(c, req.ToWalletId, req.Amount)

	wc.logActivity.AddLog(utils.AddLog(c, status, message, c.Request.RequestURI))
	utils.WriteResponse(c, status, message, data)
}

func (wc *WalletController) TransferHandler(c *gin.Context) {
	var req dto.WalletTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wc.logActivity.AddLog(utils.AddLog(c, http.StatusBadRequest, constant.InvalidRequest, c.Request.RequestURI))
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}

	status, message, err := wc.service.Transfer(c, req)

	wc.logActivity.AddLog(utils.AddLog(c, status, message, c.Request.RequestURI))
	utils.WriteResponse(c, status, message, err)
}

func (wc *WalletController) WithdrawHandler(c *gin.Context) {
	var req dto.WalletWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wc.logActivity.AddLog(utils.AddLog(c, http.StatusBadRequest, constant.InvalidRequest, c.Request.RequestURI))
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}
	status, message, data := wc.service.WithdrawWallet(c, req.Amount)

	wc.logActivity.AddLog(utils.AddLog(c, status, message, c.Request.RequestURI))
	utils.WriteResponse(c, status, message, data)
}
