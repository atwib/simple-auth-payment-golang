package middleware

import (
	"net/http"
	"strings"
	"test-mnc/constant"
	"test-mnc/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedUnauthorized, nil)
			return
		}

		// Pisahkan "Bearer" dan token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedUnauthorizedFormat, nil)
			return
		}
		tokenString := parts[1]

		isBlacklisted, err := utils.IsTokenBlacklisted(authHeader)
		if err != nil {
			utils.WriteResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		if isBlacklisted {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthFailedBlacklistToken, nil)
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthInvalidOrExpiredToken, nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthCouldNotParseToken, nil)
			return
		}

		if sub, ok := claims["sub"].(string); ok {
			c.Set("sub", sub)
		} else {
			utils.WriteResponse(c, http.StatusUnauthorized, constant.ErrAuthInvalidTokenClaims, nil)
			return
		}

		c.Next()
	}
}
