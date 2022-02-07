package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
)

type DeleteInfo struct {
	Bucket    string
	Key       string
	AfterDays string
}

func getAfterDaysOfInt(after string) (int, error) {
	if len(after) == 0 {
		return -1, nil
	}
	return strconv.Atoi(after)
}

func Delete(info DeleteInfo) {
	afterDays, err := getAfterDaysOfInt(info.AfterDays)
	if err != nil {
		log.ErrorF("delete after days invalid:%v", err)
		return
	}

	result, err := object.Delete(object.DeleteApiInfo{
		Bucket:    info.Bucket,
		Key:       info.Key,
		AfterDays: afterDays,
	})

	if err != nil {
		log.ErrorF("delete error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("delete error:%s", result.Error)
		return
	}
}

type BatchDeleteInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

// BatchDelete 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDelete(info BatchDeleteInfo) {
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
			key := items[0]
			putTime := ""
			if len(items) > 1 {
				putTime = items[1]
			}
			if key != "" {
				return object.DeleteApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Condition: batch.OperationCondition{
						PutTime: putTime,
					},
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.DeleteApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s\n", apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
			log.ErrorF("Delete '%s' when put time:'%d' Failed, Code: %d, Error: %s\n", apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s\n", apiInfo.Key, apiInfo.Condition.PutTime)
			log.AlertF("Delete '%s' when put time:'%d' success\n", apiInfo.Key, apiInfo.Condition.PutTime)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch delete error:%v:", err)
	}).Start()
}

// BatchDeleteAfter 延迟批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDeleteAfter(info BatchDeleteInfo) {
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
			key := items[0]
			after := ""
			if len(items) > 1 {
				after = items[1]
			}
			afterDays, err := getAfterDaysOfInt(after)
			if err != nil {
				log.ErrorF("parse after days error:%v from:%s", err, after)
			} else if key != "" {
				return object.DeleteApiInfo{
					Bucket:    info.Bucket,
					Key:       key,
					AfterDays: afterDays,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.DeleteApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s\n", apiInfo.Key, apiInfo.AfterDays, result.Code, result.Error)
			log.ErrorF("Expire '%s' => '%d' Failed, Code: %d, Error: %s\n", apiInfo.Key, apiInfo.AfterDays, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s\n", apiInfo.Key, apiInfo.AfterDays)
			log.AlertF("Expire '%s' => '%d' success\n", apiInfo.Key, apiInfo.AfterDays)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch delete after error:%v:", err)
	}).Start()
}
