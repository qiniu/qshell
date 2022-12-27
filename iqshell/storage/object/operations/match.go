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
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"path/filepath"
)

type MatchInfo object.MatchApiInfo

func (info *MatchInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	if len(info.LocalFile) == 0 {
		return alert.CannotEmptyError("LocalFile", "")
	}
	return nil
}

func Match(cfg *iqshell.Config, info MatchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	_, err := object.Match(object.MatchApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Match  Failed, [%s:%s] => '%s', Error:%v",
			info.Bucket, info.Key, info.LocalFile, err)
	} else {
		log.InfoF("Match Success, [%s:%s] => '%s'",
			info.Bucket, info.Key, info.LocalFile)
	}
}

type BatchMatchInfo struct {
	BatchInfo    batch.Info
	Bucket       string
	LocalFileDir string
}

func (info *BatchMatchInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if len(info.LocalFileDir) == 0 {
		return alert.CannotEmptyError("LocalFileDir", "")
	}

	if path, e := filepath.Abs(info.LocalFileDir); e != nil || len(path) == 0 {
		return data.NewEmptyError().AppendDescF("LocalFileDir invalid, err:%s", e)
	} else {
		info.LocalFileDir = path
		log.DebugF("LocalFileDir:%s", info.LocalFileDir)
	}
	return nil
}

func BatchMatch(cfg *iqshell.Config, info BatchMatchInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", cfg.CmdCfg.CmdId, info.Bucket, info.LocalFileDir, info.BatchInfo.InputFile))
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

	dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
	if info.BatchInfo.EnableRecord {
		log.DebugF("batch match recorder:%s", dbPath)
	} else {
		log.Debug("batch match recorder:Not Enable")
	}

	metric := &batch.Metric{}
	metric.Start()
	lineParser := bucket.NewListLineParser()
	flow.New(info.BatchInfo.Info).
		WorkProviderWithFile(info.BatchInfo.InputFile,
			info.BatchInfo.EnableStdin,
			flow.NewItemsWorkCreator(info.BatchInfo.ItemSeparate, 1, func(items []string) (work flow.Work, err *data.CodeError) {
				listObject, e := lineParser.Parse(items)
				if e != nil {
					return nil, e
				}

				if len(listObject.Key) == 0 {
					return nil, alert.Error("key invalid", "")
				}

				return &object.MatchApiInfo{
					Bucket:         info.Bucket,
					Key:            listObject.Key,
					ServerFileHash: listObject.Hash,
					LocalFile:      filepath.Join(info.LocalFileDir, listObject.Key),
				}, nil
			})).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*object.MatchApiInfo)
				return object.Match(*in)
			}), nil
		})).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		SetOverseerEnable(info.BatchInfo.EnableRecord).
		SetDBOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: &object.MatchApiInfo{},
				},
				Result: &object.MatchResult{},
				Err:    nil,
			}
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err == nil {
				return false, nil
			}

			if !info.BatchInfo.RecordRedoWhileError {
				return false, workRecord.Err
			}

			result, _ := workRecord.Result.(*object.MatchResult)
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
			metric.PrintProgress("Batching:" + work.Data)

			operationResult, _ := result.(*object.MatchResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					exporter.Success().ExportF("%s", work.Data)
					log.InfoF("Skip line:%s because have done and success", work.Data)
				} else {
					metric.AddFailureCount(1)
					exporter.Fail().ExportF("%s", work.Data)
					log.InfoF("Skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
				log.InfoF("Skip line:%s because:%v", work.Data, err)
			}
		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			metric.AddCurrentCount(1)
			metric.AddSuccessCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			in, _ := workInfo.Work.(*object.MatchApiInfo)
			exporter.Success().ExportF("%s\t \t%s", in.Key, in.ServerFileHash)
			log.InfoF("Match Success, [%s:%s] => '%s'", info.Bucket, in.Key, in.LocalFile)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.AddFailureCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
			if in, ok := workInfo.Work.(*object.MatchApiInfo); ok {
				log.ErrorF("Match Failed, [%s:%s] => '%s', Error: %s", info.Bucket, in.Key, in.LocalFile, err)
			} else {
				log.ErrorF("Match Failed, %s, Error: %s", workInfo.Data, err)
			}
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	log.InfoF("job dir:%s, there is a cache related to this command in this folder, which will also be used next time the same command is executed. If you are sure that you don’t need it, you can delete this folder.", workspace.GetJobDir())

	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		data.SetCmdStatusError()
		log.ErrorF("save batch match result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save batch match result to path:%s", resultPath)
	}

	log.Info("--------------- Batch Match Result ---------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("--------------------------------------------------")

	if !metric.IsCompletedSuccessfully() {
		data.SetCmdStatusError()
	}
}
