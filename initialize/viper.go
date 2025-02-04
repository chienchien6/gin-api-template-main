package initialize

import (
	"RCSP/global"
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ViperInit 初始化viper配置解析包，函数可接受命令行参数
func InitViperConfig() {
	var configFile string
	// 读取配置文件优先级: 命令行 > 默认值
	flag.StringVar(&configFile, "c", global.ConfigFile, "配置配置")
	if len(configFile) == 0 {
		// 读取默认配置文件
		panic("配置文件不存在！")
	}
	// 读取配置文件
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置解析失败:%s", err))
	}
	// 动态监测配置文件
	v.WatchConfig()
	//當配置文件發生變更時，執行指定的函數，並輸出變更的提示信息。
	v.OnConfigChange(func(in fsnotify.Event) { //fsnotify: 用於監控文件系統事件（如文件變更）
		fmt.Println("配置文件发生改变")

		//Unmarshal是一个用于将JSON格式的字符串解码为相应数据结构的函数。
		//初始加載配置文件到 global.GvaConfig 結構中
		if err := v.Unmarshal(&global.GvaConfig); err != nil {
			panic(fmt.Errorf("配置重载失败:%s", err))
		}
	})

	if err := v.Unmarshal(&global.GvaConfig); err != nil {
		panic(fmt.Errorf("配置重载失败:%s", err))
	}
	// 设置配置文件
	global.GvaConfig.App.ConfigFile = configFile

	/*
		//fmt.Printf("\n\n\n[   global.GvaConfig: %#v   ]\n\n\n", global.GvaConfig)
		//fmt.Println(reflect.TypeOf(global.GvaConfig).Kind())
		data := reflect.TypeOf(global.GvaConfig)
		num := reflect.ValueOf(global.GvaConfig)
		for i := 0; i < data.NumField(); i++ {
			fmt.Printf("%#v\n", num.Field(i).Interface())
			//println(data.Field(i).Name, ":", num.Field(i).Interface())
			//data2 := reflect.TypeOf(global.GvaConfig)
			//num2 := reflect.ValueOf(global.GvaConfig)
		}
	*/
}
