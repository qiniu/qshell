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

func convertFreezeAfterDaysToInt(freezeAfterDays string) (int, error) {
	if len(freezeAfterDays) == 0 {
		return 0, alert.CannotEmptyError("FreezeAfterDays", "")
	}

	if freezeAfterDaysInt, err := strconv.Atoi(freezeAfterDays); err != nil {
		return 0, alert.Error("FreezeAfterDays is invalid:"+err.Error(), "")
	} else {
		if freezeAfterDaysInt > 0 || freezeAfterDaysInt < 8 {
			return freezeAfterDaysInt, nil
		}
		return 0, alert.Error("FreezeAfterDays must between 1 and 7, include 1 and 7", "")
	}
}

type RestoreArchiveInfo struct {
	Bucket             string
	Key                string
	FreezeAfterDays    string
	freezeAfterDaysInt int
}

func (info *RestoreArchiveInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.FreezeAfterDays) == 0 {
		return alert.CannotEmptyError("FreezeAfterDays", "")
	}

	if freezeAfterDaysInt, err := convertFreezeAfterDaysToInt(info.FreezeAfterDays); err != nil {
		return err
	} else {
		info.freezeAfterDaysInt = freezeAfterDaysInt
	}

	return nil
}

func RestoreArchive(cfg *iqshell.Config, info RestoreArchiveInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.RestoreArchive(object.RestoreArchiveApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		FreezeAfterDays: info.freezeAfterDaysInt,
	})
	if err != nil {
		log.ErrorF("Restore archive error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Restore archive error:%v", result.Error)
		return
	}

	log.InfoF("Restore archive [%s:%s] FreezeAfterDays:%s success", info.Bucket, info.Key, info.FreezeAfterDays)
}

type BatchRestoreArchiveInfo struct {
	BatchInfo          batch.Info
	Bucket             string
	FreezeAfterDays    string
	freezeAfterDaysInt int
}

func (info *BatchRestoreArchiveInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if freezeAfterDaysInt, err := convertFreezeAfterDaysToInt(info.FreezeAfterDays); err != nil {
		return err
	} else {
		info.freezeAfterDaysInt = freezeAfterDaysInt
	}
	return nil
}

func BatchRestoreArchive(cfg *iqshell.Config, info BatchRestoreArchiveInfo) {
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
			if len(key) > 0 {
				return object.RestoreArchiveApiInfo{
					Bucket:          info.Bucket,
					Key:             key,
					FreezeAfterDays: info.freezeAfterDaysInt,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		apiInfo, ok := (operation).(object.RestoreArchiveApiInfo)
		if !ok {
			return
		}

		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%d\t%s", apiInfo.Key, result.Code, result.Error)
			log.ErrorF("Restore archive [%s:%s] FreezeAfterDays:%s Failed, Code: %d, Error: %s",
				apiInfo.Bucket, apiInfo.Key, apiInfo.FreezeAfterDays,
				result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s", apiInfo.Key)
			log.InfoF("Restore archive [%s:%s] FreezeAfterDays:%s success", apiInfo.Bucket, apiInfo.Key, info.FreezeAfterDays)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch move error:%v:", err)
	}).Start()
}
