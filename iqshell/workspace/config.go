package workspace

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/config"
	"github.com/qiniu/qshell/v2/iqshell/data"
)

func DefaultConfig() *config.Config {
	return &config.Config{
		Credentials: auth.Credentials{
			AccessKey: "",
			SecretKey: nil,
		},
		UseHttps:    data.TrueString,
		Hosts: config.Hosts{
			UC:  "",
			Api: "",
			Rs:  "",
			Rsf: "",
			Io:  "",
			Up:  "",
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
