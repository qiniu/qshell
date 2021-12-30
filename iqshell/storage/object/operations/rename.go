package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

// rename 实际用的还是 move

type RenameInfo object.MoveApiInfo

func Rename(info RenameInfo) {
	result, err := object.Move(object.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Rename error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Rename error:%v", result.Error)
		return
	}
}

type BatchRenameInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

func BatchRename(info BatchRenameInfo) {
	handler, err := NewBatchHandler(info.BatchInfo)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewFlow(info.BatchInfo.Info).ReadOperation(func() (operation batch.Operation, complete bool) {
		var in batch.Operation = nil
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, true
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		if len(items) > 1 {
			sourceKey, destKey := items[0], items[1]
			if sourceKey != "" && destKey != "" {
				in = object.MoveApiInfo{
					SourceBucket: info.Bucket,
					SourceKey:    sourceKey,
					DestBucket:   info.Bucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Force,
				}
			}
		}
		return in, false
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.MoveApiInfo)
		if !ok {
			return
		}
		in := RenameInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s\n", in.SourceKey, in.DestKey, result.Code, result.Error)
			log.ErrorF("Rename '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s\n", in.SourceKey, in.DestKey)
			log.ErrorF("Rename '%s:%s' => '%s:%s' success\n",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch rename error:%v:", err)
	}).Start()
}
