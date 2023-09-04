package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type UploadConfig struct {
	UpHost    string `json:"up_host,omitempty"`
	BindUpIp  string `json:"bind_up_ip,omitempty"`
	BindRsIp  string `json:"bind_rs_ip,omitempty"`
	BindNicIp string `json:"bind_nic_ip,omitempty"` //local network interface card config

	SrcDir                 string `json:"src_dir,omitempty"`
	FileList               string `json:"file_list,omitempty"`
	IgnoreDir              bool   `json:"ignore_dir,omitempty"`
	SkipFilePrefixes       string `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes       string `json:"skip_path_prefixes,omitempty"`
	SkipFixedStrings       string `json:"skip_fixed_strings,omitempty"`
	SkipSuffixes           string `json:"skip_suffixes,omitempty"`
	FileEncoding           string `json:"file_encoding,omitempty"`
	Bucket                 string `json:"bucket,omitempty"`
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
	DisableForm            bool   `json:"disable_form,omitempty"`
	WorkerCount            int    `json:"work_count,omitempty"` // 分片上传并发数
	RecordRoot             string `json:"record_root,omitempty"`
	SequentialReadFile     bool   `json:"sequential_read_file"` // 文件顺序读

	Policy *storage.PutPolicy `json:"policy"`
}

func DefaultUploadConfig() UploadConfig {
	return UploadConfig{
		UpHost:                 "",
		BindUpIp:               "",
		BindRsIp:               "",
		BindNicIp:              "",
		SrcDir:                 "",
		FileList:               "",
		IgnoreDir:              false,
		SkipFilePrefixes:       "",
		SkipPathPrefixes:       "",
		SkipFixedStrings:       "",
		SkipSuffixes:           "",
		FileEncoding:           "",
		Bucket:                 "",
		ResumableAPIV2:         false,
		ResumableAPIV2PartSize: 4 * 1024 * 1024,
		PutThreshold:           8 * 1024 * 1024,
		KeyPrefix:              "",
		Overwrite:              false,
		CheckExists:            false,
		CheckHash:              false,
		CheckSize:              false,
		RescanLocal:            false,
		FileType:               0,
		DeleteOnSuccess:        false,
		DisableResume:          false,
		DisableForm:            false,
		WorkerCount:            3,
		RecordRoot:             "",
		Policy:                 nil,
	}
}

func (up *UploadConfig) IsIgnoreDir() bool {
	return up.IgnoreDir
}

func (up *UploadConfig) IsResumeAPIV2() bool {
	return up.ResumableAPIV2
}

func (up *UploadConfig) IsOverwrite() bool {
	return up.Overwrite
}

func (up *UploadConfig) IsCheckExists() bool {
	return up.CheckExists
}

func (up *UploadConfig) IsCheckHash() bool {
	return up.CheckHash
}

func (up *UploadConfig) IsCheckSize() bool {
	return up.CheckSize
}

func (up *UploadConfig) IsRescanLocal() bool {
	return up.RescanLocal
}

func (up *UploadConfig) IsDeleteOnSuccess() bool {
	return up.DeleteOnSuccess
}

func (up *UploadConfig) IsDisableResume() bool {
	return up.DisableResume
}

func (up *UploadConfig) IsDisableForm() bool {
	return up.DisableForm
}

func (up *UploadConfig) JobId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s", up.SrcDir, up.Bucket, up.FileList))
}

func (up *UploadConfig) Check() *data.CodeError {
	// 验证大小
	if up.ResumableAPIV2PartSize == 0 {
		up.ResumableAPIV2PartSize = data.BLOCK_SIZE
	} else if up.ResumableAPIV2PartSize < int64(utils.MB) {
		up.ResumableAPIV2PartSize = utils.MB
	} else if up.ResumableAPIV2PartSize > int64(utils.GB) {
		up.ResumableAPIV2PartSize = utils.GB
	}

	if len(up.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if len(up.SrcDir) == 0 {
		return alert.CannotEmptyError("SrcDir", "")
	}

	srcFileInfo, err := os.Stat(up.SrcDir)
	if err != nil {
		return data.NewEmptyError().AppendDesc("invalid SrcDir:" + err.Error())
	}

	if !srcFileInfo.IsDir() {
		return data.NewEmptyError().AppendDescF("SrcDir should be a directory: %s", up.SrcDir)
	}

	if len(up.FileList) > 0 {
		fileListInfo, err := os.Stat(up.FileList)
		if err != nil {
			return data.NewEmptyError().AppendDesc("invalid FileList:" + err.Error())
		}

		if fileListInfo.IsDir() {
			return data.NewEmptyError().AppendDescF("FileList should be a file: %s", up.FileList)
		}
	}

	if up.FileType < 0 || up.FileType > 3 {
		return data.NewEmptyError().AppendDesc("wrong Filetype, It should be one of 0, 1, 2, 3")
	}

	if up.Policy != nil {
		//if (up.Policy.CallbackURL == "" && up.Policy.CallbackHost != "") ||
		//	(up.Policy.CallbackURL != "" && up.Policy.CallbackHost == "") {
		//	return data.NewEmptyError().AppendDesc("callbackUrls and callback must exist at the same time")
		//}

		if up.Policy.CallbackURL != "" {
			callbackUrls := strings.Replace(up.Policy.CallbackURL, ",", ";", -1)
			up.Policy.CallbackURL = callbackUrls
			if len(up.Policy.CallbackBody) == 0 {
				up.Policy.CallbackBody = "key=$(key)&hash=$(etag)"
			}
			if len(up.Policy.CallbackBodyType) == 0 {
				up.Policy.CallbackBodyType = "application/x-www-form-urlencoded"
			}
		}
	}

	return nil
}

func (up *UploadConfig) HitByPathPrefixes(localFileRelativePath string) (hit bool, pathPrefix string) {

	if len(up.SkipPathPrefixes) > 0 {
		//unpack skip prefix
		pathPrefixes := strings.Split(up.SkipPathPrefixes, ",")
		for _, prefix := range pathPrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			if strings.HasPrefix(localFileRelativePath, prefix) {
				pathPrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (up *UploadConfig) HitByFilePrefixes(localFileRelativePath string) (hit bool, filePrefix string) {
	if len(up.SkipFilePrefixes) > 0 {
		//unpack skip prefix
		filePrefixes := strings.Split(up.SkipFilePrefixes, ",")
		for _, prefix := range filePrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			localFileName := filepath.Base(localFileRelativePath)
			if strings.HasPrefix(localFileName, prefix) {
				filePrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (up *UploadConfig) HitByFixesString(localFileRelativePath string) (hit bool, hitFixedStr string) {
	if len(up.SkipFixedStrings) > 0 {
		//unpack fixed strings
		fixedStrings := strings.Split(up.SkipFixedStrings, ",")
		for _, fixedStr := range fixedStrings {
			if strings.TrimSpace(fixedStr) == "" {
				continue
			}

			if strings.Contains(localFileRelativePath, fixedStr) {
				hitFixedStr = fixedStr
				hit = true
				break
			}
		}
	}
	return

}

func (up *UploadConfig) HitBySuffixes(localFileRelativePath string) (hit bool, hitSuffix string) {
	if len(up.SkipSuffixes) > 0 {
		suffixes := strings.Split(up.SkipSuffixes, ",")
		for _, suffix := range suffixes {
			if strings.TrimSpace(suffix) == "" {
				continue
			}

			if strings.HasSuffix(localFileRelativePath, suffix) {
				hitSuffix = suffix
				hit = true
				break
			}
		}
	}
	return
}
