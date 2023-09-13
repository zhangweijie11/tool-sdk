package option

import (
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
)

type Option struct {
	ExecutorIns global.ExecutorInterface
	ValidModels []interface{}
}

var defaultOption = &Option{
	ExecutorIns: global.ValidExecutorIns,
	ValidModels: []interface{}{&models.Work{}, &models.Task{}, &models.Result{}},
}

func GetDefaultOption() *Option {
	return defaultOption
}
