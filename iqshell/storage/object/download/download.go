package download

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

type ApiInfo struct {
	Url            string // 文件下载的 url 【必填】
	Domain         string // 文件下载的 domain 【必填】
	ToFile         string // 文件保存的路径 【必填】
	FileSize       int64  // 文件大小，有值则会检测文件大小 【选填】
	FileModifyTime int64  // 文件修改时间 【选填】
	StatusDBPath   string // 下载状态缓存的 db 路径 【选填】
	Referer        string // 请求 header 中的 Referer 【选填】
	FileEncoding   string // 文件编码方式 【选填】
	Bucket         string // 文件所在 bucket，用于验证 hash 【选填】
	Key            string // 文件被保存的 key，用于验证 hash 【选填】
	FileHash       string // 文件 hash，有值则会检测 hash 【选填】
}

type ApiResult struct {
	FileAbsPath string // 文件被保存的绝对路径
	IsUpdate    bool   // 是否为接续下载
	IsExist     bool   // 是否为已存在
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info ApiInfo) (res ApiResult, err error) {
	if len(info.ToFile) == 0 {
		err = errors.New("the filename saved after downloading is empty")
		return
	}

	f, err := createDownloadFiles(info.ToFile, info.FileEncoding)
	if err != nil {
		return
	}

	// 检查文件是否已存在，如果存在是否符合预期
	dbChecker := &dbHandler{
		DBFilePath:           info.StatusDBPath,
		FilePath:             f.toAbsFile,
		FileHash:             info.FileHash,
		FileSize:             info.FileSize,
		FileServerUpdateTime: info.FileModifyTime,
	}
	err = dbChecker.init()
	if err != nil {
		return
	}

	shouldDownload := true

	// 文件存在则检查文件状态
	fileStatus, err := os.Stat(f.toAbsFile)
	tempFileStatus, tempErr := os.Stat(f.tempFile)
	if err == nil || os.IsExist(err) || tempErr == nil || os.IsExist(tempErr) {
		// 检查服务端文件是否变更
		if cErr := dbChecker.checkInfoOfDB(); cErr != nil {
			log.WarningF("Local file `%s` exist for key `%s`, but not match:%v", f.toAbsFile, info.Key, cErr)
			if e := f.clean(); e != nil {
				log.WarningF("Local file `%s` exist for key `%s`, clean error:%v", f.toAbsFile, info.Key, e)
			}
			if sErr := dbChecker.saveInfoToDB(); sErr != nil {
				log.WarningF("Local file `%s` exist for key `%s`, save info to db clean error:%v", f.toAbsFile, info.Key, sErr)
			}
		}
		if tempFileStatus != nil && tempFileStatus.Size() > 0 {
			res.IsUpdate = true
		}
		// 文件是否已下载完成，如果完成跳过下载阶段，直接验证
		if fileStatus != nil && fileStatus.Size() == info.FileSize {
			res.IsExist = true
			shouldDownload = false
		}
	} else {
		if sErr := dbChecker.saveInfoToDB(); sErr != nil {
			log.WarningF("Local file `%s` not exist for key `%s`, save info to db clean error:%v", f.toAbsFile, info.Key, sErr)
		}
	}

	res.FileAbsPath = f.toAbsFile

	// 下载
	if shouldDownload {
		downloader := &Downloader{
			files:   f,
			Url:     info.Url,
			Domain:  info.Domain,
			Referer: info.Referer,
		}
		err = downloader.Download()
		if err != nil {
			return
		}
		err = dbChecker.saveInfoToDB()
		if err != nil {
			return
		}
	}

	// 检查下载后的数据是否符合预期
	err = (&LocalFileInfo{
		File:                f.toAbsFile,
		Bucket:              info.Bucket,
		Key:                 info.Key,
		FileHash:            info.FileHash,
		FileSize:            info.FileSize,
		RemoveFileWhenError: true,
	}).CheckDownloadFile()
	if err != nil {
		return
	}

	return
}
