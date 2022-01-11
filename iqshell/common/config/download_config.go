package config

import "strings"

// DownloadConfig qdownload子命令用到的配置参数
type DownloadConfig struct {
	ThreadCount  int    `json:"thread_count"`
	FileEncoding string `json:"file_encoding"`
	KeyFile      string `json:"key_file"`
	DestDir      string `json:"dest_dir"`
	Bucket       string `json:"bucket"`
	Prefix       string `json:"prefix,omitempty"`
	Suffixes     string `json:"suffixes,omitempty"`
	IoHost       string `json:"io_host,omitempty"`
	Public       bool   `json:"public,omitempty"`
	CheckHash    bool   `json:"check_hash"`

	//down from cdn
	Referer   string `json:"referer,omitempty"`
	CdnDomain string `json:"cdn_domain,omitempty"`
	UseHttps  bool   `json:"use_https,omitempty"`

	//log settings
	RecordRoot string `json:"record_root,omitempty"`
	LogLevel   string `json:"log_level,omitempty"`
	LogFile    string `json:"log_file,omitempty"`
	LogRotate  int    `json:"log_rotate,omitempty"`
	LogStdout  bool   `json:"log_stdout,omitempty"`

	BatchNum int
}

func (d *DownloadConfig) Init() {
	if d.BatchNum <= 0 {
		d.BatchNum = 1000
	}
}

// 获取一个存储空间的下载域名， 默认使用用户配置的域名，如果没有就使用接口随机选择一个下载域名
func (d *DownloadConfig) DownloadDomain() (domain string) {
	if d.CdnDomain != "" {
		domain = d.CdnDomain
	} else if d.IoHost != "" {
		domain = d.IoHost
	}
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	return
}
