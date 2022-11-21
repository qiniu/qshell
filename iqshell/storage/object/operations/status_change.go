package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
)

type ForbiddenInfo struct {
	Bucket      string
	Key         string
	UnForbidden bool
}

func (info *ForbiddenInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func (info *ForbiddenInfo) getStatus() int {
	// 0:启用  1:禁用
	if info.UnForbidden {
		return 0
	} else {
		return 1
	}
}

func (info *ForbiddenInfo) getStatusDesc() string {
	// 0:启用  1:禁用
	if info.UnForbidden {
		return "启用"
	} else {
		return "禁用"
	}
}

func ForbiddenObject(cfg *iqshell.Config, info ForbiddenInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.ChangeStatus(&object.ChangeStatusApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Status: info.getStatus(),
	})

	statusDesc := info.getStatusDesc()
	if err != nil || result == nil {
		data.SetCmdStatusError()
		log.ErrorF("Change status Failed, [%s:%s] => %s, Error: %v",
			info.Bucket, info.Key, statusDesc, err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Change status Success, [%s:%s] => %s",
			info.Bucket, info.Key, statusDesc)
	} else {
		data.SetCmdStatusError()
		log.ErrorF("Change status Failed, [%s:%s] => %s, Code:%d, Error:%s",
			info.Bucket, info.Key, statusDesc, result.Code, result.Error)
	}
}

type BatchChangeStatusInfo struct {
	BatchInfo   batch.Info
	Bucket      string
	UnForbidden bool
}

func (info *BatchChangeStatusInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func (info *BatchChangeStatusInfo) getStatus() int {
	// 0:启用  1:禁用
	if info.UnForbidden {
		return 0
	} else {
		return 1
	}
}

func BatchChangeStatus(cfg *iqshell.Config, info BatchChangeStatusInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", cfg.CmdCfg.CmdId, info.Bucket, info.BatchInfo.InputFile))
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

	statusInt := info.getStatus()
	batch.NewHandler(info.BatchInfo).
		EmptyOperation(func() flow.Work {
			return &object.ChangeStatusApiInfo{}
		}).
		SetFileExport(exporter).
		ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
			if len(items) > 0 {
				key := items[0]
				return &object.ChangeStatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Status: statusInt,
				}, nil
			}
			return nil, alert.Error("need more than one param", "")
		}).
		OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
			in, ok := (operation).(*object.ChangeStatusApiInfo)
			if !ok {
				data.SetCmdStatusError()
				log.ErrorF("Change status Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
				return
			}
			if result.IsSuccess() {
				log.InfoF("Change status Success, [%s:%s] => '%d'", in.Bucket, in.Key, in.Status)
			} else {
				data.SetCmdStatusError()
				log.ErrorF("Change status Failed, [%s:%s] => %d, Code: %d, Error: %s",
					in.Bucket, in.Key, in.Status, result.Code, result.Error)
			}
		}).
		OnError(func(err *data.CodeError) {
			data.SetCmdStatusError()
			log.ErrorF("batch change status error:%v:", err)
		}).Start()
}
