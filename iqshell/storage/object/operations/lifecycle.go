package operations

import (
	"fmt"
	"path/filepath"

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
)

type ChangeLifecycleInfo object.ChangeLifecycleApiInfo

func (info *ChangeLifecycleInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}

	if info.ToIAAfterDays == 0 &&
		info.ToArchiveIRAfterDays == 0 &&
		info.ToArchiveAfterDays == 0 &&
		info.ToDeepArchiveAfterDays == 0 &&
		info.DeleteAfterDays == 0 {
		return data.NewEmptyError().AppendDesc("must set at least one value of lifecycle")
	}

	return nil
}

func ChangeLifecycle(cfg *iqshell.Config, info *ChangeLifecycleInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: info,
	}); !shouldContinue {
		return
	}

	result, err := object.ChangeLifecycle((*object.ChangeLifecycleApiInfo)(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("ChangeLifecycle Failed, [%s:%s], Error:%v",
			info.Bucket, info.Key, err)
		return
	}

	if result.IsSuccess() {
		lifecycleValues := []int{
			info.ToIAAfterDays, info.ToIntelligentTieringAfterDays, info.ToArchiveIRAfterDays, info.ToArchiveAfterDays,
			info.ToDeepArchiveAfterDays, info.DeleteAfterDays,
		}
		lifecycleDescs := []string{
			"to IA storage", "to IntelligentTiering storage", "to ARCHIVE_IR storage", "to ARCHIVE storage",
			"to DEEP_ARCHIVE storage", "delete",
		}
		log.InfoF("Change lifecycle Success, [%s:%s]", info.Bucket, info.Key)
		for i := 0; i < len(lifecycleValues); i++ {
			lifecycleValue := lifecycleValues[i]
			lifecycleDesc := lifecycleDescs[i]
			if lifecycleValue == 0 {
				continue
			}
			if lifecycleValue == -1 {
				log.InfoF("● cancel %s", lifecycleDesc)
			} else {
				log.InfoF("● %s after %d days", lifecycleDesc, lifecycleValue)
			}
		}
	} else {
		data.SetCmdStatusError()
		log.ErrorF("Change lifecycle Failed, [%s:%s], Code:%d, Error:%s",
			info.Bucket, info.Key, result.Code, result.Error)
	}
}

type BatchChangeLifecycleInfo struct {
	BatchInfo                     batch.Info //
	Bucket                        string     //
	ToIAAfterDays                 int        // 转换到 低频存储类型，设置为 -1 表示取消
	ToArchiveIRAfterDays          int        // 转换到 归档直读存储类型， 设置为 -1 表示取消
	ToArchiveAfterDays            int        // 转换到 归档存储类型， 设置为 -1 表示取消
	ToDeepArchiveAfterDays        int        // 转换到 深度归档存储类型， 设置为 -1 表示取消
	ToIntelligentTieringAfterDays int        // 转换到 智能分层存储类型， 设置为 -1 表示取消
	DeleteAfterDays               int        // 过期删除，删除后不可恢复，设置为 -1 表示取消
}

func (info *BatchChangeLifecycleInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if info.ToIAAfterDays == 0 &&
		info.ToArchiveIRAfterDays == 0 &&
		info.ToArchiveAfterDays == 0 &&
		info.ToDeepArchiveAfterDays == 0 &&
		info.ToIntelligentTieringAfterDays == 0 &&
		info.DeleteAfterDays == 0 {
		return data.NewEmptyError().AppendDesc("must set at least one value of lifecycle")
	}

	return nil
}

func BatchChangeLifecycle(cfg *iqshell.Config, info BatchChangeLifecycleInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%d:%d:%d:%d:%d:%d:%s", cfg.CmdCfg.CmdId, info.Bucket,
			info.ToIAAfterDays, info.ToArchiveIRAfterDays, info.ToArchiveAfterDays,
			info.ToDeepArchiveAfterDays, info.ToIntelligentTieringAfterDays, info.DeleteAfterDays,
			info.BatchInfo.InputFile))
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
		data.SetCmdStatusError()
		return
	}

	lineParser := bucket.NewListLineParser()
	batch.NewHandler(info.BatchInfo).
		EmptyOperation(func() flow.Work {
			return &object.ChangeLifecycleApiInfo{}
		}).
		SetFileExport(exporter).
		ItemsToOperation(func(items []string) (operation batch.Operation, err *data.CodeError) {
			listObject, e := lineParser.Parse(items)
			if e != nil {
				return nil, e
			}
			if len(listObject.Key) == 0 {
				return nil, alert.Error("key invalid", "")
			}
			return &object.ChangeLifecycleApiInfo{
				Bucket:                        info.Bucket,
				Key:                           listObject.Key,
				ToIAAfterDays:                 info.ToIAAfterDays,
				ToArchiveIRAfterDays:          info.ToArchiveIRAfterDays,
				ToArchiveAfterDays:            info.ToArchiveAfterDays,
				ToDeepArchiveAfterDays:        info.ToDeepArchiveAfterDays,
				ToIntelligentTieringAfterDays: info.ToIntelligentTieringAfterDays,
				DeleteAfterDays:               info.DeleteAfterDays,
			}, nil
		}).
		OnResult(func(operationInfo string, operation batch.Operation, result *batch.OperationResult) {
			apiInfo, ok := (operation).(*object.ChangeLifecycleApiInfo)
			if !ok {
				data.SetCmdStatusError()
				log.ErrorF("Change lifecycle Failed, %s, Code: %d, Error: %s", operationInfo, result.Code, result.Error)
				return
			}
			in := (*ChangeLifecycleInfo)(apiInfo)
			if result.IsSuccess() {
				log.InfoF("Change lifecycle Success, [%s:%s] => '%d:%d:%d:%d:%d'", in.Bucket, in.Key,
					in.ToIAAfterDays, in.ToArchiveIRAfterDays, in.ToArchiveAfterDays, in.ToDeepArchiveAfterDays, in.DeleteAfterDays)
			} else {
				data.SetCmdStatusError()
				log.ErrorF("Change lifecycle Failed, [%s:%s], Code: %d, Error: %s", in.Bucket, in.Key, result.Code, result.Error)
			}
		}).
		OnError(func(err *data.CodeError) {
			data.SetCmdStatusError()
			log.ErrorF("Batch Change lifecycle error:%v", err)
		}).Start()
}
