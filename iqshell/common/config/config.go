package config

import (
	"encoding/json"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type Config struct {
	CmdId       string            `json:"-"` // 命令 Id
	Credentials *auth.Credentials `json:"-"`
	UseHttps    *data.Bool         `json:"use_https,omitempty"`
	Hosts       *Hosts            `json:"hosts,omitempty"`
	Up          *Up               `json:"up,omitempty"`
	Download    *Download         `json:"download,omitempty"`
}

func (c *Config) IsUseHttps() bool {
	if c.UseHttps == nil {
		return false
	}
	return c.UseHttps.Value()
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

	c.CmdId = utils.GetNotEmptyStringIfExist(c.CmdId, from.CmdId)
	c.UseHttps = data.GetNotEmptyBoolIfExist(c.UseHttps, from.UseHttps)
	if !c.HasCredentials() {
		c.Credentials = from.Credentials
	}

	if from.Hosts != nil {
		if c.Hosts == nil {
			c.Hosts = &Hosts{}
		}
		c.Hosts.merge(from.Hosts)
	}

	if from.Up != nil {
		if c.Up == nil {
			c.Up = &Up{}
		}
		c.Up.merge(from.Up)
	}

	if from.Download != nil {
		if c.Download == nil {
			c.Download = &Download{}
		}
		c.Download.merge(from.Download)
	}
}

func (c *Config) String() string {
	if desc, err := json.MarshalIndent(c, "", "\t"); err == nil {
		return string(desc)
	} else {
		return ""
	}
}
