package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

type LoadInfo struct {
	UserConfigPath   string
	GlobalConfigPath string
}

func Load(info LoadInfo) error {
	if len(info.GlobalConfigPath) > 0 {
		globalConfigViper = viper.New()
		globalConfigViper.SetConfigFile(info.GlobalConfigPath)
	}

	if len(info.UserConfigPath) > 0 {
		userConfigViper = viper.New()
		userConfigViper.SetConfigFile(info.UserConfigPath)
	}

	if err := userConfigViper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return errors.New("read user config error:" + err.Error())
		}
	}

	if err := globalConfigViper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return errors.New("read global config error:" + err.Error())
		}
	}

	return nil
}
