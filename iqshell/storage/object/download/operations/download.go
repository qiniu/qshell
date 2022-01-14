package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"time"
)

type DownloadInfo struct {
	download.ApiInfo
	IsPublic bool // 是否是公有云
}

func DownloadFile(info DownloadInfo) {
	_, _ = downloadFile(info)
}

func downloadFile(info DownloadInfo) (download.ApiResult, error) {
	log.InfoF("Download start:%s => %s", info.Url, info.ToFile)

	// 构造下载 url
	if info.IsPublic {
		info.Url = download.PublicUrl(download.UrlApiInfo{
			BucketDomain: info.Domain,
			Key:          info.Key,
			UseHttps:     workspace.GetConfig().IsUseHttps(),
		})
	} else {
		info.Url = download.PrivateUrl(download.UrlApiInfo{
			BucketDomain: info.Domain,
			Key:          info.Key,
			UseHttps:     workspace.GetConfig().IsUseHttps(),
		})
	}

	startTime := time.Now().UnixNano() / 1e6
	res, err := download.Download(info.ApiInfo)
	if err != nil {
		log.ErrorF("Download  failed:%s => %s error:%v", info.Url, info.ToFile, err)
		return res, err
	}

	fileStatus, err := os.Stat(res.FileAbsPath)
	if err != nil {
		log.ErrorF("Download  failed:%s => %s get file status error:%v", info.Url, info.ToFile, err)
		return res, err
	}
	if fileStatus == nil {
		log.ErrorF("Download  failed:%s => %s download speed: can't get file status", info.Url, info.ToFile)
		return res, err
	}

	endTime := time.Now().UnixNano() / 1e6
	duration := float64(endTime - startTime) / 1000
	speed := fmt.Sprintf("%.2fKB/s", float64(fileStatus.Size())/duration/1024)
	if res.IsExist {
		log.Alert("Download skip because file exist:%s => %s", info.Url, res.FileAbsPath)
	} else if res.IsUpdate {
		log.Alert("Download update success:%s => %s speed:%s", info.Url, res.FileAbsPath, speed)
	} else {
		log.Alert("Download success:%s => %s speed:%s", info.Url, res.FileAbsPath, speed)
	}

	return res, nil
}
