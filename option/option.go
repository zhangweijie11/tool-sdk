package option

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/routers"
)

type Option struct {
	ExecutorIns  global.ExecutorInterface
	ValidModels  []interface{}
	ValidRouters []func(*gin.Engine) gin.IRoutes
	ValidConfig  []byte
}

var defaultOption = &Option{
	ExecutorIns:  global.ValidExecutorIns,
	ValidModels:  []interface{}{&models.Work{}, &models.Task{}, &models.Result{}},
	ValidRouters: []func(*gin.Engine) gin.IRoutes{routers.InitPingRouter, routers.InitWorkRouter},
}

func GetDefaultOption() *Option {
	return defaultOption
}
