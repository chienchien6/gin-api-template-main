package main

import (
	"RCSP/core"
	"RCSP/global"
	"RCSP/initialize"
)

func main() {

	// 初始化配置
	//initialize.InitViperConfig()

	initialize.InitConfig()
	//fmt.Println(global.GvaConfig.Test.Key)

	global.GvaLogger.Info(global.GvaConfig.Test.Key)

	// 启动服务
	core.RunServer()

}
