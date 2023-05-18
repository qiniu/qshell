package config

import "github.com/qiniu/go-sdk/v7/auth"

func GetUser() *Config {
	return &Config{
		Credentials: &auth.Credentials{
			AccessKey: getAccessKey(ConfigTypeUser),
			SecretKey: []byte(getSecretKey(ConfigTypeUser)),
		},
		UseHttps: getIsUseHttps(ConfigTypeUser),
		Hosts: &Hosts{
			UC:    GetUcHosts(ConfigTypeUser),
			Api:   GetApiHosts(ConfigTypeUser),
			Rs:    GetRsHosts(ConfigTypeUser),
			Rsf:   GetRsfHosts(ConfigTypeUser),
			Io:    GetIoHosts(ConfigTypeUser),
			Up:    GetUpHosts(ConfigTypeUser),
			IoSrc: GetIoSrcHosts(ConfigTypeUser),
		},
	}
}
