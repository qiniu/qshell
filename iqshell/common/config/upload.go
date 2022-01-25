package config

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"os"
	"path/filepath"
	"strings"
)

type Up struct {
	*LogSetting

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
	RecordRoot             string `json:"record_root"`

	Tasks  *Tasks             `json:"tasks,omitempty"`
	Retry  *Retry             `json:"retry,omitempty"`
	Policy *storage.PutPolicy `json:"policy"`
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	up.LogSetting.merge(from.LogSetting)

	if up.PutThreshold == 0 {
		up.PutThreshold = from.PutThreshold
	}

	if up.Tasks == nil {
		up.Tasks = from.Tasks
	} else {
		up.Tasks.merge(from.Tasks)
	}

	if up.Retry == nil {
		up.Retry = from.Retry
	} else {
		up.Retry.merge(from.Retry)
	}
}

func (up *Up) Check() error {
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
