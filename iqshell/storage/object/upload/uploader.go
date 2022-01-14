package upload

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"strings"
)

type ApiInfo struct {
	FilePath         string        // 文件路径，可为网络资源，也可为本地资源
	FileSize         int64         // 文件大小
	FileModifyTime   int64         // 本地文件修改时间
	CheckExist       bool          // 检查服务端是否已存在
	CheckHash        bool          // 是否检查 hash, 检查是会对比服务端文件 hash
	CheckSize        bool          // 是否检查文件大小，检查是会对比服务端文件大小
	FileStatusDBPath string        // 文件上传状态想你想保存的 db 路径
	ToBucket         string        // 文件保存至 bucket 的名称
	SaveKey          string        // 文件保存的名称
	TokenProvider    func() string // token provider
}

func (u *ApiInfo) isNetworkSource() bool {
	return strings.HasPrefix(u.FilePath, "http://") || strings.HasPrefix(u.FilePath, "https://")
}

type UploadResult struct {
	Key    string `json:"key"`
	FSize  int64  `json:"fsize"`
	Hash   string `json:"hash"`
	IsSkip bool   `json:"is_skip"` // 是否被 skip
}

type Uploader interface {
	upload(info ApiInfo) (UploadResult, error)
}

func Upload(info ApiInfo) (res UploadResult, err error) {

	d := &dbHandler{
		DBFilePath:     info.FileStatusDBPath,
		FilePath:       info.FilePath,
		FileUpdateTime: info.FileModifyTime,
	}
	err = d.init()
	if err != nil {
		log.WarningF("upload: db init error:%v", err)
	}

	// 检查服务端是否存在
	if info.CheckExist {
		if info.CheckHash {

		}
		if info.CheckSize {

		}
	} else {
		// 检查本地是否保存有上传数据，已上传会被保存。
		err = d.checkInfoOfDB()
		if err == nil {
			res.Key = info.SaveKey
			res.FSize = info.FileSize
			res.IsSkip = true
			log.InfoF("upload: file already upload:%s", d.FilePath)
			return
		}
		log.WarningF("upload check info:%v", err)
	}

	log.InfoF("upload: start upload file:%s", d.FilePath)
	cfg := workspace.GetConfig()
	if info.isNetworkSource() {
		res, err = uploadNetworkSource(info, cfg)
	} else {
		res, err = uploadLocalSource(info, cfg)
	}
	log.InfoF("upload:   end upload file:%s error:%v", d.FilePath, err)

	if err != nil {
		err =  errors.New("upload error:" + err.Error())
		return
	}

	err = d.saveInfoToDB()
	if err != nil {
		err = errors.New("upload: save upload info to db error:" + err.Error())
		return
	}

	return res,nil
}

func uploadNetworkSource(info ApiInfo, cfg *config.Config) (result UploadResult, err error) {

	return
}

func uploadLocalSource(info ApiInfo, cfg *config.Config) (result UploadResult, err error) {
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
			TryTimes:   3,
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
			TryTimes:   3,
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
