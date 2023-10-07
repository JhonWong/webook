package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	initVipper()
	initLogger()

	server := InitWebServer()

	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initVipper() {
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 重新加载配置内容
		fmt.Sprintln("Config file changed")
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
