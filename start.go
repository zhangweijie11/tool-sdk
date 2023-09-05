package tool_sdk

import (
	web "gitlab.example.com/zhangweijie/tool-sdk/cmd"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/initizlize"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
)

func Start(paramsInterface ...global.ParamsInterface) {
	err := initizlize.LoadConfig("config.yaml")
	if err != nil {
		logger.Panic("加载配置文件出现错误", err)
	}
	if len(paramsInterface) > 0 {
		global.CheckParamsInterface(paramsInterface[0])
	}

	err = web.Start()
	if err != nil {
		logger.Panic("启动服务出现错误", err)
	}
}
