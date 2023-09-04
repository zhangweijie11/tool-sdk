package main

import (
	web "tool-sdk/cmd"
	"tool-sdk/initizlize"
	"tool-sdk/middleware/logger"
)

func main() {
	err := initizlize.LoadConfig("config.yaml")
	if err != nil {
		logger.Panic("加载配置文件出现错误", err)
	}
	err = web.Start()
	if err != nil {
		logger.Panic("启动服务出现错误", err)
	}
}
