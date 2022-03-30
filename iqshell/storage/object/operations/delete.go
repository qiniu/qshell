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

type DeleteInfo struct {
	Bucket string
	Key    string
}

func (info *DeleteInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func Delete(cfg *iqshell.Config, info DeleteInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Delete(object.DeleteApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		DeleteAfterDays: 0,
	})

	if err != nil {
		log.ErrorF("Delete Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("Delete Failed, [%s:%s], Code:%d, Error:%s",
			info.Bucket, info.Key, result.Code, result.Error)
		return
	}

	log.InfoF("Delete Success, [%s:%s]", info.Bucket, info.Key)
}

type BatchDeleteInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchDeleteInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

// BatchDelete 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDelete(cfg *iqshell.Config, info BatchDeleteInfo) {
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
		key := ""
		putTime := ""

		if len(items) > 0 {
			key = items[0]
		}
		if len(key) == 0 {
			return nil, true
		}

		if len(items) > 1 {
			putTime = items[1]
		}
		// list 结果格式 14902611578248790
		if len(putTime) != 17 && len(items) > 3 && len(items[3]) == 17 {
			putTime = items[3]
		}
		return object.DeleteApiInfo{
			Bucket: info.Bucket,
			Key:    key,
			Condition: batch.OperationCondition{
				PutTime: putTime,
			},
		}, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.DeleteApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s", apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
			if len(apiInfo.Condition.PutTime) == 0 {
				log.ErrorF("Delete Failed, [%s:%s], Code: %d, Error: %s",
					apiInfo.Bucket, apiInfo.Key, result.Code, result.Error)
			} else {
				log.ErrorF("Delete Failed, [%s:%s], PutTime:'%s', Code: %d, Error: %s",
					apiInfo.Bucket, apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
			}
		} else {
			handler.Export().Success().ExportF("%s\t%s", apiInfo.Key, apiInfo.Condition.PutTime)
			if len(apiInfo.Condition.PutTime) == 0 {
				log.InfoF("Delete Success, [%s:%s]", apiInfo.Bucket, apiInfo.Key)
			} else {
				log.InfoF("Delete Success, [%s:%s], PutTime:'%s'", apiInfo.Bucket, apiInfo.Key, apiInfo.Condition.PutTime)
			}
		}
	}).OnError(func(err error) {
		log.ErrorF("Batch delete error:%v:", err)
	}).Start()
}

type DeleteAfterInfo struct {
	Bucket    string
	Key       string
	AfterDays string
}

func (info *DeleteAfterInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.AfterDays) == 0 {
		return alert.CannotEmptyError("DeleteAfterDays", "")
	}
	return nil
}

func getAfterDaysOfInt(after string) (int, error) {
	if len(after) == 0 {
		return 0, nil
	}
	return strconv.Atoi(after)
}

func DeleteAfter(cfg *iqshell.Config, info DeleteAfterInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	afterDays, err := getAfterDaysOfInt(info.AfterDays)
	if err != nil {
		log.ErrorF("DeleteAfterDays invalid:%v", err)
		return
	}

	result, err := object.Delete(object.DeleteApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		DeleteAfterDays: afterDays,
	})

	if err != nil {
		log.ErrorF("Expire Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("Expire Failed, [%s:%s], Code:%d, Error:%s",
			info.Bucket, info.Key, result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Expire Success, [%s:%s], '%s'天后删除", info.Bucket, info.Key, info.AfterDays)
	}
}

// BatchDeleteAfter 延迟批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchDeleteAfter(cfg *iqshell.Config, info BatchDeleteInfo) {
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
					Bucket:          info.Bucket,
					Key:             key,
					DeleteAfterDays: afterDays,
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
			handler.Export().Fail().ExportF("%s\t%s\t%d\t%s", apiInfo.Key, apiInfo.DeleteAfterDays, result.Code, result.Error)
			log.ErrorF("Expire Failed, [%s:%s], '%d'天后删除, Code: %d, Error: %s", apiInfo.Bucket, apiInfo.Key, apiInfo.DeleteAfterDays, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%s", apiInfo.Key, apiInfo.DeleteAfterDays)
			log.InfoF("Expire Success, [%s:%s], '%d'天后删除", apiInfo.Bucket, apiInfo.Key, apiInfo.DeleteAfterDays)
		}
	}).OnError(func(err error) {
		log.ErrorF("Batch expire error:%v:", err)
	}).Start()
}
