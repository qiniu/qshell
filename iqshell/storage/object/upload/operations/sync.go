package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"time"
)

type SyncUploadInfo struct {
	ResourceUrl string
	Bucket      string
	Key         string
	UpHostIp    string
	IsResumeV2  bool
}

func SyncUpload(info SyncUploadInfo) {
	if len(info.Bucket) == 0 {
		log.Error(alert.CannotEmpty("bucket", ""))
		return
	}

	if len(info.ResourceUrl) == 0 {
		log.Error(alert.CannotEmpty("resource url", ""))
		return
	}

	if len(info.Key) == 0 {
		if key, err := utils.KeyFromUrl(info.ResourceUrl); err != nil {
			log.ErrorF("get path as key: %v\n", err)
			os.Exit(data.STATUS_ERROR)
		} else {
			info.Key = key
		}
	}

	//sync
	tStart := time.Now()
	syncRet, sErr := upload.Sync(info.ResourceUrl, info.Bucket, info.Key, info.UpHostIp, info.IsResumeV2)
	if sErr != nil {
		log.Error("%v", sErr)
		os.Exit(data.STATUS_ERROR)
	}

	fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", info.ResourceUrl, info.Bucket, info.Key, time.Since(tStart))
	fmt.Println("Hash:", syncRet.Hash)
	fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, utils.FormatFileSize(syncRet.Fsize))
	fmt.Println("Mime:", syncRet.MimeType)
}
