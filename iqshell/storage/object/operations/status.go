package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"time"
)

type StatusInfo object.StatusApiInfo

func Status(info StatusInfo) {
	result, err := object.Status(object.StatusApiInfo(info))
	if err != nil {
		log.ErrorF("Stat error:%v", err)
		return
	}
	log.Alert(getStatusInfo(info, result))
}

type BatchStatusInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

func BatchStatus(info BatchStatusInfo) {
	handler, err := NewBatchHandler(info.BatchInfo)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewFlow(info.BatchInfo.Info).ReadOperation(func() (operation batch.Operation, complete bool) {
		var in rs.BatchOperation = nil
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, true
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				in = object.StatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
				}
			}
		}
		return in, false
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.StatusApiInfo)
		if !ok {
			return
		}
		in := StatusInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%d\t%s\n", in.Key, result.Code, result.Error)
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

func getStatusInfo(info StatusInfo, status batch.OperationResult) string {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", info.Bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", info.Key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Hash:", status.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", status.FSize, utils.FormatFileSize(status.FSize))

	putTime := time.Unix(0, status.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", status.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", status.MimeType)
	if status.Type == 0 {
		statInfo += fmt.Sprintf("%-20s%d -> 标准存储\r\n", "FileType:", status.Type)
	} else {
		statInfo += fmt.Sprintf("%-20s%d -> 低频存储\r\n", "FileType:", status.Type)
	}
	return statInfo
}
