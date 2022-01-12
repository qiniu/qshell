package config

import (
	"encoding/json"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type Config struct {
	Credentials *auth.Credentials `json:"-"`
	UseHttps    string            `json:"use_https,omitempty"`
	Hosts       *Hosts            `json:"hosts,omitempty"`
	Up          *Up               `json:"up,omitempty"`
	Download    *Download         `json:"download,omitempty"`
}

func (c *Config) IsUseHttps() bool {
	return c.UseHttps == data.FalseString
}

func (c *Config) HasCredentials() bool {
	return c.Credentials != nil && len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
}

func (c *Config) GetRegion() *storage.Region {
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

func (c *Config) Merge(from *Config) {
	if from == nil {
		return
	}

	if !c.HasCredentials() {
		c.Credentials = from.Credentials
	}

	if len(c.UseHttps) == 0 {
		c.UseHttps = from.UseHttps
	}

	if c.Hosts == nil {
		c.Hosts = from.Hosts
	} else {
		c.Hosts.merge(from.Hosts)
	}

	if c.Up == nil {
		c.Up = from.Up
	} else {
		c.Up.merge(from.Up)
	}

	if c.Download == nil {
		c.Download = from.Download
	} else {
		c.Download.merge(from.Download)
	}
}

func (c *Config) GetLogConfig() *LogSetting {
	if c.Up.LogSetting != nil {
		return c.Up.LogSetting
	}

	if c.Download.LogSetting != nil {
		return c.Up.LogSetting
	}
	
	return nil
}

func (c *Config) String() string {
	if desc, err := json.MarshalIndent(c, "", "\t"); err == nil {
		return string(desc)
	} else {
		return ""
	}
}
