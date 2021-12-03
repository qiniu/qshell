package workspace

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func defaultConfig() *config.Config {
	return &config.Config{
		Credentials: auth.Credentials{
			AccessKey: "",
			SecretKey: nil,
		},
		UseHttps: data.TrueString,
		Hosts: config.Hosts{
			Rs:  []string{"rs.qiniu.com"},
			Rsf: []string{"rsf.qiniu.com"},
			Api: []string{"api.qiniu.com"},
			UC:  []string{"uc.qbox.me"},
		},
		Up: config.Up{
			PutThreshold:        1024 * 1024 * 4,
			ChunkSize:           1024 * 1024 * 2,
			ResumeApiVersion:    data.ResumeApiV1,
			FileConcurrentParts: 10,
			Tasks: config.Tasks{
				ConcurrentCount:       3,
				StopWhenOneTaskFailed: data.FalseString,
			},
			Retry: config.Retry{
				Max:      1,
				Interval: 1000,
			},
		},
		Download: config.Download{
			Tasks: config.Tasks{
				ConcurrentCount:       3,
				StopWhenOneTaskFailed: data.FalseString,
			},
			Retry: config.Retry{
				Max:      1,
				Interval: 1000,
			},
		},
	}
}
