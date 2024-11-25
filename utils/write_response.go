package utils

import (
	"test-mnc/dto"

	"github.com/gin-gonic/gin"
)

func WriteResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := dto.CommonResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
	c.Abort()
}
