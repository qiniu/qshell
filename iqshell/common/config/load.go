package config

import "github.com/spf13/viper"

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

func Load(options ...Option) {
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
}

type info struct {
	userConfigPath   string
	globalConfigPath string
}