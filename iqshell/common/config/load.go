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

func LoadGlobalConfig(globalConfigPath string) *data.CodeError {
	if len(globalConfigPath) > 0 {
		globalConfigViper = viper.New()
		globalConfigViper.SetConfigFile(globalConfigPath)
		err := globalConfigViper.ReadInConfig()
		log.DebugF("read global config error:%v", err)
	}
	return nil
}

func LoadUserConfig(userConfigPath string) *data.CodeError {
	if len(userConfigPath) > 0 {
		userConfigViper = viper.New()
		userConfigViper.SetConfigFile(userConfigPath)
		err := userConfigViper.ReadInConfig()
		log.DebugF("read user config error:%v", err)
	}
	return nil
}
