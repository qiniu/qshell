package config

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
	"path/filepath"
	"strings"
)

type Up struct {
	*LogSetting

	UpHost    *data.String `json:"up_host,omitempty"`
	BindUpIp  *data.String `json:"bind_up_ip,omitempty"`
	BindRsIp  *data.String `json:"bind_rs_ip,omitempty"`
	BindNicIp *data.String `json:"bind_nic_ip,omitempty"` //local network interface card config

	SrcDir                 *data.String `json:"src_dir,omitempty"`
	FileList               *data.String `json:"file_list,omitempty"`
	IgnoreDir              *data.Bool   `json:"ignore_dir,omitempty"`
	SkipFilePrefixes       *data.String `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes       *data.String `json:"skip_path_prefixes,omitempty"`
	SkipFixedStrings       *data.String `json:"skip_fixed_strings,omitempty"`
	SkipSuffixes           *data.String `json:"skip_suffixes,omitempty"`
	FileEncoding           *data.String `json:"file_encoding,omitempty"`
	Bucket                 *data.String `json:"bucket,omitempty"`
	ResumableAPIV2         *data.Bool   `json:"resumable_api_v2,omitempty"`
	ResumableAPIV2PartSize *data.Int64  `json:"resumable_api_v2_part_size,omitempty"`
	PutThreshold           *data.Int64  `json:"put_threshold,omitempty"`
	KeyPrefix              *data.String `json:"key_prefix,omitempty"`
	Overwrite              *data.Bool   `json:"overwrite,omitempty"`
	CheckExists            *data.Bool   `json:"check_exists,omitempty"`
	CheckHash              *data.Bool   `json:"check_hash,omitempty"`
	CheckSize              *data.Bool   `json:"check_size,omitempty"`
	RescanLocal            *data.Bool   `json:"rescan_local,omitempty"`
	FileType               *data.Int    `json:"file_type,omitempty"`
	DeleteOnSuccess        *data.Bool   `json:"delete_on_success,omitempty"`
	DisableResume          *data.Bool   `json:"disable_resume,omitempty"`
	DisableForm            *data.Bool   `json:"disable_form,omitempty"`
	WorkerCount            *data.Int    `json:"work_count,omitempty"` // 分片上传并发数
	RecordRoot             *data.String `json:"record_root,omitempty"`

	Tasks  *Tasks             `json:"-"`
	Retry  *Retry             `json:"-"`
	Policy *storage.PutPolicy `json:"policy"`
}

func (up *Up) IsIgnoreDir() bool {
	if up.IgnoreDir == nil {
		return false
	}
	return up.IgnoreDir.Value()
}

func (up *Up) IsResumeAPIV2() bool {
	if up.ResumableAPIV2 == nil {
		return false
	}
	return up.ResumableAPIV2.Value()
}

func (up *Up) IsOverwrite() bool {
	if up.Overwrite == nil {
		return false
	}
	return up.Overwrite.Value()
}

func (up *Up) IsCheckExists() bool {
	if up.CheckExists == nil {
		return false
	}
	return up.CheckExists.Value()
}

func (up *Up) IsCheckHash() bool {
	if up.CheckHash == nil {
		return false
	}
	return up.CheckHash.Value()
}

func (up *Up) IsCheckSize() bool {
	if up.CheckSize == nil {
		return false
	}
	return up.CheckSize.Value()
}

func (up *Up) IsRescanLocal() bool {
	if up.RescanLocal == nil {
		return false
	}
	return up.RescanLocal.Value()
}

func (up *Up) IsDeleteOnSuccess() bool {
	if up.DeleteOnSuccess == nil {
		return false
	}
	return up.DeleteOnSuccess.Value()
}

func (up *Up) IsDisableResume() bool {
	if up.DisableResume == nil {
		return false
	}
	return up.DisableResume.Value()
}

func (up *Up) IsDisableForm() bool {
	if up.DisableForm == nil {
		return false
	}
	return up.DisableForm.Value()
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	up.SrcDir = data.GetNotEmptyStringIfExist(up.SrcDir, from.SrcDir)
	up.FileList = data.GetNotEmptyStringIfExist(up.FileList, from.FileList)
	up.IgnoreDir = data.GetNotEmptyBoolIfExist(up.IgnoreDir, from.IgnoreDir)
	up.SkipFilePrefixes = data.GetNotEmptyStringIfExist(up.SkipFilePrefixes, from.SkipFilePrefixes)
	up.SkipPathPrefixes = data.GetNotEmptyStringIfExist(up.SkipPathPrefixes, from.SkipPathPrefixes)
	up.SkipFixedStrings = data.GetNotEmptyStringIfExist(up.SkipFixedStrings, from.SkipFixedStrings)
	up.SkipSuffixes = data.GetNotEmptyStringIfExist(up.SkipSuffixes, from.SkipSuffixes)

	up.UpHost = data.GetNotEmptyStringIfExist(up.UpHost, from.UpHost)
	up.BindUpIp = data.GetNotEmptyStringIfExist(up.BindUpIp, from.BindUpIp)
	up.BindRsIp = data.GetNotEmptyStringIfExist(up.BindRsIp, from.BindRsIp)
	up.BindNicIp = data.GetNotEmptyStringIfExist(up.BindNicIp, from.BindNicIp)

	up.FileEncoding = data.GetNotEmptyStringIfExist(up.FileEncoding, from.FileEncoding)
	up.Bucket = data.GetNotEmptyStringIfExist(up.Bucket, from.Bucket)
	up.ResumableAPIV2 = data.GetNotEmptyBoolIfExist(up.ResumableAPIV2, from.ResumableAPIV2)
	up.ResumableAPIV2PartSize = data.GetNotEmptyInt64IfExist(up.ResumableAPIV2PartSize, from.ResumableAPIV2PartSize)
	up.PutThreshold = data.GetNotEmptyInt64IfExist(up.PutThreshold, from.PutThreshold)
	up.KeyPrefix = data.GetNotEmptyStringIfExist(up.KeyPrefix, from.KeyPrefix)
	up.Overwrite = data.GetNotEmptyBoolIfExist(up.Overwrite, from.Overwrite)
	up.CheckExists = data.GetNotEmptyBoolIfExist(up.CheckExists, from.CheckExists)
	up.CheckHash = data.GetNotEmptyBoolIfExist(up.CheckHash, from.CheckHash)
	up.CheckSize = data.GetNotEmptyBoolIfExist(up.CheckSize, from.CheckSize)
	up.RescanLocal = data.GetNotEmptyBoolIfExist(up.RescanLocal, from.RescanLocal)
	up.FileType = data.GetNotEmptyIntIfExist(up.FileType, from.FileType)
	up.DeleteOnSuccess = data.GetNotEmptyBoolIfExist(up.DeleteOnSuccess, from.DeleteOnSuccess)
	up.DisableResume = data.GetNotEmptyBoolIfExist(up.DisableResume, from.DisableResume)
	up.DisableForm = data.GetNotEmptyBoolIfExist(up.DisableForm, from.DisableForm)
	up.WorkerCount = data.GetNotEmptyIntIfExist(up.WorkerCount, from.WorkerCount)
	up.RecordRoot = data.GetNotEmptyStringIfExist(up.RecordRoot, from.RecordRoot)

	if from.LogSetting != nil {
		if up.LogSetting == nil {
			up.LogSetting = &LogSetting{}
		}
		up.LogSetting.merge(from.LogSetting)
	}

	if from.Policy != nil {
		if up.Policy == nil {
			up.Policy = &storage.PutPolicy{}
		}
		mergeUploadPolicy(from.Policy, up.Policy)
	}

	if from.Tasks != nil {
		if up.Tasks == nil {
			up.Tasks = &Tasks{}
		}
		up.Tasks.merge(from.Tasks)
	}

	if from.Retry != nil {
		if up.Retry == nil {
			up.Retry = &Retry{}
		}
		up.Retry.merge(from.Retry)
	}
}

func (up *Up) JobId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s", up.SrcDir.Value(), up.Bucket.Value(), up.FileList.Value()))
}

func (up *Up) GetLogLevel() int {
	if up.LogLevel == nil {
		return log.LevelInfo
	}

	//init log level
	logLevel := log.LevelInfo
	switch up.LogLevel.Value() {
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarning
	case "error":
		logLevel = log.LevelError
	default:
		logLevel = log.LevelInfo
	}
	return int(logLevel)
}

func (up *Up) GetLogRotate() int {
	logRotate := 1
	if data.NotEmpty(up.LogRotate) {
		logRotate = up.LogRotate.Value()
	}
	return logRotate
}

func (up *Up) Check() error {
	// 验证大小
	if up.ResumableAPIV2PartSize == nil {
		up.ResumableAPIV2PartSize = data.NewInt64(data.BLOCK_SIZE)
	} else if up.ResumableAPIV2PartSize.Value() < int64(utils.MB) {
		up.ResumableAPIV2PartSize = data.NewInt64(utils.MB)
	} else if up.ResumableAPIV2PartSize.Value() > int64(utils.GB) {
		up.ResumableAPIV2PartSize = data.NewInt64(utils.GB)
	}

	if data.Empty(up.Bucket) {
		return alert.CannotEmptyError("Bucket", "")
	}

	if data.Empty(up.SrcDir) {
		return alert.CannotEmptyError("SrcDir", "")
	}

	srcFileInfo, err := os.Stat(up.SrcDir.Value())
	if err != nil {
		return errors.New("invalid SrcDir:" + err.Error())
	}

	if !srcFileInfo.IsDir() {
		return fmt.Errorf("SrcDir should be a directory: %s", *up.SrcDir)
	}

	if data.NotEmpty(up.FileList) {
		fileListInfo, err := os.Stat(up.FileList.Value())
		if err != nil {
			return fmt.Errorf("invalid FileList:%s error:%v", up.FileList.Value(), err.Error())
		}

		if fileListInfo.IsDir() {
			return fmt.Errorf("FileList should be a file: %s", up.FileList.Value())
		}
	}

	if up.FileType.Value() != 1 && up.FileType.Value() != 0 {
		return errors.New("wrong Filetype, It should be one of 1, 2, 3")
	}

	if up.Policy != nil {
		if (up.Policy.CallbackURL == "" && up.Policy.CallbackHost != "") ||
			(up.Policy.CallbackURL != "" && up.Policy.CallbackHost == "") {
			return errors.New("callbackUrls and callback must exist at the same time")
		}

		if up.Policy.CallbackHost != "" && up.Policy.CallbackURL != "" {
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

func (up *Up) HitByPathPrefixes(localFileRelativePath string) (hit bool, pathPrefix string) {

	if data.NotEmpty(up.SkipPathPrefixes) {
		//unpack skip prefix
		pathPrefixes := strings.Split(up.SkipPathPrefixes.Value(), ",")
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

func (up *Up) HitByFilePrefixes(localFileRelativePath string) (hit bool, filePrefix string) {
	if data.NotEmpty(up.SkipFilePrefixes) {
		//unpack skip prefix
		filePrefixes := strings.Split(up.SkipFilePrefixes.Value(), ",")
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

func (up *Up) HitByFixesString(localFileRelativePath string) (hit bool, hitFixedStr string) {
	if data.NotEmpty(up.SkipFixedStrings) {
		//unpack fixed strings
		fixedStrings := strings.Split(up.SkipFixedStrings.Value(), ",")
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

func (up *Up) HitBySuffixes(localFileRelativePath string) (hit bool, hitSuffix string) {
	if data.NotEmpty(up.SkipSuffixes) {
		suffixes := strings.Split(up.SkipSuffixes.Value(), ",")
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
