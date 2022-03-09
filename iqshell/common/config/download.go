package config

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
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

	// 是否使用 getfile api，私有云使用
	GetFileApi *data.Bool `json:"get_file_api"`

	// 下载状态保存路径
	RecordRoot *data.String `json:"record_root,omitempty"`

	BatchNum *data.Int `json:"-"`

	Tasks *Tasks `json:"tasks,omitempty"`
	Retry *Retry `json:"retry,omitempty"`
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

	if data.Empty(d.Bucket) && data.Empty(d.IoHost) && data.Empty(d.CdnDomain) {
		return alert.Error("bucket / io_host / cdn_domain one them should has value)", "")
	}

	if data.Empty(d.BatchNum) {
		d.BatchNum = data.NewInt(1000)
	}
	return d.LogSetting.Check()
}
