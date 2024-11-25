package main

import (
	"test-mnc/controller"
	"test-mnc/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	authController := controller.NewAuthController("data/users.json")
	walletController := controller.NewWalletController("data/wallet.json")

	r := gin.Default()
	r.POST("/api/auth/register", authController.RegisterHandler)
	r.POST("/api/auth/login", authController.LoginHandler)
	r.POST("/api/auth/refresh-token", authController.RefreshTokenHandler)
	r.GET("/api/auth/me", middleware.AuthMiddleware(), authController.SelfHandler)
	r.POST("/api/auth/reset-password", middleware.AuthMiddleware(), authController.ResetPasswordHandler)
	r.POST("/api/auth/logout", middleware.AuthMiddleware(), authController.LogoutHandler)

	r.POST("/api/wallet/register", middleware.AuthMiddleware(), walletController.RegisterHandler)
	r.POST("/api/wallet/transfer", middleware.AuthMiddleware(), walletController.TransferHandler)
	r.POST("/api/wallet/transfer-confirm", middleware.AuthMiddleware(), walletController.TransferConfirmHandler)
	r.POST("/api/wallet/top-up", middleware.AuthMiddleware(), walletController.TopUpHandler)
	r.POST("/api/wallet/withdraw", middleware.AuthMiddleware(), walletController.WithdrawHandler)
	r.Run(":8080")
}
