package main

import (
	web "gitlab.example.com/zhangweijie/tool-sdk/cmd"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/initizlize"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
)

// Start 程序开始函数
func main() {
	err := initizlize.LoadConfig("config.yaml")
	if err != nil {
		logger.Panic("加载配置文件出现错误", err)
	}
	global.ValidParamsIns = global.NewParamsIns()
	err = web.Start()
	if err != nil {
		logger.Panic("启动服务出现错误", err)
	}
}
