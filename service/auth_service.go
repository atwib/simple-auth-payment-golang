package service

import (
	"errors"
	"test-mnc/constant"
	"test-mnc/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"test-mnc/model"
	"test-mnc/repository"
)

var jwtKey = []byte("secret_key")
var refreshJwtKey = []byte("refresh_secret_key")

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(user model.User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New(constant.ErrAuthCannotBeEmpty)
	}
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New(constant.ErrAuthFailedHashPassword)
	}

	user.Password = hashPassword
	return s.repo.AddUser(user)
}

func (s *AuthService) GenerateTokens(userId, username string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iss": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iss": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(refreshJwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) VerifyRefreshToken(refreshTokenString string) (string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshJwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New(constant.ErrAuthInvalidRefreshToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New(constant.ErrAuthCouldNotParseToken)
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New(constant.ErrAuthInvalidTokenClaims)
	}

	return userId, nil
}

func (s *AuthService) AuthenticateUser(username, password string) (string, string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", errors.New(constant.ErrAuthInvalidUsernameOrPassword)
	}

	verify := utils.VerifyHashed(user.Password, password)
	if !verify {
		return "", "", errors.New(constant.ErrAuthInvalidUsernameOrPassword)
	}

	accessTokenString, refreshTokenString, err := s.GenerateTokens(user.UserId, user.Username)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (string, string, error) {

	userId, err := s.VerifyRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByUserId(userId)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := s.GenerateTokens(user.UserId, user.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ResetPassword(userId, oldPassword, newPassword string) error {
	user, err := s.repo.GetByUserId(userId)
	if err != nil {
		return err
	}

	verify := utils.VerifyHashed(user.Password, oldPassword)
	if !verify {
		return errors.New(constant.ErrAuthInvalidOldPassword)
	}

	hashPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	return s.repo.UpdateUser(user)
}

func (s *AuthService) GetByUserId(userId string) (*model.User, error) {
	return s.repo.GetByUserId(userId)
}

func (s *AuthService) Logout(tokenString string) error {
	return utils.AddTokenToBlacklist(tokenString)
}
