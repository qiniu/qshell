package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Config struct {
	Credentials auth.Credentials
	UseHttps    string
	Hosts       Hosts
	Up          Up
	Download    Download
}

func (c Config) IsUseHttps() bool {
	return c.UseHttps == data.FalseString
}

func (c Config) HasCredentials() bool {
	return len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
}

func (c Config) GetRegion() *storage.Region {
	return &storage.Region{
		SrcUpHosts: c.Hosts.Up,
		CdnUpHosts: c.Hosts.Up,
		RsHost:     c.Hosts.GetOneRs(),
		RsfHost:    c.Hosts.GetOneRsf(),
		ApiHost:    c.Hosts.GetOneApi(),
		IovipHost:  c.Hosts.GetOneIo(),
	}
}

type Hosts struct {
	UC  []string
	Api []string
	Rs  []string
	Rsf []string
	Io  []string
	Up  []string
}

func (h Hosts) GetOneUc() string {
	return getOneHostFromStringArray(h.UC)
}

func (h Hosts) GetOneApi() string {
	return getOneHostFromStringArray(h.Api)
}

func (h Hosts) GetOneRs() string {
	return getOneHostFromStringArray(h.Rs)
}

func (h Hosts) GetOneRsf() string {
	return getOneHostFromStringArray(h.Rsf)
}

func (h Hosts) GetOneIo() string {
	return getOneHostFromStringArray(h.Io)
}

func (h Hosts) GetOneUp() string {
	return getOneHostFromStringArray(h.Up)
}

func getOneHostFromStringArray(hosts []string) string {
	if len(hosts) > 0 {
		return hosts[0]
	} else {
		return ""
	}
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
