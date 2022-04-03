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
	"strconv"
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
	if err != nil {
		log.ErrorF("Change status Failed, [%s:%s] => %s, Error: %v",
			info.Bucket, info.Key, info.getStatus(), statusDesc)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("Change status Failed, [%s:%s] => %s, Code:%s, Error: %s",
			info.Bucket, info.Key, statusDesc, result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.Info("Change status Success, [%s:%s] => %s",
			info.Bucket, info.Key, statusDesc)
	}
}

type BatchChangeStatusInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchChangeStatusInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchChangeStatus(cfg *iqshell.Config, info BatchChangeStatusInfo) {
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
			key, status := items[0], items[1]
			statusInt, e := strconv.Atoi(status)
			if e != nil {
				return nil, data.NewEmptyError().AppendDescF("parse status error:%v", e)
			} else if key != "" && status != "" {
				return &object.ChangeStatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Status: statusInt,
				}, nil
			}
		}
		return nil, alert.Error("need more than one param", "")
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		in, ok := (operation).(*object.ChangeStatusApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Change status Failed, [%s:%s] => %d, Code: %d, Error: %s",
				in.Bucket, in.Key, in.Status, result.Code, result.Error)
		} else {
			exporter.Success().Export(operationInfo)
			log.InfoF("Change status Success, [%s:%s] => '%d'", in.Bucket, in.Key, in.Status)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("batch change status error:%v:", err)
	}).Start()
}
