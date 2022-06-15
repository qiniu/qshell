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
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
	"strconv"
)

type DeleteInfo struct {
	Bucket string
	Key    string
}

func (info *DeleteInfo) Check() *data.CodeError {
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

	result, err := object.Delete(&object.DeleteApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		DeleteAfterDays: 0,
	})

	if err != nil || result == nil {
		log.ErrorF("Delete Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Delete Success, [%s:%s]", info.Bucket, info.Key)
	} else {
		log.ErrorF("Delete Failed, [%s:%s], Code:%d, Error:%s",
			info.Bucket, info.Key, result.Code, result.Error)
	}
}

type BatchDeleteInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchDeleteInfo) Check() *data.CodeError {
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

	lineParser := bucket.NewListLineParser()
	batch.NewHandler(info.BatchInfo).EmptyOperation(func() flow.Work {
		return &object.DeleteApiInfo{}
	}).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		listObject, e := lineParser.Parse(items)
		if e != nil {
			return nil, e
		}

		if len(listObject.Key) == 0 {
			return nil, alert.Error("key invalid", "")
		}

		return &object.DeleteApiInfo{
			Bucket: info.Bucket,
			Key:    listObject.Key,
			Condition: batch.OperationCondition{
				PutTime: listObject.PutTimeString(),
			},
		}, nil
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.DeleteApiInfo)
		if !ok {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Delete Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
			return
		}
		if result.IsSuccess() {
			exporter.Success().Export(operationInfo)
			if len(apiInfo.Condition.PutTime) == 0 {
				log.InfoF("Delete Success, [%s:%s]", apiInfo.Bucket, apiInfo.Key)
			} else {
				log.InfoF("Delete Success, [%s:%s], PutTime:'%s'", apiInfo.Bucket, apiInfo.Key, apiInfo.Condition.PutTime)
			}
		} else {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			if len(apiInfo.Condition.PutTime) == 0 {
				log.ErrorF("Delete Failed, [%s:%s], Code: %d, Error: %s",
					apiInfo.Bucket, apiInfo.Key, result.Code, result.Error)
			} else {
				log.ErrorF("Delete Failed, [%s:%s], PutTime:'%s', Code: %d, Error: %s",
					apiInfo.Bucket, apiInfo.Key, apiInfo.Condition.PutTime, result.Code, result.Error)
			}
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch delete error:%v:", err)
	}).Start()
}

type DeleteAfterInfo struct {
	Bucket    string
	Key       string
	AfterDays string
}

func (info *DeleteAfterInfo) Check() *data.CodeError {
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

func getAfterDaysOfInt(after string) (int, *data.CodeError) {
	if len(after) == 0 {
		return 0, nil
	}
	i, err := strconv.Atoi(after)
	return i, data.ConvertError(err)
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

	result, err := object.Delete(&object.DeleteApiInfo{
		Bucket:          info.Bucket,
		Key:             info.Key,
		DeleteAfterDays: afterDays,
	})

	if err != nil || result == nil {
		log.ErrorF("Expire Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if result.IsSuccess() {
		log.InfoF("Expire Success, [%s:%s], '%s'天后删除", info.Bucket, info.Key, info.AfterDays)
	} else {
		log.ErrorF("Expire Failed, [%s:%s], Code:%d, Error:%s",
			info.Bucket, info.Key, result.Code, result.Error)
	}
}

func BatchDeleteAfter(cfg *iqshell.Config, info BatchDeleteInfo) {
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
		return &object.DeleteApiInfo{}
	}).ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
		after := ""
		key := items[0]
		if len(key) == 0 {
			return nil, alert.Error("key invalid", "")
		}

		if len(items) > 1 {
			after = items[1]
		}
		afterDays, err := getAfterDaysOfInt(after)
		if err != nil {
			return nil, data.NewEmptyError().AppendDescF("parse after days error:%v from:%s", err, after)
		}

		return &object.DeleteApiInfo{
			Bucket:          info.Bucket,
			Key:             key,
			DeleteAfterDays: afterDays,
		}, nil
	}).OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
		apiInfo, ok := (operation).(*object.DeleteApiInfo)
		if !ok {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Delete Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
			return
		}
		if result.IsSuccess() {
			exporter.Success().Export(operationInfo)
			log.InfoF("Expire Success, [%s:%s], '%d'天后删除", apiInfo.Bucket, apiInfo.Key, apiInfo.DeleteAfterDays)
		} else {
			exporter.Fail().ExportF("%s%s%d-%s", operationInfo, flow.ErrorSeparate, result.Code, result.Error)
			log.ErrorF("Expire Failed, [%s:%s], '%d'天后删除, Code: %d, Error: %s", apiInfo.Bucket, apiInfo.Key, apiInfo.DeleteAfterDays, result.Code, result.Error)
		}
	}).OnError(func(err *data.CodeError) {
		log.ErrorF("Batch expire error:%v:", err)
	}).Start()
}
