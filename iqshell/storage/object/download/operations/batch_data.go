package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type DownloadCfg struct {
	FileEncoding           string `json:"file_encoding,omitempty"`
	KeyFile                string `json:"key_file,omitempty"`
	DestDir                string `json:"dest_dir,omitempty"`
	Bucket                 string `json:"bucket,omitempty"`
	Prefix                 string `json:"prefix,omitempty"`
	SavePathHandler        string `json:"save_path_handler"`
	Suffixes               string `json:"suffixes,omitempty"`
	IoHost                 string `json:"io_host,omitempty"`
	Public                 bool   `json:"public,omitempty"`
	CheckHash              bool   `json:"check_hash,omitempty"`
	EnableSlice            bool   `json:"enable_slice"`
	SliceFileSizeThreshold int64  `json:"slice_file_size_threshold"`
	SliceSize              int64  `json:"slice_size"`
	SliceConcurrentCount   int    `json:"slice_concurrent_count"`

	//down from cdn
	Referer   string `json:"referer,omitempty"`
	CdnDomain string `json:"cdn_domain,omitempty"`

	// 是否使用 getfile api，私有云使用
	GetFileApi bool `json:"get_file_api"`

	// 当遇到错误时删除临时文件
	RemoveTempWhileError bool `json:"remove_temp_while_error"`

	// 下载状态保存路径
	RecordRoot string `json:"record_root,omitempty"`
}

func DefaultDownloadCfg() DownloadCfg {
	return DownloadCfg{
		FileEncoding:           "",
		KeyFile:                "",
		DestDir:                "",
		Bucket:                 "",
		Prefix:                 "",
		Suffixes:               "",
		IoHost:                 "",
		Public:                 false,
		CheckHash:              false,
		Referer:                "",
		CdnDomain:              "",
		GetFileApi:             false,
		EnableSlice:            false,
		SliceFileSizeThreshold: 40 * utils.MB,
		SliceSize:              4 * utils.MB,
		SliceConcurrentCount:   10,
		RemoveTempWhileError:   false,
		RecordRoot:             "",
	}
}

func (d *DownloadCfg) JobId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s", d.DestDir, d.Bucket, d.KeyFile))
}

func (d *DownloadCfg) Check() *data.CodeError {
	if len(d.Bucket) == 0 {
		return alert.CannotEmptyError("bucket", "")
	}
	return nil
}
