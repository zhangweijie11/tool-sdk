package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
)

// TechnologiesCreateApi 入库指纹数据
func TechnologiesCreateApi(c *gin.Context) {
	schemas.SuccessCreate(c, "SUCCESS")
	return
}
