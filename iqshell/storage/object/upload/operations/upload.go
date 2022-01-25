package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"time"
)

type UploadInfo struct {
	upload.ApiInfo
	IsPublic bool // 是否是公有云
}

func UploadFile(info UploadInfo)  {
	_, err := uploadFile(info.ApiInfo)
	if err != nil {
		log.Error(err)
	}
}

func uploadFile(info upload.ApiInfo) (upload.ApiResult, error) {
	startTime := time.Now().UnixNano() / 1e6
	res, err := upload.Upload(info)
	if err != nil {
		log.ErrorF("Upload  failed:%s => [%s:%s] error:%v", info.FilePath, info.ToBucket, info.SaveKey, err)
		return res, err
	}
	endTime := time.Now().UnixNano() / 1e6

	duration := float64(endTime - startTime) / 1000
	speed := fmt.Sprintf("%.2fKB/s", float64(res.FSize)/duration/1024)
	if res.IsSkip {
		log.Alert("Upload skip because file exist:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)
	} else {
		log.Alert("Download success%s => [%s:%s] speed:%s", info.FilePath, info.ToBucket, speed)
	}

	return res, nil
}


