package upload

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path"
	"time"
)

type ApiInfo struct {
	FilePath         string        // 文件路径，可为网络资源，也可为本地资源
	ToBucket         string        // 文件保存至 bucket 的名称
	SaveKey          string        // 文件保存的名称
	MimeType         string        // 文件类型
	FileType         int           // 存储状态
	CheckExist       bool          // 检查服务端是否已存在此文件
	CheckHash        bool          // 是否检查 hash, 检查是会对比服务端文件 hash
	CheckSize        bool          // 是否检查文件大小，检查是会对比服务端文件大小
	Overwrite        bool          // 当遇到服务端文件已存在时，是否使用本地文件覆盖之服务端的文件
	UpHost           string        // 上传使用的域名
	FileStatusDBPath string        // 文件上传状态信息保存的 db 路径
	TokenProvider    func() string // token provider
	TryTimes         int           // 失败时，最多重试次数【可选】
	TryInterval      time.Duration // 重试间隔时间 【可选】
	FileSize         int64         // 待上传文件的大小, 如果不配置会动态读取 【可选】
	FileModifyTime   int64         // 本地文件修改时间, 如果不配置会动态读取 【可选】
	DisableForm      bool          // 不使用 form 上传 【可选】
	DisableResume    bool          // 不使用分片上传 【可选】
	UseResumeV2      bool          // 分片上传时是否使用分片 v2 上传 【可选】
	ChunkSize        int64         // 分片上传时的分片大小
	PutThreshold     int64         // 分片上传时上传阈值
}

func (a *ApiInfo) init() (err error) {
	if len(a.FilePath) == 0 {
		return errors.New(alert.CannotEmpty("upload file path", ""))
	}

	// 获取文件信息
	if a.FileSize == 0 || a.FileModifyTime == 0 {
		if utils.IsNetworkSource(a.FilePath) {
			a.FileSize, err = utils.NetworkFileLength(a.FilePath)
			if err != nil {
				return fmt.Errorf("get network file:%s size error:%v", a.FilePath, err)
			}
		} else {
			localFileStatus, err := os.Stat(a.FilePath)
			if err != nil {
				return fmt.Errorf("get local file:%s status error:%v", a.FilePath, err)
			}
			a.FileSize = localFileStatus.Size()
			a.FileModifyTime = localFileStatus.ModTime().UnixNano() / 100 // 兼容老版本：Unit is 100ns
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
	FSize          int64  `json:"file_size"` // 文件大小
	Hash           string `json:"hash"`      // 文件 etag
	IsSkip         bool   `json:"-"`         // 是否被 skip
	IsNotOverWrite bool   `json:"-"`         // 是否因未开启 overwrite 而未被覆盖之前的上传
}

func ApiResultFormat() string {
	return `{"key":"$(key)","hash":"$(etag)","file_size":$(fsize),"mime_type":"$(mimeType)"}`
}

type Uploader interface {
	upload(info ApiInfo) (ApiResult, error)
}

func Upload(info ApiInfo) (res ApiResult, err error) {
	err = info.init()
	if err != nil {
		log.WarningF("upload: info init error:%v", err)
	}

	d := &dbHandler{
		DBFilePath:     info.FileStatusDBPath,
		FilePath:       info.FilePath,
		FileUpdateTime: info.FileModifyTime,
	}
	err = d.init()
	if err != nil {
		log.WarningF("upload: db init error:%v", err)
	}

	exist := false
	match := false
	if info.CheckExist {
		checker := &serverChecker{
			Bucket:     info.ToBucket,
			Key:        info.SaveKey,
			FilePath:   info.FilePath,
			FileSize:   info.FileSize,
			CheckExist: info.CheckExist,
			CheckHash:  info.CheckHash,
			CheckSize:  info.CheckSize,
		}

		// 检查服务端的数据
		exist, match, err = checker.check()
		if err != nil {
			log.WarningF("upload server check error:%v", err.Error())
		}
	} else {
		// 检查本地数据
		exist, match, err = d.checkInfoOfDB()
		if err != nil {
			log.WarningF("upload db check error:%v", err.Error())
		}
	}

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
			res.IsNotOverWrite = true
			return
		}
	}

	log.DebugF("upload: start upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	res, err = uploadSource(info)
	log.DebugF("upload:   end upload:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)

	if err != nil {
		err = errors.New("upload error:" + err.Error())
		return
	}

	err = d.saveInfoToDB()
	if err != nil {
		log.WarningF("upload: save upload info to db error:%v", err)
	}

	return res, nil
}

func uploadSource(info ApiInfo) (ApiResult, error) {
	storageCfg := workspace.GetStorageConfig()
	var up Uploader
	if utils.IsNetworkSource(info.FilePath) {
		up = networkSourceUploader(info, storageCfg)
	} else {
		up = localSourceUploader(info, storageCfg)
	}
	return up.upload(info)
}

func localSourceUploader(info ApiInfo, storageCfg *storage.Config) (up Uploader) {
	if info.DisableResume || (!info.DisableForm && info.FileSize < info.PutThreshold) {
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

func networkSourceUploader(info ApiInfo, storageCfg *storage.Config) (up Uploader) {
	return newConveyorUploader(storageCfg)
}
