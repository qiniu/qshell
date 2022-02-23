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
		log.ErrorF("Stat error:%v", err)
		return
	}
	log.Alert(getResultInfo(info.Bucket, info.Key, result))
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
			log.ErrorF("Status '%s' Failed, Code: %d, Error: %s", in.Key, result.Code, result.Error)
		} else {
			status := fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d",
				in.Key, result.FSize, result.Hash, result.MimeType, result.PutTime, result.Type)
			handler.Export().Success().Export(status)
			log.Alert(status)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch Status error:%v:", err)
	}).Start()
}

var objectTypes = []string{"标准存储", "低频存储", "归档存储", "深度归档存储"}

func getResultInfo(bucket, key string, status batch.OperationResult) string {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "FileHash:", status.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", status.FSize, utils.FormatFileSize(status.FSize))

	putTime := time.Unix(0, status.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", status.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", status.MimeType)

	typeString := "未知类型"
	if status.Type >= 0 && status.Type < len(objectTypes) {
		typeString = objectTypes[status.Type]
	}
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "FileType:", status.Type, typeString)

	return statInfo
}
