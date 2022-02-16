package main

import (
	"ccgpgov/common"
	"fmt"
	"os"
	"time"
	_ "github.com/joho/godotenv/autoload"
	"go.pfgit.cn/letsgo/xdev"
)
func Print() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
var Log = xdev.XNewLoggerDefault("./log/" + common.APP_NAME + ".log")
func main() {
	Log.Info(common.APP_NAME + "/" + common.APP_VERSION)

	if err := common.ReadConfig(); err != nil {
		os.Exit(1)
	}

	Log.Info("==========================config init ok=============================")
	Log.Infof("%+v", common.Config)
	xdev.SetLogLevel(common.Config.LogLevel)

	//TODO:
	Spider()
	IRISMonitor()
}
