package workspace

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func defaultConfig() *config.Config {
	return &config.Config{
		Credentials: &auth.Credentials{
			AccessKey: "",
			SecretKey: nil,
		},
		UseHttps: data.NewBool(false),
		Hosts: &config.Hosts{
			UC: []string{"uc.qbox.me"},
		},
		Log: &config.LogSetting{
			LogLevel:  data.NewString(config.InfoKey),
			LogFile:   nil,
			LogRotate: data.NewInt(7),
			LogStdout: data.NewBool(true),
		},
	}
}

func checkConfig(cfg *config.Config) (err *data.CodeError) {
	// host
	configHostCount := 0
	if len(cfg.Hosts.Api) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Rs) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Rsf) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Io) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Up) > 0 {
		configHostCount += 1
	}
	if configHostCount != 0 && configHostCount != 5 {
		err = data.NewEmptyError().AppendDesc("hosts: api/rs/rsf/io/up should config all")
	}
	return
}
