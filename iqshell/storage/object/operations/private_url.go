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
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultDeadline = 3600
)

type PrivateUrlInfo struct {
	PublicUrl string
	Deadline  string
}

func (p PrivateUrlInfo) WorkId() string {
	return p.PublicUrl
}

func (p *PrivateUrlInfo) Check() *data.CodeError {
	if len(p.PublicUrl) == 0 {
		return alert.CannotEmptyError("PublicUrl", "")
	}
	return nil
}

func (p PrivateUrlInfo) getDeadlineOfInt() (int64, *data.CodeError) {
	if len(p.Deadline) == 0 {
		return time.Now().Add(time.Second * DefaultDeadline).Unix(), nil
	}

	if val, err := strconv.ParseInt(p.Deadline, 10, 64); err != nil {
		return 0, data.NewEmptyError().AppendDesc("invalid deadline")
	} else {
		return val, nil
	}
}

func PrivateUrl(cfg *iqshell.Config, info PrivateUrlInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	deadline, err := info.getDeadlineOfInt()
	if err != nil {
		log.Error(err)
		return
	}

	url, err := download.PublicUrlToPrivate(download.PublicUrlToPrivateApiInfo{
		PublicUrl: info.PublicUrl,
		Deadline:  deadline,
	})

	log.Alert(url)
}

type BatchPrivateUrlInfo struct {
	BatchInfo batch.Info
	Deadline  string
}

func (info *BatchPrivateUrlInfo) Check() *data.CodeError {
	return nil
}

func BatchPrivateUrl(cfg *iqshell.Config, info BatchPrivateUrlInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", cfg.CmdCfg.CmdId, info.Deadline, info.BatchInfo.InputFile))
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

	dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
	if info.BatchInfo.EnableRecord {
		log.DebugF("batch sign recorder:%s", dbPath)
	} else {
		log.Debug("batch sign recorder:Not Enable")
	}

	metric := &batch.Metric{}
	metric.Start()
	flow.New(info.BatchInfo.Info).
		WorkProviderWithFile(info.BatchInfo.InputFile,
			info.BatchInfo.EnableStdin,
			flow.NewItemsWorkCreator(info.BatchInfo.ItemSeparate, 1, func(items []string) (work flow.Work, err *data.CodeError) {
				url := items[0]
				if url == "" {
					return nil, alert.Error("url invalid", "")
				}

				urlToSign := strings.TrimSpace(url)
				if urlToSign == "" {
					return nil, alert.Error("url invalid after TrimSpace", "")
				}
				return &PrivateUrlInfo{
					PublicUrl: url,
					Deadline:  info.Deadline,
				}, nil
			})).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*PrivateUrlInfo)
				if deadline, gErr := in.getDeadlineOfInt(); gErr == nil {
					if r, pErr := download.PublicUrlToPrivate(download.PublicUrlToPrivateApiInfo{
						PublicUrl: in.PublicUrl,
						Deadline:  deadline,
					}); pErr != nil {
						return nil, pErr
					} else {
						return r, nil
					}
				} else {
					return nil, gErr
				}
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
					Work: &PrivateUrlInfo{},
				},
				Result: &download.PublicUrlToPrivateApiResult{},
				Err:    nil,
			}
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if !info.BatchInfo.RecordRedoWhileError {
				return false, nil
			}

			if workRecord.Err != nil {
				return true, workRecord.Err
			}
			result, _ := workRecord.Result.(*download.PublicUrlToPrivateApiResult)
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

			operationResult, _ := result.(*download.PublicUrlToPrivateApiResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					exporter.Success().Export(work.Data)
					log.DebugF("Skip line:%s because have done and success", work.Data)
				} else {
					metric.AddFailureCount(1)
					exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
					log.DebugF("Skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
				log.DebugF("Skip line:%s because:%v", work.Data, err)
			}

		}).
		OnWorkSuccess(func(work *flow.WorkInfo, result flow.Result) {
			metric.AddCurrentCount(1)
			metric.AddSuccessCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			r, _ := result.(*download.PublicUrlToPrivateApiResult)
			exporter.Success().Export(work.Data)
			log.Alert(r.Url)
		}).
		OnWorkFail(func(work *flow.WorkInfo, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.AddFailureCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
			log.Error(err)
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save batch sign result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save batch sign result to path:%s", resultPath)
	}

	log.Info("\n--------------- Batch Sign Result ---------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("-------------------------------------------------")
}
