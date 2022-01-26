package config

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
	"path/filepath"
	"strings"
)

type Up struct {
	*LogSetting

	SrcDir           string `json:"-"`
	FileList         string `json:"-"`
	IgnoreDir        string `json:"-"`
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
	ResumableAPIV2         string   `json:"resumable_api_v2,omitempty"`
	ResumableAPIV2PartSize int64  `json:"resumable_api_v2_part_size,omitempty"`
	PutThreshold           int64  `json:"put_threshold,omitempty"`
	KeyPrefix              string `json:"key_prefix,omitempty"`
	Overwrite              string   `json:"overwrite,omitempty"`
	CheckExists            string   `json:"check_exists,omitempty"`
	CheckHash              string   `json:"check_hash,omitempty"`
	CheckSize              string   `json:"check_size,omitempty"`
	RescanLocal            string   `json:"rescan_local,omitempty"`
	FileType               int    `json:"file_type,omitempty"`
	DeleteOnSuccess        string   `json:"delete_on_success,omitempty"`
	DisableResume          string   `json:"disable_resume,omitempty"`
	DisableForm            string   `json:"disable_form"`
	WorkerCount            int    `json:"work_count"` // 分片上传并发数
	RecordRoot             string `json:"record_root"`

	Tasks  *Tasks             `json:"tasks,omitempty"`
	Retry  *Retry             `json:"retry,omitempty"`
	Policy *storage.PutPolicy `json:"policy"`
}

func (up *Up)IsIgnoreDir() bool {
	return up.IgnoreDir == data.TrueString
}

func (up *Up)IsResumableAPIV2() bool {
	return up.ResumableAPIV2 == data.TrueString
}

func (up *Up)IsOverwrite() bool {
	return up.Overwrite == data.TrueString
}

func (up *Up)IsCheckExists() bool {
	return up.CheckExists == data.TrueString
}

func (up *Up)IsCheckHash() bool {
	return up.CheckHash == data.TrueString
}

func (up *Up)IsCheckSize() bool {
	return up.CheckSize == data.TrueString
}

func (up *Up)IsRescanLocal() bool {
	return up.RescanLocal == data.TrueString
}

func (up *Up)IsDeleteOnSuccess() bool {
	return up.DeleteOnSuccess == data.TrueString
}

func (up *Up)IsDisableResume() bool {
	return up.DisableResume == data.TrueString
}

func (up *Up)IsDisableForm() bool {
	return up.DisableForm == data.TrueString
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	up.SrcDir = utils.GetNotEmptyStringIfExist(up.SrcDir, from.SrcDir)
	up.FileList = utils.GetNotEmptyStringIfExist(up.FileList, from.FileList)
	up.IgnoreDir = utils.GetNotEmptyStringIfExist(up.IgnoreDir, from.IgnoreDir)
	up.SkipFilePrefixes = utils.GetNotEmptyStringIfExist(up.SkipFilePrefixes, from.SkipFilePrefixes)
	up.SkipPathPrefixes = utils.GetNotEmptyStringIfExist(up.SkipPathPrefixes, from.SkipPathPrefixes)
	up.SkipFixedStrings = utils.GetNotEmptyStringIfExist(up.SkipFixedStrings, from.SkipFixedStrings)
	up.SkipSuffixes = utils.GetNotEmptyStringIfExist(up.SkipSuffixes, from.SkipSuffixes)

	up.UpHost = utils.GetNotEmptyStringIfExist(up.UpHost, from.UpHost)
	up.BindUpIp = utils.GetNotEmptyStringIfExist(up.BindUpIp, from.BindUpIp)
	up.BindRsIp = utils.GetNotEmptyStringIfExist(up.BindRsIp, from.BindRsIp)
	up.BindNicIp = utils.GetNotEmptyStringIfExist(up.BindNicIp, from.BindNicIp)

	up.PutThreshold = utils.GetNotZeroInt64IfExist(up.PutThreshold, from.PutThreshold)

	if from.Tasks != nil {
		if up.LogSetting == nil {
			up.LogSetting = &LogSetting{}
		}
		up.LogSetting.merge(from.LogSetting)
	}

	if up.PutThreshold == 0 {
		up.PutThreshold = from.PutThreshold
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
	return utils.Md5Hex(fmt.Sprintf("%s:%s", up.SrcDir, up.Bucket))
}

func (up *Up) GetLogLevel() int {

	//init log level
	logLevel := log.LevelInfo
	switch up.LogLevel {
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
	if up.LogRotate > 0 {
		logRotate = up.LogRotate
	}
	return logRotate
}

func (up *Up) Check() error {
	// 验证大小
	if up.ResumableAPIV2PartSize <= 0 {
		up.ResumableAPIV2PartSize = data.BLOCK_SIZE
	} else if up.ResumableAPIV2PartSize < int64(utils.MB) {
		up.ResumableAPIV2PartSize = int64(utils.MB)
	} else if up.ResumableAPIV2PartSize > int64(utils.GB) {
		up.ResumableAPIV2PartSize = int64(utils.GB)
	}

	if up.FileType != 1 && up.FileType != 0 {
		return errors.New("wrong Filetype, It should be 0 or 1")
	}

	srcFileInfo, err := os.Stat(up.SrcDir)
	if err != nil {
		return errors.New("upload config error for parameter `SrcDir`:" + err.Error())
	}

	if !srcFileInfo.IsDir() {
		return errors.New("upload src dir should be a directory")
	}

	if up.Bucket == "" {
		return errors.New("upload config no `bucket` specified")
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

	if up.SkipPathPrefixes != "" {
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

func (up *Up) HitByFilePrefixes(localFileRelativePath string) (hit bool, filePrefix string) {
	if up.SkipFilePrefixes != "" {
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

func (up *Up) HitByFixesString(localFileRelativePath string) (hit bool, hitFixedStr string) {
	if up.SkipFixedStrings != "" {
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

func (up *Up) HitBySuffixes(localFileRelativePath string) (hit bool, hitSuffix string) {
	if up.SkipSuffixes != "" {
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
