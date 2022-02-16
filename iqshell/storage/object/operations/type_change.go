package operations

import (
	"errors"
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
		return -1, errors.New("parse type error:" + err.Error())
	}

	if ret < 0 || ret > 2 {
		return -1, errors.New("type must be 0 / 1 or 2")
	}
	return ret, nil
}

func ChangeType(info ChangeTypeInfo) {
	t, err := info.getTypeOfInt()
	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	result, err := object.ChangeType(object.ChangeTypeApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Type:   t,
	})

	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Type:%v", result.Error)
		return
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
		return alert.CannotEmptyError("bucket", "")
	}
	return nil
}

func BatchChangeType(info BatchChangeTypeInfo) {
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
				log.ErrorF("parse type error:%v", err)
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
			handler.Export().Fail().ExportF("%s\t%d\t%d\t%s\n", in.Key, in.Type, result.Code, result.Error)
			log.ErrorF("Change Type '%s' => '%d' Failed, Code: %d, Error: %s",
				in.Key, in.Type, result.Code, result.Error)
		} else {
			handler.Export().Success().ExportF("%s\t%d\n", in.Key, in.Type)
			log.InfoF("Change Type '%s' => '%d' success", in.Key, in.Type)
		}
	}).OnError(func(err error) {
		log.ErrorF("batch change Type error:%v:", err)
	}).Start()
}
