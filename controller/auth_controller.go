package controller

import (
	"net/http"
	"test-mnc/constant"

	"test-mnc/dto"
	"test-mnc/model"
	"test-mnc/repository"
	"test-mnc/service"
	"test-mnc/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(dbFile string) *AuthController {
	repo := repository.NewAuthRepository(dbFile)
	authService := service.NewAuthService(repo)
	return &AuthController{service: authService}
}

func (uc *AuthController) RegisterHandler(c *gin.Context) {
	var req dto.AuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}
	user := model.User{
		UserId:   uuid.New().String(),
		Username: req.Username,
		Password: req.Password,
	}

	err := uc.service.RegisterUser(user)
	if err != nil {
		utils.WriteResponse(c, http.StatusInternalServerError, constant.ErrAuthFailedRegisterUser, gin.H{"error": err.Error()})
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthRegister, map[string]interface{}{"user": user})
}

func (uc *AuthController) LoginHandler(c *gin.Context) {
	var req dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := uc.service.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedLogin, gin.H{"error": err.Error()})
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthLogin, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *AuthController) RefreshTokenHandler(c *gin.Context) {
	req := dto.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := uc.service.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.WriteResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthRefreshToken, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *AuthController) SelfHandler(c *gin.Context) {
	userId := c.GetString("sub")

	if userId == "" {
		utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedUnauthorized, nil)
		return
	}

	user, err := uc.service.GetByUserId(userId)
	if err != nil {
		utils.WriteResponse(c, http.StatusInternalServerError, constant.ErrAuthFailedSelf, gin.H{"error": err.Error()})
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthSelf, dto.UserResponse{
		UserId:   user.UserId,
		Username: user.Username,
	})
}

func (uc *AuthController) ResetPasswordHandler(c *gin.Context) {
	userId := c.GetString("sub")

	if userId == "" {
		utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedUnauthorized, nil)
		return
	}

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteResponse(c, http.StatusBadRequest, constant.InvalidRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.service.ResetPassword(userId, req.OldPassword, req.NewPassword)
	if err != nil {
		utils.WriteResponse(c, http.StatusInternalServerError, constant.ErrAuthFailedResetPassword, gin.H{"error": err.Error()})
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthResetPassword, nil)
}

func (uc *AuthController) LogoutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedUnauthorized, nil)
		return
	}

	err := uc.service.Logout(tokenString)
	if err != nil {
		utils.WriteResponse(c, http.StatusInternalServerError, constant.ErrAuthFailedBlacklistToken, gin.H{"error": err.Error()})
		return
	}

	utils.WriteResponse(c, http.StatusOK, constant.SucAuthLogout, nil)
}
