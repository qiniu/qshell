package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
)

type ForbiddenInfo struct {
	Bucket      string
	Key         string
	UnForbidden bool
}

func (info *ForbiddenInfo) Check() error {
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

	result, err := object.ChangeStatus(object.ChangeStatusApiInfo{
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

func (info *BatchChangeStatusInfo) Check() error {
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
			key, status := items[0], items[1]
			statusInt, err := strconv.Atoi(status)
			if err != nil {
				log.ErrorF("parse status error:", err)
			} else if key != "" && status != "" {
				return object.ChangeStatusApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Status: statusInt,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		in, ok := (operation).(object.ChangeStatusApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%d\t%d\t%s", in.Key, in.Status, result.Code, result.Error)
			log.ErrorF("Change status Failed, [%s:%s] => %d, Code: %d, Error: %s",
				in.Bucket, in.Key, in.Status, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%d", in.Key, in.Status)
			log.InfoF("Change status Success, [%s:%s] => '%d'", in.Bucket, in.Key, in.Status)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch change status error:%v:", err)
	}).Start()
}
