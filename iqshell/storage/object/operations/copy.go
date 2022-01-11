package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type CopyInfo object.CopyApiInfo

func Copy(info CopyInfo) {
	result, err := object.Copy(object.CopyApiInfo(info))
	if err != nil {
		log.ErrorF("Copy error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Copy error:%v", result.Error)
		return
	}
}

type BatchCopyInfo struct {
	BatchInfo    batch.Info
	SourceBucket string
	DestBucket   string
}

func BatchCopy(info BatchCopyInfo) {
	handler, err := group.NewHandler(info.BatchInfo.Info)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewFlow(info.BatchInfo).ReadOperation(func() (operation batch.Operation, complete bool) {
		var in batch.Operation = nil
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, true
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) > 0 {
			srcKey, destKey := items[0], items[0]
			if len(items) > 1 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				in = object.CopyApiInfo{
					SourceBucket: info.SourceBucket,
					SourceKey:    srcKey,
					DestBucket:   info.DestBucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Force,
				}
			}
		}
		return in, false
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.CopyApiInfo)
		if !ok {
			return
		}
		in := CopyInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s\n", in.SourceKey, in.DestKey, result.Code, result.Error)
			log.ErrorF("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s\n", in.SourceKey, in.DestKey)
			log.ErrorF("Copy '%s:%s' => '%s:%s' success\n",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch copy error:%v:", err)
	}).Start()
}
