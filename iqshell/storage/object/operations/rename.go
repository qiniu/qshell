package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

// rename 实际用的还是 move

type RenameInfo object.MoveApiInfo

func (info *RenameInfo) Check() error {
	if len(info.SourceBucket) == 0 {
		return alert.CannotEmptyError("SourceBucket", "")
	}
	if len(info.SourceKey) == 0 {
		return alert.CannotEmptyError("SourceKey", "")
	}
	if len(info.DestBucket) == 0 {
		return alert.CannotEmptyError("DestBucket", "")
	}
	if len(info.DestKey) == 0 {
		return alert.CannotEmptyError("DestKey", "")
	}
	return nil
}

func Rename(cfg *iqshell.Config, info RenameInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Move(object.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Rename error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Rename error:%v", result.Error)
		return
	}

	log.InfoF("Rename '%s:%s' => '%s:%s' success",
		info.SourceBucket, info.SourceKey,
		info.DestBucket, info.DestKey)
}

type BatchRenameInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchRenameInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchRename(cfg *iqshell.Config, info BatchRenameInfo) {
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
		if len(items) > 1 {
			sourceKey, destKey := items[0], items[1]
			if sourceKey != "" && destKey != "" {
				return object.MoveApiInfo{
					SourceBucket: info.Bucket,
					SourceKey:    sourceKey,
					DestBucket:   info.Bucket,
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
		in := RenameInfo(apiInfo)
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s", in.SourceKey, in.DestKey, result.Code, result.Error)
			log.ErrorF("Rename '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s", in.SourceKey, in.DestKey)
			log.InfoF("Rename '%s:%s' => '%s:%s' success",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch rename error:%v:", err)
	}).Start()
}
