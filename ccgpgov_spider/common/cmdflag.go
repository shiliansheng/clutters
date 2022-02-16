package common

import (
	"flag"
	"fmt"
	"os"
)

type XCMDFlag struct {
	h, H bool
	v    bool

	//TODO:根据需要增加其他命令行参数
}

var CMDFlag XCMDFlag

func init() {
	flag.BoolVar(&CMDFlag.h, "h", false, "帮助信息")
	flag.BoolVar(&CMDFlag.H, "H", false, "帮助信息")
	flag.BoolVar(&CMDFlag.v, "v", false, "显示版本号")

	//TODO:其他命令行参数设置

	flag.Parse()

	if CMDFlag.h || CMDFlag.H {
		flag.Usage()
		os.Exit(1)
	}

	if CMDFlag.v {
		fmt.Println(APP_NAME + "/" + APP_VERSION)
		os.Exit(1)
	}
}
