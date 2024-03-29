package routers

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/api"
)

// InitWorkRouter  初始化任务模块的路由
func InitWorkRouter(engine *gin.Engine) gin.IRoutes {
	var group = engine.Group("/works")
	{
		group.POST("", api.WorkCreateApi)
		group.DELETE("", api.WorkDeleteApi)
		group.GET("", api.WorkGetInfoApi)
		group.PATCH("/pause", api.WorkPauseApi)
		group.PATCH("/stop", api.WorkStopApi)
		group.PATCH("/restart", api.WorkRestartApi)
		group.POST("/list", api.WorkListApi)
		group.POST("/result/list", api.WorkResultListApi)
	}
	return group
}
