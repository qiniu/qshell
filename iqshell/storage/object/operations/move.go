package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type MoveInfo object.MoveApiInfo

func Move(info MoveInfo) {
	result, err := object.Move(object.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Move error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Move error:%v", result.Error)
		return
	}
}

type BatchMoveInfo struct {
	BatchInfo    BatchInfo
	SourceBucket string
	DestBucket   string
}

func BatchMove(info BatchMoveInfo) {
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
			srcKey, destKey := items[0], items[0]
			if len(items) > 1 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				in = &object.MoveApiInfo{
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
		apiInfo, ok := (operation).(object.MoveApiInfo)
		if !ok {
			return
		}

		in := MoveInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s\n", in.SourceKey, in.DestKey, result.Code, result.Error)
			log.ErrorF("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s\n", in.SourceKey, in.DestKey)
			log.ErrorF("Move '%s:%s' => '%s:%s' success\n",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch move error:%v:", err)
	}).Start()
}
