package operations

import (
	"fmt"
	"path/filepath"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type MoveInfo object.MoveApiInfo

func (info *MoveInfo) Check() *data.CodeError {
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

func Move(cfg *iqshell.Config, info MoveInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Move((*object.MoveApiInfo)(&info))
	if err != nil || result == nil {
		data.SetCmdStatusError()
		log.ErrorF("Move Failed, [%s:%s] => [%s:%s], Error: %v",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Move Success, [%s:%s] => [%s:%s]",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey)
	} else {
		data.SetCmdStatusError()
		log.ErrorF("Move Failed, [%s:%s] => [%s:%s], Code: %d, Error: %s",
			info.SourceBucket, info.SourceKey,
			info.DestBucket, info.DestKey,
			result.Code, result.Error)
	}
}

type BatchMoveInfo struct {
	BatchInfo    batch.Info
	SourceBucket string
	DestBucket   string
}

func (info *BatchMoveInfo) Check() *data.CodeError {
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

func BatchMove(cfg *iqshell.Config, info BatchMoveInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", cfg.CmdCfg.CmdId, info.SourceBucket, info.DestBucket, info.BatchInfo.InputFile))
		return filepath.Join(cmdPath, jobId)
	}
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
		return
	}

	batch.NewHandler(info.BatchInfo).
		SetFileExport(exporter).
		EmptyOperation(func() flow.Work {
			return &object.MoveApiInfo{}
		}).
		ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
			srcKey, destKey := items[0], items[0]
			if len(items) > 1 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				return &object.MoveApiInfo{
					SourceBucket: info.SourceBucket,
					SourceKey:    srcKey,
					DestBucket:   info.DestBucket,
					DestKey:      destKey,
					Force:        info.BatchInfo.Overwrite,
				}, nil
			}
			return nil, alert.Error("key invalid", "")
		}).
		OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
			apiInfo, ok := (operation).(*object.MoveApiInfo)
			if !ok {
				data.SetCmdStatusError()
				log.ErrorF("Change mimetype Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
				return
			}

			if result.IsSuccess() {
				log.InfoF("Move Success, [%s:%s] => [%s:%s]",
					apiInfo.SourceBucket, apiInfo.SourceKey,
					apiInfo.DestBucket, apiInfo.DestKey)
			} else {
				data.SetCmdStatusError()
				log.ErrorF("Move Failed, [%s:%s] => [%s:%s], Code: %d, Error: %s",
					apiInfo.SourceBucket, apiInfo.SourceKey,
					apiInfo.DestBucket, apiInfo.DestKey,
					result.Code, result.Error)
			}
		}).
		OnError(func(err *data.CodeError) {
			data.SetCmdStatusError()
			log.ErrorF("Batch move error:%v:", err)
		}).Start()
}
