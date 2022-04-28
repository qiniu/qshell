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
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
)

type FetchInfo object.FetchApiInfo

func (info *FetchInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.FromUrl) == 0 {
		return alert.CannotEmptyError("RemoteResourceUrl", "")
	}
	return nil
}

func Fetch(cfg *iqshell.Config, info FetchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Fetch(object.FetchApiInfo(info))
	if err != nil {
		log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error:%v",
			info.FromUrl, info.Bucket, info.Key, err)
	} else {
		log.InfoF("Fetch Success, '%s' => [%s:%s]", info.FromUrl, info.Bucket, info.Key)
		log.AlertF("Key:%s", result.Key)
		log.AlertF("FileHash:%s", result.Hash)
		log.AlertF("Fsize: %d (%s)", result.Fsize, utils.FormatFileSize(result.Fsize))
		log.AlertF("Mime:%s", result.MimeType)
	}
}

type BatchFetchInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchFetchInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

//BatchFetch 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchFetch(cfg *iqshell.Config, info BatchFetchInfo) {
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

	// overseer， 数组源 不类型不记录中间状态
	var overseer flow.Overseer
	if info.BatchInfo.EnableRecord {
		dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
		log.DebugF("batch fetch recorder:%s", dbPath)
		if overseer, err = flow.NewDBRecordOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: nil,
				},
				Result: &object.FetchResult{},
				Err:    nil,
			}
		}); err != nil {
			log.ErrorF("batch fetch create overseer error:%v", err)
			return
		}
	} else {
		log.Debug("batch fetch recorder:Not Enable")
	}

	metric := &batch.Metric{}
	metric.Start()
	flow.New(info.BatchInfo.Info).
		WorkProviderWithFile(info.BatchInfo.InputFile,
			info.BatchInfo.EnableStdin,
			flow.NewItemsWorkCreator(info.BatchInfo.ItemSeparate, 1, func(items []string) (work flow.Work, err *data.CodeError) {
				key := ""
				fromUrl := items[0]
				if len(items) > 1 {
					key = items[1]
				} else if k, e := utils.KeyFromUrl(fromUrl); e == nil {
					key = k
				}
				if len(key) == 0 || len(fromUrl) == 0 {
					return nil, alert.Error("key or fromUrl invalid", "")
				}

				return &object.FetchApiInfo{
					Bucket:  info.Bucket,
					Key:     key,
					FromUrl: fromUrl,
				}, nil
			})).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*object.FetchApiInfo)
				return object.Fetch(*in)
			}), nil
		})).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		SetOverseer(overseer).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err == nil {
				return false, nil
			}

			if !info.BatchInfo.RecordRedoWhileError {
				return false, workRecord.Err
			}

			result, _ := workRecord.Result.(*object.FetchResult)
			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDesc("result is invalid")
			}
			return false, nil
		}).
		OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching")

			operationResult, _ := result.(*object.FetchResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					log.DebugF("Skip line:%s because have done and success", work.Data)
				} else {
					metric.AddFailureCount(1)
					log.DebugF("Skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
				log.DebugF("Skip line:%s because:%v", work.Data, err)
			}

		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			metric.AddCurrentCount(1)
			metric.AddSuccessCount(1)
			metric.PrintProgress("Batching")

			in, _ := workInfo.Work.(*object.FetchApiInfo)
			exporter.Success().ExportF("%s\t%s", in.FromUrl, in.Bucket)
			log.InfoF("Fetch Success, '%s' => [%s:%s]", in.FromUrl, info.Bucket, in.Key)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.AddFailureCount(1)
			metric.PrintProgress("Batching")

			exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
			if in, ok := workInfo.Work.(*object.FetchApiInfo); ok {
				log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error: %v", in.FromUrl, in.Bucket, in.Key, err)
			} else {
				log.ErrorF("Fetch Failed, %s, Error: %s", workInfo.Data, err)
			}
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save batch fetch result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save batch fetch result to path:%s", resultPath)
	}

	log.Info("--------------- Batch Result ---------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("--------------------------------------------")
}
