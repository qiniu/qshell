package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type UploadConfig struct {
	UpHost    string `json:"up_host,omitempty"`
	BindUpIp  string `json:"bind_up_ip,omitempty"`
	BindRsIp  string `json:"bind_rs_ip,omitempty"`
	BindNicIp string `json:"bind_nic_ip,omitempty"` // local network interface card config

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
	SequentialReadFile     bool   `json:"sequential_read_file"`   // 文件顺序读
	Accelerate             bool   `json:"uploading_acceleration"` // 开启上传加速

	// 唯一属主标识。特殊场景下非常有用，例如根据 App-Client 标识给图片或视频打水印。
	EndUser string `json:"end_user,omitempty"`

	// 上传成功后，七牛云向业务服务器发送 POST 请求的 URL。必须是公网上可以正常进行 POST 请求并能响应 HTTP/1.1 200 OK 的有效 URL。
	// 另外，为了给客户端有一致的体验，我们要求 callbackUrl 返回包 Content-Type 为 “application/json”，即返回的内容必须是合法的
	// JSON 文本。出于高可用的考虑，本字段允许设置多个 callbackUrl（用英文符号 ; 分隔），在前一个 callbackUrl 请求失败的时候会依次
	// 重试下一个 callbackUrl。一个典型例子是：http://<ip1>/callback;http://<ip2>/callback，并同时指定下面的 callbackHost 字段。
	// 在 callbackUrl 中使用 ip 的好处是减少对 dns 解析的依赖，可改善回调的性能和稳定性。指定 callbackUrl，必须指定 callbackbody，
	// 且值不能为空。
	CallbackURL string `json:"callback_url,omitempty"`

	// 上传成功后，七牛云向业务服务器发送回调通知时的 Host 值。与 callbackUrl 配合使用，仅当设置了 callbackUrl 时才有效。
	CallbackHost string `json:"callback_host,omitempty"`

	// 上传成功后，七牛云向业务服务器发送 Content-Type: application/x-www-form-urlencoded 的 POST 请求。业务服务器可以通过直接读取
	// 请求的 query 来获得该字段，支持魔法变量和自定义变量。callbackBody 要求是合法的 url query string。
	// 例如key=$(key)&hash=$(etag)&w=$(imageInfo.width)&h=$(imageInfo.height)。如果callbackBodyType指定为application/json，
	// 则callbackBody应为json格式，例如:{“key”:"$(key)",“hash”:"$(etag)",“w”:"$(imageInfo.width)",“h”:"$(imageInfo.height)"}。
	CallbackBody string `json:"callback_body,omitempty"`

	// 上传成功后，七牛云向业务服务器发送回调通知 callbackBody 的 Content-Type。默认为 application/x-www-form-urlencoded，也可设置
	// 为 application/json。
	CallbackBodyType string `json:"callback_body_type,omitempty"`

	// 资源上传成功后触发执行的预转持久化处理指令列表。fileType=2或3（上传归档存储或深度归档存储文件）时，不支持使用该参数。支持魔法变量和自
	// 定义变量。每个指令是一个 API 规格字符串，多个指令用;分隔。请参阅persistenOps详解与示例。同时添加 persistentPipeline 字段，使用专
	// 用队列处理，请参阅persistentPipeline。
	PersistentOps string `json:"persistent_ops,omitempty"`

	// 接收持久化处理结果通知的 URL。必须是公网上可以正常进行 POST 请求并能响应 HTTP/1.1 200 OK 的有效 URL。该 URL 获取的内容和持久化处
	// 理状态查询的处理结果一致。发送 body 格式是 Content-Type 为 application/json 的 POST 请求，需要按照读取流的形式读取请求的 body
	// 才能获取。
	PersistentNotifyURL string `json:"persistent_notify_url,omitempty"`

	// 转码队列名。资源上传成功后，触发转码时指定独立的队列进行转码。为空则表示使用公用队列，处理速度比较慢。建议使用专用队列。
	PersistentPipeline string `json:"persistent_pipeline,omitempty"`

	// saveKey 的优先级设置。为 true 时，saveKey不能为空，会忽略客户端指定的key，强制使用saveKey进行文件命名。参数不设置时，
	// 默认值为false
	ForceSaveKey bool `json:"force_save_key,omitempty"`

	// 开启 MimeType 侦测功能，并按照下述规则进行侦测；如不能侦测出正确的值，会默认使用 application/octet-stream 。
	// 默认设为 0 时：如上传端指定了 MimeType 则直接使用该值，否则按如下顺序侦测 MimeType 值：
	//		1. 检查文件扩展名；
	//		2. 检查 Key 扩展名；
	//		3. 侦测内容。
	// 设为 1 时：则忽略上传端传递的文件 MimeType 信息，并按如下顺序侦测 MimeType 值：
	//		1. 侦测内容；
	//		2. 检查文件扩展名；
	//		3. 检查 Key 扩展名。
	// 设为 -1 时：无论上传端指定了何值直接使用该值。
	DetectMime int `json:"detect_mime,omitempty"`

	CallbackFetchKey uint8 `json:"callback_fetch_key,omitempty"`

	DeleteAfterDays int `json:"delete_after_days,omitempty"`

	// 上传单链接限速，单位：bit/s；范围：819200 - 838860800（即800Kb/s - 800Mb/s），如果超出该范围将返回 400 错误
	TrafficLimit uint64 `json:"traffic_limit,omitempty"`
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
		ResumableAPIV2:         true,
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

	if up.CallbackURL != "" {
		callbackUrls := strings.Replace(up.CallbackURL, ",", ";", -1)
		up.CallbackURL = callbackUrls
		if len(up.CallbackBody) == 0 {
			up.CallbackBody = "key=$(key)&hash=$(etag)"
		}
		if len(up.CallbackBodyType) == 0 {
			up.CallbackBodyType = "application/x-www-form-urlencoded"
		}
	}

	return nil
}

func (up *UploadConfig) HitByPathPrefixes(localFileRelativePath string) (hit bool, pathPrefix string) {
	if len(up.SkipPathPrefixes) > 0 {
		// unpack skip prefix
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
		// unpack skip prefix
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
		// unpack fixed strings
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
