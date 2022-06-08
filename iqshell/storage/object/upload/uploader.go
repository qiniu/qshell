package upload

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"os"
	"path"
	"sync"
	"time"
)

type ApiInfo struct {
	FilePath            string            // 文件路径，可为网络资源，也可为本地资源
	ToBucket            string            // 文件保存至 bucket 的名称
	SaveKey             string            // 文件保存的名称
	MimeType            string            // 文件类型
	FileType            int               // 存储状态
	CheckExist          bool              // 检查服务端是否已存在此文件
	CheckHash           bool              // 是否检查 hash, 检查是会对比服务端文件 hash
	CheckSize           bool              // 是否检查文件大小，检查是会对比服务端文件大小
	Overwrite           bool              // 当遇到服务端文件已存在时，是否使用本地文件覆盖之服务端的文件
	UpHost              string            // 上传使用的域名
	TokenProvider       func() string     // token provider
	TryTimes            int               // 失败时，最多重试次数【可选】
	TryInterval         time.Duration     // 重试间隔时间 【可选】
	LocalFileSize       int64             // 待上传文件的大小, 如果不配置会动态读取 【可选】
	LocalFileModifyTime int64             // 待上传文件修改时间, 如果不配置会动态读取 【可选】
	DisableForm         bool              // 不使用 form 上传 【可选】
	DisableResume       bool              // 不使用分片上传 【可选】
	UseResumeV2         bool              // 分片上传时是否使用分片 v2 上传 【可选】
	ResumeWorkerCount   int               // 分片上传 worker 数量
	ChunkSize           int64             // 分片上传时的分片大小
	PutThreshold        int64             // 分片上传时上传阈值
	Progress            progress.Progress // 上传进度回调
}

func (a *ApiInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", a.ToBucket, a.SaveKey, a.FilePath)
}

func (a *ApiInfo) Check() (err *data.CodeError) {
	if len(a.FilePath) == 0 {
		return alert.CannotEmptyError("upload file path", "")
	}

	// 获取文件信息
	if a.LocalFileSize == 0 || a.LocalFileModifyTime == 0 {
		if utils.IsNetworkSource(a.FilePath) {
			a.LocalFileSize, err = utils.NetworkFileLength(a.FilePath)
			if err != nil {
				return data.NewEmptyError().AppendDescF("get network file:%s size error:%v", a.FilePath, err)
			}
		} else {
			localFileStatus, err := os.Stat(a.FilePath)
			if err != nil {
				return data.NewEmptyError().AppendDescF("get local file:%s status error:%v", a.FilePath, err)
			}
			a.LocalFileSize = localFileStatus.Size()
			a.LocalFileModifyTime = localFileStatus.ModTime().UnixNano() / 100 // 兼容老版本：Unit is 100ns
		}
	}

	if a.TryTimes == 0 {
		a.TryTimes = 3
	}

	if a.TryInterval == 0 {
		a.TryInterval = time.Second
	}

	if len(a.SaveKey) == 0 {
		a.SaveKey = path.Base(a.FilePath)
	}

	return nil
}

type ApiResult struct {
	Key            string `json:"key"`
	MimeType       string `json:"mime_type"` // 文件类型
	ServerFileSize int64  `json:"file_size"`     // 文件大小
	ServerFileHash string `json:"hash"`      // 文件 etag
	ServerPutTime  int64  `json:"put_time"`  // 文件上传时间
	IsSkip         bool   `json:"-"`         // 是否被 skip
	IsNotOverwrite bool   `json:"-"`         // 是否因未开启 overwrite 而未覆盖之前的上传
	IsOverwrite    bool   `json:"-"`         // 覆盖之前的上传
}

var _ flow.Result = (*ApiResult)(nil)

func (a *ApiResult) IsValid() bool {
	return len(a.Key) > 0 && len(a.MimeType) > 0 && len(a.ServerFileHash) > 0
}

func ApiResultFormat() string {
	return `{"key":"$(key)","hash":"$(etag)","file_size":$(fsize),"mime_type":"$(mimeType)"}`
}

