package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
)

// PingApi 服务连通性测试
func PingApi(c *gin.Context) {
	schemas.Success(c, "pong")
	return
}
