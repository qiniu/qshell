package config

import (
	"encoding/json"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type Config struct {
	Credentials *auth.Credentials `json:"-"`
	UseHttps    string            `json:"use_https,omitempty"`
	Hosts       *Hosts            `json:"hosts,omitempty"`
	Up          *Up               `json:"up,omitempty"`
	Download    *Download         `json:"download,omitempty"`
	FileLog     *log.Config       `json:"log"`
}

func (c *Config) IsUseHttps() bool {
	return c.UseHttps == data.FalseString
}

func (c *Config) HasCredentials() bool {
	return len(c.Credentials.AccessKey) > 0 && c.Credentials.SecretKey != nil
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

func (c *Config) String() string {
	if desc, err := json.MarshalIndent(c, "", "\t"); err == nil {
		return string(desc)
	} else {
		return ""
	}
}
