package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"time"
)

type StatusInfo object.StatusApiInfo

func (info *StatusInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func Status(cfg *iqshell.Config, info StatusInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Status(object.StatusApiInfo(info))
	if err != nil {
		log.ErrorF("Status Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Status Success, [%s:%s]", info.Bucket, info.Key)
		log.Alert(getResultInfo(info.Bucket, info.Key, result))
	}
}

type BatchStatusInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchStatusInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchStatus(cfg *iqshell.Config, info BatchStatusInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	handler, err := group.NewHandler(info.BatchInfo.Info)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewFlow(info.BatchInfo).ReadOperation(func() (operation batch.Operation, hasMore bool) {
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, false
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				return object.StatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.StatusApiInfo)
		if !ok {
			return
		}
		in := StatusInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%d\t%s", in.Key, result.Code, result.Error)
			log.ErrorF("Status Failed, [%s:%s], Code: %d, Error: %s", in.Bucket, in.Key, result.Code, result.Error)
		} else {
			status := fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d",
				in.Key, result.FSize, result.Hash, result.MimeType, result.PutTime, result.Type)
			handler.Export().Success().Export(status)
			log.Alert(status)
		}
	}).OnError(func(err error) {
		log.ErrorF("Batch Status error:%v:", err)
	}).Start()
}

func getResultInfo(bucket, key string, status object.StatusResult) string {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "FileHash:", status.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", status.FSize, utils.FormatFileSize(status.FSize))

	putTime := time.Unix(0, status.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", status.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", status.MimeType)

	resoreStatus := ""
	if status.RestoreStatus > 0 {
		if status.RestoreStatus == 1 {
			resoreStatus = "解冻中"
		} else if status.RestoreStatus == 2 {
			resoreStatus = "解冻完成"
		}
	}
	if len(resoreStatus) > 0 {
		statInfo += fmt.Sprintf("%-20s%d(%s)\r\n", "RestoreStatus:", status.RestoreStatus, resoreStatus)
	}

	if status.Expiration > 0 {
		expiration := time.Unix(status.Expiration, 0)
		statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Expiration:", status.Expiration, expiration.String())
	}

	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "FileType:", status.Type, getStorageTypeDescription(status.Type))

	return statInfo
}

var objectTypes = []string{"标准存储", "低频存储", "归档存储", "深度归档存储"}

func getStorageTypeDescription(storageType int) string {
	typeString := "未知类型"
	if storageType >= 0 && storageType < len(objectTypes) {
		typeString = objectTypes[storageType]
	}
	return typeString
}
