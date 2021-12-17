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

type UpPolicy storage.PutPolicy
type Up struct {
	SrcDir           string `json:"-"`
	FileList         string `json:"-"`
	IgnoreDir        bool   `json:"-"`
	SkipFilePrefixes string `json:"-"`
	SkipPathPrefixes string `json:"-"`
	SkipFixedStrings string `json:"-"`
	SkipSuffixes     string `json:"-"`

	UpHost    string `json:"up_host,omitempty"`
	BindUpIp  string `json:"bind_up_ip,omitempty"`
	BindRsIp  string `json:"bind_rs_ip,omitempty"`
	BindNicIp string `json:"bind_nic_ip,omitempty"` //local network interface card config

	FileEncoding           string `json:"file_encoding"`
	Bucket                 string `json:"bucket"`
	ResumableAPIV2         bool   `json:"resumable_api_v2,omitempty"`
	ResumableAPIV2PartSize int64  `json:"resumable_api_v2_part_size,omitempty"`
	PutThreshold           int64  `json:"put_threshold,omitempty"`
	KeyPrefix              string `json:"key_prefix,omitempty"`
	Overwrite              bool   `json:"overwrite,omitempty"`
	CheckExists            bool   `json:"check_exists,omitempty"`
	CheckHash              bool   `json:"check_hash,omitempty"`
	CheckSize              bool   `json:"check_size,omitempty"`
	RescanLocal            bool   `json:"rescan_local,omitempty"`
	FileType               int    `json:"file_type,omitempty"`
	DeleteOnSuccess        bool   `json:"delete_on_success,omitempty"`
	DisableResume          bool   `json:"disable_resume,omitempty"`

	//log settings
	LogLevel  string `json:"log_level,omitempty"`
	LogFile   string `json:"log_file,omitempty"`
	LogRotate int    `json:"log_rotate,omitempty"`
	LogStdout bool   `json:"log_stdout,omitempty"`

	Tasks  Tasks    `json:"tasks,omitempty"`
	Retry  Retry    `json:"retry,omitempty"`
	Policy UpPolicy `json:"policy"`
}

type Download struct {
	Tasks Tasks `json:"tasks,omitempty"`
	Retry Retry `json:"retry,omitempty"`
}

type Tasks struct {
	ConcurrentCount       int    `json:"concurrent_count,omitempty"`
	StopWhenOneTaskFailed string `json:"stop_when_one_task_failed,omitempty"`
}
