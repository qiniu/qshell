package config

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/viper"
)

type LoadInfo struct {
	UserConfigPath   string
	GlobalConfigPath string
}

func Load(info LoadInfo) *data.CodeError {
	if len(info.GlobalConfigPath) > 0 {
		globalConfigViper = viper.New()
		globalConfigViper.SetConfigFile(info.GlobalConfigPath)
		err := globalConfigViper.ReadInConfig()
		log.DebugF("read global config error:%v", err)
	}

	if len(info.UserConfigPath) > 0 {
		userConfigViper = viper.New()
		userConfigViper.SetConfigFile(info.UserConfigPath)
		err := userConfigViper.ReadInConfig()
		log.DebugF("read user config error:%v", err)
	}

	return nil
}
