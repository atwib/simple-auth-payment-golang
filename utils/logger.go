package utils

import (
	"test-mnc/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddLog(c *gin.Context, status int, message, endpoint string) model.LogActivity {
	return model.LogActivity{
		LogId:     uuid.New().String(),
		UserId:    c.GetString("sub"),
		Status:    status,
		Message:   message,
		Endpoint:  endpoint,
		CreatedAt: time.Now().Local().String(),
	}
}
