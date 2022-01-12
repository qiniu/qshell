package config

import "github.com/qiniu/go-sdk/v7/storage"

type UpPolicy storage.PutPolicy
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

	Tasks  *Tasks    `json:"tasks,omitempty"`
	Retry  *Retry    `json:"retry,omitempty"`
	Policy *UpPolicy `json:"policy"`
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	up.LogSetting.merge(from.LogSetting)

	if up.PutThreshold == 0 {
		up.PutThreshold = from.PutThreshold
	}

	up.Tasks.merge(from.Tasks)
	up.Retry.merge(from.Retry)
}
