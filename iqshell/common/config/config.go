package config

import (
	"encoding/json"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Config struct {
	Credentials auth.Credentials `json:"-"`
	UseHttps    string           `json:"use_https,omitempty"`
	Hosts       Hosts            `json:"hosts,omitempty"`
	Up          Up               `json:"up,omitempty"`
	Download    Download         `json:"download,omitempty"`
}

func (c Config) IsUseHttps() bool {
	return c.UseHttps == data.FalseString
}

func (c Config) HasCredentials() bool {
	return len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
}

func (c Config) GetRegion() *storage.Region {
	if len(c.Hosts.Api) == 0 && len(c.Hosts.Rs) == 0 && len(c.Hosts.Rsf) == 0 &&
		len(c.Hosts.Io) == 0 && len(c.Hosts.Up) == 0 {
		return nil
	}

	return &storage.Region{
		SrcUpHosts: c.Hosts.Up,
		CdnUpHosts: c.Hosts.Up,
		RsHost:     c.Hosts.GetOneRs(),
		RsfHost:    c.Hosts.GetOneRsf(),
		ApiHost:    c.Hosts.GetOneApi(),
		IovipHost:  c.Hosts.GetOneIo(),
	}
}

func (c Config) String() string {
	if desc, err := json.MarshalIndent(c, "", "\t"); err == nil {
		return string(desc)
	} else {
		return ""
	}
}

type Hosts struct {
	UC  []string `json:"uc,omitempty"`
	Api []string `json:"api,omitempty"`
	Rs  []string `json:"rs,omitempty"`
	Rsf []string `json:"rsf,omitempty"`
	Io  []string `json:"io,omitempty"`
	Up  []string `json:"up,omitempty"`
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
	Max      int `json:"max,omitempty"`
	Interval int `json:"interval,omitempty"`
}

type Up struct {
	PutThreshold        int    `json:"put_threshold,omitempty"`
	ChunkSize           int    `json:"chunk_size,omitempty"`
	ResumeApiVersion    string `json:"resume_api_version,omitempty"`
	FileConcurrentParts int    `json:"file_concurrent_parts"`
	Tasks               Tasks  `json:"tasks,omitempty"`
	Retry               Retry  `json:"retry,omitempty"`
}

type Download struct {
	Tasks Tasks `json:"tasks,omitempty"`
	Retry Retry `json:"retry,omitempty"`
}

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed,omitempty"`
}
