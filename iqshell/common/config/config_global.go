package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
)

func GetGlobal() *Config {
	return &Config{
		Credentials: &auth.Credentials{
			AccessKey: getAccessKey(ConfigTypeGlobal),
			SecretKey: []byte(getSecretKey(ConfigTypeGlobal)),
		},
		UseHttps: getIsUseHttps(ConfigTypeGlobal),
		Hosts: &Hosts{
			UC:  GetUcHosts(ConfigTypeGlobal),
			Api: GetApiHosts(ConfigTypeGlobal),
			Rs:  GetRsHosts(ConfigTypeGlobal),
			Rsf: GetRsfHosts(ConfigTypeGlobal),
			Io:  GetIoHosts(ConfigTypeGlobal),
			Up:  GetUpHosts(ConfigTypeGlobal),
		},
	}
}
