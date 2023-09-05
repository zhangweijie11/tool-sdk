package routers

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/api"
)

// InitWorkRouter  初始化任务模块的路由
func InitWorkRouter(engine *gin.Engine) gin.IRoutes {
	var group = engine.Group("/works")
	{
		group.POST("", api.CreateWorkApi)
	}
	return group
}
