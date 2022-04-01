package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
)

func convertFreezeAfterDaysToInt(freezeAfterDays string) (int, *data.CodeError) {
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

func (info *RestoreArchiveInfo) Check() *data.CodeError {
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

	result, err := object.RestoreArchive(&object.RestoreArchiveApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		FreezeAfterDays: info.freezeAfterDaysInt,
	})
	if err != nil {
		log.ErrorF("Restore archive Failed, [%s:%s], FreezeAfterDays:%s, Error: %v",
			info.Bucket, info.Key, info.FreezeAfterDays, err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Restore archive Failed, [%s:%s], FreezeAfterDays:%s, Code: %d, Error: %s",
			info.Bucket, info.Key, info.FreezeAfterDays,
			result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Restore archive Success, [%s:%s], FreezeAfterDays:%s",
			info.Bucket, info.Key, info.FreezeAfterDays)
	}
}

type BatchRestoreArchiveInfo struct {
	BatchInfo          batch.Info
	Bucket             string
	FreezeAfterDays    string
	freezeAfterDaysInt int
}

func (info *BatchRestoreArchiveInfo) Check() *data.CodeError {
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

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}

	batch.NewHandler(info.BatchInfo).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		key := items[0]
		if len(key) > 0 {
			return &object.RestoreArchiveApiInfo{
				Bucket:          info.Bucket,
				Key:             key,
				FreezeAfterDays: info.freezeAfterDaysInt,
			}, nil
		}
		return nil, alert.Error("key invalid", "")
	}).OnResult(func(operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.RestoreArchiveApiInfo)
		if !ok {
			return
		}

		if result.Code != 200 || result.Error != "" {
			exporter.Fail().ExportF("%s\t%d\t%s", apiInfo.Key, result.Code, result.Error)
			log.ErrorF("Restore archive Failed, [%s:%s], FreezeAfterDays:%d, Code: %d, Error: %s",
				apiInfo.Bucket, apiInfo.Key, apiInfo.FreezeAfterDays,
				result.Code, result.Error)
		} else {
			exporter.Success().ExportF("%s", apiInfo.Key)
			log.InfoF("Restore archive Success, [%s:%d], FreezeAfterDays:%s",
				apiInfo.Bucket, apiInfo.Key, apiInfo.FreezeAfterDays)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch restore archive error:%v:", err)
	}).Start()
}
