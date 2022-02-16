package common

import (
	"strings"

	"go.pfgit.cn/letsgo/xdev"
)

type XConfig struct {
	ConnStr  	string 		`key:"common.db_conn" commit:"false"`
	LogLevel 	string 		`key:"common.log_level" default:"INFO" commit:"false"`
	TimeSpan 	int			`key:"common.time_span" commit:"false"`
	MonitorPort string		`key:"common.monitor_port" commit:"false"`
	TaskRunTime	string 		`key:"common.task_run_time" commit:"false"`
	LocalIp		string		`commit:"false"`
	Account  	[]string	`key:"atinfo.account" commit:"false"`
	Nickname 	[]string	`key:"atinfo.nickname" commit:"false"`
}

var Config XConfig
func ReadConfig() error {
	err := xdev.ReadConfig(APP_CONFIG_PATH, APP_NAME, &Config)
	dbIP := Config.ConnStr
	Config.LocalIp = xdev.GetLocalIP(dbIP[strings.IndexByte(dbIP, '@') + 1 : strings.LastIndexByte(dbIP, ':') - 1])
	return err
}
