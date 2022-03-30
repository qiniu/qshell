package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
)

type ChangeTypeInfo struct {
	Bucket string
	Key    string
	Type   string
}

func (info *ChangeTypeInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.Type) == 0 {
		return alert.CannotEmptyError("Type", "")
	}
	return nil
}

func (info *ChangeTypeInfo) getTypeOfInt() (int, error) {
	if len(info.Type) == 0 {
		return -1, errors.New(alert.CannotEmpty("type", ""))
	}

	ret, err := strconv.Atoi(info.Type)
	if err != nil {
		return -1, errors.New("Parse type error:" + err.Error())
	}

	if ret < 0 || ret > 3 {
		return -1, errors.New("type must be one of 0, 1, 2, 3")
	}
	return ret, nil
}

func ChangeType(cfg *iqshell.Config, info ChangeTypeInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	t, err := info.getTypeOfInt()
	if err != nil {
		log.ErrorF("Change Type Failed, [%s:%s] error:%v", err)
		return
	}

	result, err := object.ChangeType(object.ChangeTypeApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Type:   t,
	})

	if err != nil {
		log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Error: %v",
			info.Bucket, info.Key, t, getStorageTypeDescription(t), err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Code: %d, Error: %s",
			info.Bucket, info.Key, t, getStorageTypeDescription(t), result.Code, result.Error)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Change Type Success, [%s:%s] => '%d'(%s)", info.Bucket, info.Key, t, getStorageTypeDescription(t))
	}
}

type BatchChangeTypeInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchChangeTypeInfo) Check() error {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchChangeType(cfg *iqshell.Config, info BatchChangeTypeInfo) {
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
			key, t := items[0], items[1]
			tInt, err := strconv.Atoi(t)
			if err != nil {
				log.ErrorF("Parse type error:%v", err)
			} else if len(key) > 0 && len(t) > 0 {
				return object.ChangeTypeApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Type:   tInt,
				}, true
			}
		}
		return nil, true
	}).OnResult(func(operation batch.Operation, result batch.OperationResult) {
		in, ok := (operation).(object.ChangeTypeApiInfo)
		if !ok {
			return
		}
		if result.Code != 200 || result.Error != "" {
			handler.Export().Fail().ExportF("%s\t%d\t%d\t%s", in.Key, in.Type, result.Code, result.Error)
			log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Code: %d, Error: %s",
				info.Bucket, in.Key, in.Type, getStorageTypeDescription(in.Type), result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%d", in.Key, in.Type)
			log.InfoF("Change Type Success, [%s:%s] => '%d'(%s) ",
				info.Bucket, in.Key, in.Type, getStorageTypeDescription(in.Type))
		}
	}).OnError(func(err error) {
		log.ErrorF("Batch change Type error:%v:", err)
	}).Start()
}
