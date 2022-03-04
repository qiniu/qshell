package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"time"
)

type DownloadInfo struct {
	download.ApiInfo
	IsPublic bool // 是否是公有云
}

func (info *DownloadInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func DownloadFile(cfg *iqshell.Config, info DownloadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	// 如果 ToFile 不存在则保存在当前文件录下，文件名为：key
	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}

	if len(info.Domain) == 0 {
		log.DebugF("get domain of bucket:%s", info.Bucket)
		if d, err := bucket.DomainOfBucket(info.Bucket); err != nil {
			log.ErrorF("get bucket domain error:%v, domain can't be empty", err)
			return
		} else {
			info.Domain = d
			log.DebugF("bucket:%s domain:%s", info.Bucket, info.Domain)
		}
	}
	_, _ = downloadFile(info)
}

func downloadFile(info DownloadInfo) (download.ApiResult, error) {
	log.InfoF("Download [%s:%s] => %s", info.Bucket, info.Key, info.ToFile)

	startTime := time.Now().UnixNano() / 1e6
	res, err := download.Download(info.ApiInfo)
	if err != nil {
		log.ErrorF("Download  failed, [%s:%s] => %s error:%v", info.Bucket, info.Key, info.ToFile, err)
		return res, err
	}

	fileStatus, err := os.Stat(res.FileAbsPath)
	if err != nil {
		log.ErrorF("Download  failed, [%s:%s] => %s get file status error:%v", info.Bucket, info.Key, info.ToFile, err)
		return res, err
	}
	if fileStatus == nil {
		log.ErrorF("Download  failed, [%s:%s] => %s download speed: can't get file status", info.Bucket, info.Key, info.ToFile)
		return res, err
	}

	endTime := time.Now().UnixNano() / 1e6
	duration := float64(endTime-startTime) / 1000
	speed := fmt.Sprintf("%.2fKB/s", float64(fileStatus.Size())/duration/1024)
	if res.IsExist {
		log.AlertF("Download skip because file exist, [%s:%s] => %s", info.Bucket, info.Key, res.FileAbsPath)
	} else if res.IsUpdate {
		log.AlertF("Download update success, [%s:%s] => %s speed:%s", info.Bucket, info.Key, res.FileAbsPath, speed)
	} else {
		log.AlertF("Download success, [%s:%s] => %s speed:%s", info.Bucket, info.Key, res.FileAbsPath, speed)
	}

	return res, nil
}
