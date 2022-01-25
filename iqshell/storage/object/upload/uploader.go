package upload

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"path"
	"strings"
)

type ApiInfo struct {
	FilePath         string        // 文件路径，可为网络资源，也可为本地资源
	CheckExist       bool          // 检查服务端是否已存在
	CheckHash        bool          // 是否检查 hash, 检查是会对比服务端文件 hash
	CheckSize        bool          // 是否检查文件大小，检查是会对比服务端文件大小
	Overwrite        bool          // 当遇到服务端文件已存在时，是否使用本地文件覆盖之服务端的文件
	FileStatusDBPath string        // 文件上传状态信息保存的 db 路径
	ToBucket         string        // 文件保存至 bucket 的名称
	SaveKey          string        // 文件保存的名称
	TokenProvider    func() string // token provider
	TryTimes         int           // 失败时，最多重试次数【可选】
	FileSize         int64         // 待上传文件的大小, 如果不配置会动态读取 【可选】
	FileModifyTime   int64         // 本地文件修改时间, 如果不配置会动态读取 【可选】
}

func (a *ApiInfo) init() error {
	if len(a.FilePath) == 0 {
		return errors.New(alert.CannotEmpty("upload file path", ""))
	}

	// 获取文件信息
	if a.FileSize == 0 || a.FileModifyTime == 0 {
		localFileStatus, err := os.Stat(a.FilePath)
		if err != nil {
			return fmt.Errorf("get file:%s status error:%v", a.FilePath, err)
		}
		a.FileSize = localFileStatus.Size()
		a.FileModifyTime = localFileStatus.ModTime().UnixNano() / 100 // 兼容老版本：Unit is 100ns
	}

	if a.TryTimes == 0 {
		a.TryTimes = 3
	}

	if len(a.SaveKey) == 0 {
		a.SaveKey = path.Base(a.FilePath)
	}

	return nil
}

func (a *ApiInfo) isNetworkSource() bool {
	return strings.HasPrefix(a.FilePath, "http://") || strings.HasPrefix(a.FilePath, "https://")
}

type ApiResult struct {
	Key            string `json:"key"`
	FSize          int64  `json:"fsize"`
	Hash           string `json:"hash"`
	IsSkip         bool   `json:"is_skip"`       // 是否被 skip
	IsNotOverWrite bool   `json:"is_over_write"` // 是否因未开启 overwrite 而未被覆盖之前的上传
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
			CheckExist: info.CheckExist,
			CheckHash:  info.CheckHash,
			CheckSize:  info.CheckSize,
		}

		if checker.isNeedCheck() {
			// 检查服务端的数据
			exist, match, err = checker.check()
			if err != nil {
				err = errors.New("upload server check error:" + err.Error())
				return
			}
		} else {
			// 检查本地数据
			exist, match, err = d.checkInfoOfDB()
			if err != nil {
				log.WarningF("upload db check error:%v", err.Error())
			}
		}
	}

	if exist {
		if match {
			log.InfoF("File `%s` exists in bucket:[%s:%s], hash match, ignore this upload",
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

	log.InfoF("upload: start upload file:%s", d.FilePath)
	cfg := workspace.GetConfig()
	res, err = uploadLocalSource(info, cfg)
	log.InfoF("upload:   end upload file:%s error:%v", d.FilePath, err)

	if err != nil {
		err = errors.New("upload error:" + err.Error())
		return
	}

	err = d.saveInfoToDB()
	if err != nil {
		err = errors.New("upload: save upload info to db error:" + err.Error())
		return
	}

	return res, nil
}

func uploadLocalSource(info ApiInfo, cfg *config.Config) (result ApiResult, err error) {
	upCfg := cfg.Up
	storageCfg := workspace.GetStorageConfig()
	var up Uploader
	if info.FileSize < upCfg.PutThreshold {
		up = newFromUploader(storageCfg, &storage.PutExtra{
			Params:     nil,
			UpHost:     upCfg.UpHost,
			MimeType:   "",
			OnProgress: nil,
		})
	} else if upCfg.ResumableAPIV2 {
		up = newResumeV2Uploader(storageCfg, &storage.RputV2Extra{
			Recorder:   nil,
			Metadata:   nil,
			CustomVars: nil,
			UpHost:     upCfg.UpHost,
			MimeType:   "",
			PartSize:   upCfg.ResumableAPIV2PartSize,
			TryTimes:   info.TryTimes,
			Progresses: nil,
			Notify:     nil,
			NotifyErr:  nil,
		})
	} else {
		up = newResumeV1Uploader(storageCfg, &storage.RputExtra{
			Recorder:   nil,
			Params:     nil,
			UpHost:     upCfg.UpHost,
			MimeType:   "",
			ChunkSize:  0,
			TryTimes:   info.TryTimes,
			Progresses: nil,
			Notify:     nil,
			NotifyErr:  nil,
		})
	}

	log.DebugF("upload: start upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	result, err = up.upload(info)
	log.DebugF("upload:   end upload:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)

	return
}
