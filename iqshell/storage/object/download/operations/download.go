package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"time"
)

type DownloadInfo struct {
	Bucket        string // 文件被保存的 bucket
	Key           string // 文件被保存的 key
	ToFile        string // 文件保存的路径
	UseGetFileApi bool   //
	IsPublic      bool   //
	IoHost        string // io host, region 的 io 配置的可能为 ip, 搭配使用
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

	downloadDomain, downloadHost := getDownloadDomainAndHost(workspace.GetConfig(), &DownloadCfg{
		IoHost: info.IoHost,
		Bucket: info.Bucket,
	})
	if len(downloadDomain) == 0 && len(downloadHost) == 0 {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or io_host or bind domain to bucket", info.Bucket)
		return
	}

	log.DebugF("Download Domain:%s", downloadDomain)
	log.DebugF("Download Domain:%s", downloadHost)
	_, _ = downloadFile(&download.ApiInfo{
		IsPublic:       info.IsPublic,
		Domain:         downloadDomain,
		Host:           downloadHost,
		ToFile:         info.ToFile,
		StatusDBPath:   "",
		Referer:        "",
		FileEncoding:   "",
		Bucket:         info.Bucket,
		Key:            info.Key,
		FileModifyTime: 0,
		FileSize:       0,
		FileHash:       "",
		FromBytes:      0,
		UserGetFileApi: info.UseGetFileApi,
		Progress:       progress.NewPrintProgress(" 进度"),
	})
}

func downloadFile(info *download.ApiInfo) (download.ApiResult, error) {
	log.InfoF("Download [%s:%s] => %s", info.Bucket, info.Key, info.ToFile)
	startTime := time.Now().UnixNano() / 1e6
	res, err := download.Download(info)
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
