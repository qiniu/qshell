package operations

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"time"
)

type DownloadOneInfo download.ApiInfo

func DownloadFile(info DownloadOneInfo) {
	err := downloadFile(download.ApiInfo(info))
	if err != nil {
		return
	}
}

type BatchDownloadInfo struct {
	BatchInfo group.Info
}

func BatchDownload(info BatchDownloadInfo) {

}

func downloadFile(info download.ApiInfo) (err error) {

	file := ""
	speed := ""
	log.InfoF("Download start:%s => %s", info.Url, info.ToFile)
	defer func() {
		if err != nil {
			log.ErrorF("Download  failed:%s => %s error:%v", info.Url, info.ToFile, err)
		} else {
			log.InfoF("Download success:%s => %s speed:%s", info.Url, file, speed)
		}
	}()

	startTime := time.Now().Unix()
	file, err = download.Download(info)
	if err != nil {
		return err
	}

	fileStatus, err := os.Stat(file)
	if err != nil {
		return errors.New("download speed: get file status error:" + err.Error())
	}
	if fileStatus == nil {
		return errors.New("download speed: can't get file status")
	}

	endTime := time.Now().Unix()

	speed = fmt.Sprintf("%.2fKB/s", float64(fileStatus.Size())/float64(endTime-startTime)/1024)
	return nil
}
