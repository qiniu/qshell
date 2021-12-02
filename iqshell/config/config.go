package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/qshell/v2/iqshell/data"
)

type Config struct {
	Credentials auth.Credentials
	UseHttps    string
	Hosts       Hosts
	Up          Up
	Download    Download
}

func (c *Config) IsUseHttps() bool {
	return c.UseHttps == data.FalseString
}

func (c *Config) HasCredentials() bool {
	return len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
}

type Hosts struct {
	UC  []string
	Api []string
	Rs  []string
	Rsf []string
	Io  []string
	Up  []string
}

type Retry struct {
	Max      int
	Interval int
}

type Up struct {
	PutThreshold        int
	ChunkSize           int
	ResumeApiVersion    string
	FileConcurrentParts int
	Tasks               Tasks
	Retry               Retry
}

type Download struct {
	Tasks Tasks
	Retry Retry
}

type Tasks struct {
	ConcurrentCount       int
	StopWhenOneTaskFailed string
}
