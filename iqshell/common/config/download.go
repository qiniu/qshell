package config

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"strings"
)

type Download struct {
	*LogSetting

	ThreadCount  *data.Int    `json:"thread_count,omitempty"`
	FileEncoding *data.String `json:"file_encoding,omitempty"`
	KeyFile      *data.String `json:"key_file,omitempty"`
	DestDir      *data.String `json:"dest_dir,omitempty"`
	Bucket       *data.String `json:"bucket,omitempty"`
	Prefix       *data.String `json:"prefix,omitempty"`
	Suffixes     *data.String `json:"suffixes,omitempty"`
	IoHost       *data.String `json:"io_host,omitempty"`
	Public       *data.Bool   `json:"public,omitempty"`
	CheckHash    *data.Bool   `json:"check_hash,omitempty"`

	//down from cdn
	Referer   *data.String `json:"referer,omitempty"`
	CdnDomain *data.String `json:"cdn_domain,omitempty"`
	UseHttps  *data.Bool   `json:"use_https,omitempty"`

	// 下载状态保存路径
	RecordRoot *data.String `json:"record_root,omitempty"`

	BatchNum *data.Int `json:"-"`

	Tasks *Tasks `json:"tasks,omitempty"`
	Retry *Retry `json:"retry,omitempty"`
}

// DownloadDomain 获取一个存储空间的下载域名， 默认使用用户配置的域名，如果没有就使用接口随机选择一个下载域名
func (d *Download) DownloadDomain() (domain string) {
	if data.NotEmpty(d.CdnDomain) {
		domain = d.CdnDomain.Value()
	} else if data.NotEmpty(d.IoHost) {
		domain = d.IoHost.Value()
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

	d.ThreadCount = data.GetNotEmptyIntIfExist(d.ThreadCount, from.ThreadCount)
	d.FileEncoding = data.GetNotEmptyStringIfExist(d.FileEncoding, from.FileEncoding)
	d.KeyFile = data.GetNotEmptyStringIfExist(d.KeyFile, from.KeyFile)
	d.DestDir = data.GetNotEmptyStringIfExist(d.DestDir, from.DestDir)
	d.Bucket = data.GetNotEmptyStringIfExist(d.Bucket, from.Bucket)
	d.Prefix = data.GetNotEmptyStringIfExist(d.Prefix, from.Prefix)
	d.Suffixes = data.GetNotEmptyStringIfExist(d.Suffixes, from.Suffixes)
	d.IoHost = data.GetNotEmptyStringIfExist(d.IoHost, from.IoHost)
	d.Public = data.GetNotEmptyBoolIfExist(d.Public, from.Public)
	d.CheckHash = data.GetNotEmptyBoolIfExist(d.CheckHash, from.CheckHash)

	//down from cdn
	d.Referer = data.GetNotEmptyStringIfExist(d.Referer, from.Referer)
	d.CdnDomain = data.GetNotEmptyStringIfExist(d.CdnDomain, from.CdnDomain)
	d.UseHttps = data.GetNotEmptyBoolIfExist(d.UseHttps, from.UseHttps)

	// 下载状态保存路径
	d.RecordRoot = data.GetNotEmptyStringIfExist(d.RecordRoot, from.RecordRoot)

	d.BatchNum = data.GetNotEmptyIntIfExist(d.BatchNum, from.BatchNum)

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
	if data.Empty(d.Bucket) && data.Empty(d.KeyFile) {
		return alert.CannotEmptyError("bucket", "")
	}

	if data.Empty(d.Bucket) && len(d.DownloadDomain()) == 0 {
		return alert.Error("bucket / io_host / cdn_domain one them should has value)", "")
	}

	if data.Empty(d.BatchNum) {
		d.BatchNum = data.NewInt(1000)
	}
	return d.LogSetting.Check()
}
