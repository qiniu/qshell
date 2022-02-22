package config

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
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

	if from.Tasks != nil {
		if d.LogSetting == nil {
			d.LogSetting = &LogSetting{}
		}
		d.LogSetting.merge(from.LogSetting)
	}

	if from.Tasks != nil {
		if d.Tasks == nil {
			d.Tasks = &Tasks{}
		}
		d.Tasks.merge(from.Tasks)
	}

	d.ThreadCount = utils.GetNotZeroIntIfExist(d.ThreadCount, from.ThreadCount)
	d.FileEncoding = utils.GetNotEmptyStringIfExist(d.FileEncoding, from.FileEncoding)
	d.KeyFile = utils.GetNotEmptyStringIfExist(d.KeyFile, from.KeyFile)
	d.DestDir = utils.GetNotEmptyStringIfExist(d.DestDir, from.DestDir)
	d.Bucket = utils.GetNotEmptyStringIfExist(d.Bucket, from.Bucket)
	d.Prefix = utils.GetNotEmptyStringIfExist(d.Prefix, from.Prefix)
	d.Suffixes = utils.GetNotEmptyStringIfExist(d.Suffixes, from.Suffixes)
	d.IoHost = utils.GetNotEmptyStringIfExist(d.IoHost, from.IoHost)
	d.Public = utils.GetTrueBoolValueIfExist(d.Public, from.Public)
	d.CheckHash = utils.GetTrueBoolValueIfExist(d.CheckHash, from.CheckHash)

	//down from cdn
	d.Referer = utils.GetNotEmptyStringIfExist(d.Referer, from.Referer)
	d.CdnDomain = utils.GetNotEmptyStringIfExist(d.CdnDomain, from.CdnDomain)
	d.UseHttps = utils.GetTrueBoolValueIfExist(d.UseHttps, from.UseHttps)

	// 下载状态保存路径
	d.RecordRoot = utils.GetNotEmptyStringIfExist(d.RecordRoot, from.RecordRoot)

	d.BatchNum = utils.GetNotZeroIntIfExist(d.BatchNum, from.BatchNum)

	if from.Retry != nil {
		if d.Retry == nil {
			d.Retry = &Retry{}
		}
		d.Retry.merge(from.Retry)
	}
}

func (d *Download) JobId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s", d.DestDir, d.Bucket, d.KeyFile))
}

func (d *Download) Check() error {
	if len(d.Bucket) == 0 && len(d.KeyFile) == 0 {
		return alert.CannotEmptyError("bucket", "")
	}

	if len(d.Bucket) == 0 && len(d.DownloadDomain()) == 0 {
		return alert.Error("bucket / io_host / cdn_domain one them should has value)", "")
	}

	if d.BatchNum <= 0 {
		d.BatchNum = 1000
	}
	return d.LogSetting.Check()
}
