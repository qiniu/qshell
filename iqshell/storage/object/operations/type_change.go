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
	"strconv"
)

type ChangeTypeInfo struct {
	Bucket string
	Key    string
	Type   string
}

func (info *ChangeTypeInfo) Check() *data.CodeError {
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

func (info *ChangeTypeInfo) getTypeOfInt() (int, *data.CodeError) {
	if len(info.Type) == 0 {
		return -1, data.NewEmptyError().AppendDesc(alert.CannotEmpty("type", ""))
	}

	ret, err := strconv.Atoi(info.Type)
	if err != nil {
		return -1, data.NewEmptyError().AppendDesc("Parse type error:" + err.Error())
	}

	if ret < 0 || ret > 3 {
		return -1, data.NewEmptyError().AppendDesc("type must be one of 0, 1, 2, 3")
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
		log.ErrorF("Change Type Failed, [%s:%s] error:%v", info.Bucket, info.Key, err)
		return
	}

	result, err := object.ChangeType(&object.ChangeTypeApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Type:   t,
	})

	if err != nil || result == nil {
		log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Error: %v",
			info.Bucket, info.Key, t, getStorageTypeDescription(t), err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Change Type Success, [%s:%s] => '%d'(%s)", info.Bucket, info.Key, t, getStorageTypeDescription(t))
	} else {
		log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Code: %d, Error: %s",
			info.Bucket, info.Key, t, getStorageTypeDescription(t), result.Code, result.Error)
	}
}

type BatchChangeTypeInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchChangeTypeInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchChangeType(cfg *iqshell.Config, info BatchChangeTypeInfo) {
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
		return
	}

	batch.NewHandler(info.BatchInfo).EmptyOperation(func() flow.Work {
		return &object.ChangeTypeApiInfo{}
	}).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		if len(items) > 1 {
			key, t := items[0], items[1]
			if tInt, e := strconv.Atoi(t); e != nil {
				return nil, data.NewEmptyError().AppendDescF("parse type error:%v", e)
			} else if len(key) > 0 && len(t) > 0 {
				return &object.ChangeTypeApiInfo{
					Bucket: info.Bucket,
					Key:    key,
					Type:   tInt,
				}, nil
			}
		}
		return nil, alert.Error("need more than one param", "")
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		in, ok := (operation).(*object.ChangeTypeApiInfo)
		if !ok {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Change status Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
			return
		}
		if result.IsSuccess() {
			exporter.Success().Export(operationInfo)
			log.InfoF("Change Type Success, [%s:%s] => '%d'(%s) ",
				info.Bucket, in.Key, in.Type, getStorageTypeDescription(in.Type))
		} else {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Change Type Failed, [%s:%s] => '%d'(%s), Code: %d, Error: %s",
				info.Bucket, in.Key, in.Type, getStorageTypeDescription(in.Type), result.Code, result.Error)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch change Type error:%v:", err)
	}).Start()
}
