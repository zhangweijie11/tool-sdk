package routers

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/example/api"
)

func InitTechnologiesRouter(engine *gin.Engine) gin.IRoutes {
	var group = engine.Group("/technologies")
	{
		group.GET("", api.TechnologiesCreateApi)
	}
	return group
}
