package config

import (
	"github.com/spf13/viper"
)

type Option func(i *info)

func UserConfigPath(path string) Option {
	return func(i *info) {
		i.userConfigPath = path
	}
}

func GlobalConfigPath(path string) Option {
	return func(i *info) {
		i.globalConfigPath = path
	}
}

func Load(options ...Option) error {
	i := new(info)
	// 设置配置
	for _, option := range options {
		option(i)
	}

	if len(i.globalConfigPath) > 0 {
		globalConfigViper = viper.New()
		globalConfigViper.SetConfigFile(i.globalConfigPath)
	}

	if len(i.userConfigPath) > 0 {
		userConfigViper = viper.New()
		userConfigViper.SetConfigFile(i.userConfigPath)
	}

	_ = userConfigViper.ReadInConfig()
	_ = globalConfigViper.ReadInConfig()
	//if rErr := globalConfigViper.ReadInConfig(); rErr != nil {
	//	if _, ok := rErr.(viper.ConfigFileNotFoundError); !ok {
	//		return errors.New("read config file:" + rErr.Error())
	//	}
	//}
	return nil
}

type info struct {
	userConfigPath   string
	globalConfigPath string
}