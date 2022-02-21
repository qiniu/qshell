package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type MoveInfo object.MoveApiInfo

func (info *MoveInfo) Check() error {
	if len(info.SourceBucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.SourceKey) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.DestBucket) == 0 {
		return alert.CannotEmptyError("DestBucket", "")
	}
	if len(info.DestKey) == 0 {
		return alert.CannotEmptyError("DestKey", "")
	}
	return nil
}

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

	log.InfoF("Move '%s:%s' => '%s:%s' success",
		info.SourceBucket, info.SourceKey,
		info.DestBucket, info.DestKey)
}

type BatchMoveInfo struct {
	BatchInfo    batch.Info
	SourceBucket string
	DestBucket   string
}

func (info *BatchMoveInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.SourceBucket) == 0 {
		return alert.CannotEmptyError("SrcBucket", "")
	}

	if len(info.DestBucket) == 0 {
		return alert.CannotEmptyError("DestBucket", "")
	}

	return nil
}

func BatchMove(info BatchMoveInfo) {
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
			srcKey, destKey := items[0], items[0]
			if len(items) > 1 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				return object.MoveApiInfo{
					SourceBucket: info.SourceBucket,
					SourceKey:    srcKey,
					DestBucket:   info.DestBucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Force,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.MoveApiInfo)
		if !ok {
			return
		}

		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s", apiInfo.SourceKey, apiInfo.DestKey, result.Code, result.Error)
			log.ErrorF("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				apiInfo.SourceBucket, apiInfo.SourceKey,
				apiInfo.DestBucket, apiInfo.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s", apiInfo.SourceKey, apiInfo.DestKey)
			log.InfoF("Move '%s:%s' => '%s:%s' success",
				apiInfo.SourceBucket, apiInfo.SourceKey,
				apiInfo.DestBucket, apiInfo.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch move error:%v:", err)
	}).Start()
}
