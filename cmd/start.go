package cmd

import (
	"gitlab.example.com/zhangweijie/tool-sdk/core"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/initizlize"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/option"
)

// Start 程序开始函数
func Start(option *option.Option) {
	err := initizlize.LoadConfig(option.ValidConfigFile)
	if err != nil {
		logger.Panic("加载配置文件出现错误", err)
	}
	global.ValidExecutorIns = option.ExecutorIns
	global.ValidModels = option.ValidModels
	global.ValidRouter = option.ValidRouters

	err = core.Start()
	if err != nil {
		logger.Panic("启动服务出现错误", err)
	}
}
