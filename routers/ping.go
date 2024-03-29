package routers

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/api"
)

func InitPingRouter(engine *gin.Engine) gin.IRoutes {
	var group = engine.Group("/ping")
	{
		group.GET("", api.PingApi)
	}
	return group
}
