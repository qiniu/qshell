package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"strings"
)

const (
	TrueString  = "true"
	FalseString = "false"

	ResumeApiV1 = "v1"
	ResumeApiV2 = "v2"
)

type Config struct {
	Credentials auth.Credentials `json:"-"`
	UseHttps    string            `json:"use_https"`
	Hosts       Hosts             `json:"hosts"`
	Up          Up                `json:"up"`
	Download    Download          `json:"download"`
}

func (c *Config) IsUseHttps() bool {
	return c.UseHttps == TrueString
}

func (c *Config) HasCredentials() bool {
	return len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
}

type Hosts struct {
	UC  string `json:"uc"`
	Api string `json:"api"`
	Rs  string `json:"rs"`
	Rsf string `json:"rsf"`
	Io  string `json:"io"`
	Up  string `json:"up"`
}

type Retry struct {
	Max      int `json:"max"`
	Interval int `json:"interval"`
}

type Up struct {
	PutThreshold        int    `json:"put_threshold"`
	ChunkSize           int    `json:"chunk_size"`
	ResumeApiVersion    string `json:"resume_api_version"`
	FileConcurrentParts int    `json:"file_concurrent_parts"`
	Tasks               Tasks  `json:"tasks"`
	Retry               Retry  `json:"retry"`
}

type Download struct {
	Tasks Tasks `json:"tasks"`
	Retry Retry `json:"retry"`
}

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed"`
}

// 此处 host 可能包含 scheme
func parseHostArray(host string) []string {
	if len(host) == 0 {
		return nil
	}

	if !strings.Contains(host, ",") {
		return []string{host}
	}

	return strings.Split(host, ",")
}

func DefaultConfig() *Config {
	return &Config{
		Credentials: auth.Credentials{
			AccessKey: "",
			SecretKey: nil,
		},
		UseHttps:    TrueString,
		Hosts: Hosts{
			UC:  "",
			Api: "",
			Rs:  "",
			Rsf: "",
			Io:  "",
			Up:  "",
		},
		Up: Up{
			PutThreshold:        1024 * 1024 * 4,
			ChunkSize:           1024 * 1024 * 2,
			ResumeApiVersion:    ResumeApiV1,
			FileConcurrentParts: 10,
			Tasks: Tasks{
				ConcurrentCount:       3,
				StopWhenOneTaskFailed: FalseString,
			},
			Retry: Retry{
				Max:      1,
				Interval: 1000,
			},
		},
		Download: Download{
			Tasks: Tasks{
				ConcurrentCount:       3,
				StopWhenOneTaskFailed: FalseString,
			},
			Retry: Retry{
				Max:      1,
				Interval: 1000,
			},
		},
	}
}