type Uploader interface {
	upload(info *ApiInfo) (*ApiResult, *data.CodeError)
}

func Upload(info *ApiInfo) (res *ApiResult, err *data.CodeError) {
	err = info.Check()
	if err != nil {
		log.WarningF("upload: info init error:%v", err)
	}

	exist := false
	match := false
	if info.CheckExist {
		checkMode := object.MatchCheckModeFileSize
		if info.CheckHash {
			checkMode = object.MatchCheckModeFileHash
		}
		checkResult, mErr := object.Match(object.MatchApiInfo{
			Bucket:    info.ToBucket,
			Key:       info.SaveKey,
			LocalFile: info.FilePath,
			CheckMode: checkMode,
			FileSize:  info.LocalFileSize,
		})
		if checkResult != nil {
			exist = checkResult.Exist
			match = checkResult.Match
		}
		if mErr != nil {
			log.DebugF("check before upload error:%v", mErr)
		}
	}

	res = &ApiResult{}
	if exist {
		if match {
			log.InfoF("File `%s` exists in bucket:[%s:%s], and match, ignore this upload",
				info.FilePath, info.ToBucket, info.SaveKey)
			res.IsSkip = true
			return
		}

		if !info.Overwrite {
			log.WarningF("Skip upload of file `%s` => [%s:%s] because `overwrite` is false",
				info.FilePath, info.ToBucket, info.SaveKey)
			res.IsNotOverwrite = true
			return
		}
		res.IsOverwrite = true
	}

	log.DebugF("upload: start upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	res, err = uploadSource(info)
	if res == nil {
		res = &ApiResult{}
	}
	log.DebugF("upload:   end upload:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
	if err != nil {
		err = data.NewEmptyError().AppendDesc("upload error:" + err.Error())
		return
	}

	if stat, sErr := object.Status(object.StatusApiInfo{
		Bucket:   info.ToBucket,
		Key:      info.SaveKey,
		NeedPart: false,
	}); sErr != nil {
		err = data.NewEmptyError().AppendDesc("get file status after upload").AppendDescF("upload error:%v", sErr)
		return
	} else {
		res.ServerPutTime = stat.PutTime
	}

	if info.CheckHash {
		_, mErr := object.Match(object.MatchApiInfo{
			Bucket:    info.ToBucket,
			Key:       info.SaveKey,
			LocalFile: info.FilePath,
			CheckMode: object.MatchCheckModeFileHash,
			FileHash:  res.ServerFileHash,
			FileSize:  res.ServerFileSize,
		})
		return res, data.NewEmptyError().AppendDesc("check after upload ").AppendError(mErr)
	}

	return res, nil
}

var once sync.Once

func uploadSource(info *ApiInfo) (*ApiResult, *data.CodeError) {
	once.Do(func() {
		storage.SetSettings(&storage.Settings{
			TaskQsize: 0,
			Workers:   0,
			ChunkSize: 0,
			PartSize:  0,
			TryTimes:  info.ResumeWorkerCount,
		})
	})
	storageCfg := workspace.GetStorageConfig()
	var up Uploader
	if utils.IsNetworkSource(info.FilePath) {
		up = networkSourceUploader(info, storageCfg)
	} else {
		up = localSourceUploader(info, storageCfg)
	}
	return up.upload(info)
}

func localSourceUploader(info *ApiInfo, storageCfg *storage.Config) (up Uploader) {
	if info.DisableResume || (!info.DisableForm && info.LocalFileSize < info.PutThreshold) {
		up = newFromUploader(storageCfg, &storage.PutExtra{
			Params:     nil,
			UpHost:     info.UpHost,
			MimeType:   info.MimeType,
			OnProgress: nil,
		})
	} else if info.UseResumeV2 {
		up = newResumeV2Uploader(storageCfg)
	} else {
		up = newResumeV1Uploader(storageCfg)
	}
	return
}

func networkSourceUploader(info *ApiInfo, storageCfg *storage.Config) (up Uploader) {
	return newConveyorUploader(storageCfg)
}
