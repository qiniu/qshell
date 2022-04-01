package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"time"
)

type StatusInfo object.StatusApiInfo

func (info *StatusInfo) Check() *data.CodeError {
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

func (info *BatchStatusInfo) Check() *data.CodeError {
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

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewHandler(info.BatchInfo).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		key := items[0]
		if key != "" {
			return &object.StatusApiInfo{
				Bucket: info.Bucket,
				Key:    key,
			}, nil
		}
		return nil, alert.Error("key invalid", "")
	}).OnResult(func(operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.StatusApiInfo)
		if !ok {
			return
		}
		in := (*StatusInfo)(apiInfo)
		if result.Code != 200 || result.Error != "" {
			exporter.Fail().ExportF("%s\t%d\t%s", in.Key, result.Code, result.Error)
			log.ErrorF("Status Failed, [%s:%s], Code: %d, Error: %s", in.Bucket, in.Key, result.Code, result.Error)
		} else {
			status := fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d",
				in.Key, result.FSize, result.Hash, result.MimeType, result.PutTime, result.Type)
			exporter.Success().Export(status)
			log.Alert(status)
		}
	}).OnError(func(err *data.CodeError) {
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
