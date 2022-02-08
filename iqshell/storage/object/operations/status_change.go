package operations

import (
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

func (c ForbiddenInfo) getStatus() int {
	// 0:启用  1:禁用
	if c.UnForbidden {
		return 0
	} else {
		return 1
	}
}

func ForbiddenObject(info ForbiddenInfo) {
	result, err := object.ChangeStatus(object.ChangeStatusApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Status: info.getStatus(),
	})

	if err != nil {
		log.ErrorF("change stat error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("change stat error:%s", result.Error)
		return
	}
}

type BatchChangeStatusInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func BatchChangeStatus(info BatchChangeStatusInfo) {
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
			handler.Export().Fail().ExportF("%s\t%d\t%d\t%s\n", in.Key, in.Status, result.Code, result.Error)
			log.ErrorF("Change status '%s' => '%s' Failed, Code: %d, Error: %s",
				in.Key, in.Status, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%d\n", in.Key, in.Status)
			log.ErrorF("Change status '%s' => '%d' success\n", in.Key, in.Status)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch change status error:%v:", err)
	}).Start()
}
