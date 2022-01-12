package config

import (
	"strings"
)

type Download struct {
	*LogSetting

	ThreadCount  int    `json:"thread_count,omitempty"`
	FileEncoding string `json:"file_encoding,omitempty"`
	KeyFile      string `json:"key_file,omitempty"`
	DestDir      string `json:"dest_dir,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	Suffixes     string `json:"suffixes,omitempty"`
	IoHost       string `json:"io_host,omitempty"`
	Public       bool   `json:"public,omitempty"`
	CheckHash    bool   `json:"check_hash,omitempty"`

	//down from cdn
	Referer   string `json:"referer,omitempty"`
	CdnDomain string `json:"cdn_domain,omitempty"`
	UseHttps  bool   `json:"use_https,omitempty"`

	// 下载状态保存路径
	RecordRoot string `json:"record_root,omitempty"`

	BatchNum int `json:"-"`

	Tasks *Tasks `json:"tasks,omitempty"`
	Retry *Retry `json:"retry,omitempty"`
}

func (d *Download) Init() {
	if d.BatchNum <= 0 {
		d.BatchNum = 1000
	}
}

// DownloadDomain 获取一个存储空间的下载域名， 默认使用用户配置的域名，如果没有就使用接口随机选择一个下载域名
func (d *Download) DownloadDomain() (domain string) {
	if d.CdnDomain != "" {
		domain = d.CdnDomain
	} else if d.IoHost != "" {
		domain = d.IoHost
	}
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	return
}

func (d *Download) merge(from *Download) {
	if from == nil {
		return
	}
	d.LogSetting.merge(from.LogSetting)
	d.Tasks.merge(from.Tasks)
	d.Retry.merge(from.Retry)
}
