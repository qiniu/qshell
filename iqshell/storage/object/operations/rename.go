package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

// rename 实际用的还是 move

type RenameInfo object.MoveApiInfo

func (info *RenameInfo) Check() *data.CodeError {
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

	result, err := object.Move((*object.MoveApiInfo)(&info))
	if err != nil {
		log.ErrorF("Rename Failed, [%s:%s] => [%s:%s], Error: %v",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Rename Failed, [%s:%s] => [%s:%s], Code: %d, Error: %s",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Rename '%s:%s' => '%s:%s' success",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey)
	}
}

type BatchRenameInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchRenameInfo) Check() *data.CodeError {
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

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewHandler(info.BatchInfo).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		if len(items) > 1 {
			sourceKey, destKey := items[0], items[1]
			if sourceKey != "" && destKey != "" {
				return &object.MoveApiInfo{
					SourceBucket: info.Bucket,
					SourceKey:    sourceKey,
					DestBucket:   info.Bucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Force,
				}, nil
			} else {
				return nil, alert.Error("key invalid", "")
			}
		}
		return nil, alert.Error("need more than one param", "")
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.MoveApiInfo)
		if !ok {
			return
		}
		in := (*RenameInfo)(apiInfo)
		if result.Code != 200 || result.Error != "" {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Rename Failed, [%s:%s] => [%s:%s], Code: %d, Error: %s",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey,
				result.Code, result.Error)
		} else {
			exporter.Success().Export(operationInfo)
			log.InfoF("Rename Success, [%s:%s] => [%s:%s]",
				in.SourceBucket, in.SourceKey,
				in.DestBucket, in.DestKey)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch rename error:%v:", err)
	}).Start()
}
