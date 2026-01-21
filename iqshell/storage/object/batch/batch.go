package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/locker"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type Info struct {
	flow.Info
	export.FileExporterConfig

	Overwrite bool // 是否覆盖

	// 工作数据源
	WorkList      []flow.Work // 工作数据源：列表
	InputFile     string      // 工作数据源：文件
	ItemSeparate  string      // 工作数据源：每行元素按分隔符分的分隔符
	MinItemsCount int         // 工作数据源：每行元素最小数量
	EnableStdin   bool        // 工作数据源：stdin, 当 InputFile 不存在时使用 stdin

	EnableRecord             bool // 是否开启 record
	RecordRedoWhileError     bool // 重新执行任务时，如果任务已执行但是失败，则再重新执行一次。
	OperationCountPerRequest int  // 每批操作最大的子任务数
}

func (info *Info) Check() *data.CodeError {
	if err := info.Info.Check(); err != nil {
		return err
	}

	if info.MinItemsCount < 1 {
		info.MinItemsCount = 1
	}

	if info.OperationCountPerRequest <= 0 ||
		info.OperationCountPerRequest > defaultOperationCountPerRequest {
		info.OperationCountPerRequest = defaultOperationCountPerRequest
	}

	if len(info.ItemSeparate) == 0 {
		info.ItemSeparate = "\t"
	}

	return nil
}

type Handler interface {
	EmptyOperation(emptyOperation func() flow.Work) Handler
	SetFileExport(exporter *export.FileExporter) Handler
	ItemsToOperation(func(items []string) (operation Operation, err *data.CodeError)) Handler
	OnResult(func(operationInfo string, operation Operation, result *OperationResult)) Handler
	OnError(func(err *data.CodeError)) Handler
	Start()
}

func NewHandler(info Info) Handler {
	h := &handler{
		info: &info,
	}
	h.exporter = export.EmptyFileExport()
	return h
}

type handler struct {
	info                  *Info
	emptyOperation        func() flow.Work
	exporter              *export.FileExporter
	operationItemsCreator func(items []string) (operation Operation, err *data.CodeError)
	onError               func(err *data.CodeError)
	onResult              func(operationInfo string, operation Operation, result *OperationResult)
}

func (h *handler) EmptyOperation(emptyOperation func() flow.Work) Handler {
	h.emptyOperation = emptyOperation
	return h
}

func (h *handler) SetFileExport(exporter *export.FileExporter) Handler {
	h.exporter = exporter
	return h
}

func (h *handler) ItemsToOperation(reader func(items []string) (operation Operation, err *data.CodeError)) Handler {
	h.operationItemsCreator = reader
	return h
}

func (h *handler) OnResult(handler func(operationInfo string, operation Operation, result *OperationResult)) Handler {
	h.onResult = handler
	return h
}

func (h *handler) OnError(handler func(err *data.CodeError)) Handler {
	h.onError = handler
	return h
}

func (h *handler) Start() {
	isArraySource := h.info.WorkList != nil && len(h.info.WorkList) > 0
	if !isArraySource {
		if e := locker.TryLock(); e != nil {
			log.ErrorF("batch job, %v", e)
			return
		}

		unlockHandler := func() {
			if e := locker.TryUnlock(); e != nil {
				log.ErrorF("batch job, %v", e)
			}
		}
		workspace.AddCancelObserver(func(s os.Signal) {
			unlockHandler()
		})
		defer unlockHandler()
	}

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		h.onError(err)
		return
	}

	workBuilder := flow.New(h.info.Info)
	var workerBuilder *flow.WorkerProvideBuilder
	if isArraySource {
		workerBuilder = workBuilder.WorkProviderWithArray(h.info.WorkList)
	} else {
		log.DebugF("forceFlag: %v, overwriteFlag: %v, worker: %v, inputFile: %q, successFilePath: %q, failureFilePath: %q, sep: %q",
			h.info.Force, h.info.Overwrite, h.info.WorkerCount, h.info.InputFile, h.info.SuccessExportFilePath, h.info.FailExportFilePath, h.info.ItemSeparate)

		if h.operationItemsCreator == nil {
			log.Error(data.NewEmptyError().AppendDesc(alert.CannotEmpty("operation reader", "")))
			return
		}

		workerBuilder = workBuilder.WorkProviderWithFile(h.info.InputFile,
			h.info.EnableStdin,
			flow.NewItemsWorkCreator(h.info.ItemSeparate, h.info.MinItemsCount, func(items []string) (work flow.Work, err *data.CodeError) {
				return h.operationItemsCreator(items)
			}))
	}

	// overseer， EnableRecord 未开启不记录中间状态（数组类型的数据源默认关闭）
	dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
	if !isArraySource {
		if h.info.EnableRecord {
			log.DebugF("batch recorder:%s", dbPath)
		} else {
			log.Debug("batch recorder:Not Enable")
		}
	}

	metric := &Metric{}
	if isArraySource {
		metric.DisablePrintProgress()
	}
	metric.Start()
	workerBuilder.
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewWorker(func(workInfoList []*flow.WorkInfo) ([]*flow.WorkRecord, *data.CodeError) {
				recordList := make([]*flow.WorkRecord, 0, len(workInfoList))
				operationBucket := ""
				operationStringList := make([]string, 0, len(workInfoList))
				operationWorkInfoList := make([]*flow.WorkInfo, 0, len(workInfoList))
				for _, workInfo := range workInfoList {
					if operation, ok := workInfo.Work.(Operation); !ok {
						return nil, alert.Error("batch WorkerProvider, operation type conv error", "")
					} else {
						if len(operationBucket) == 0 {
							operationBucket = operation.GetBucket()
						}

						if operationString, e := operation.ToOperation(); e != nil {
							recordList = append(recordList, &flow.WorkRecord{
								WorkInfo: workInfo,
								Result:   nil,
								Err:      e,
							})
						} else {
							operationStringList = append(operationStringList, operationString)
							operationWorkInfoList = append(operationWorkInfoList, workInfo)
						}
					}
				}

				if cErr := bucket.CompleteBucketManagerRegion(bucketManager, operationBucket); cErr != nil {
					return nil, cErr
				}

				resultList, e := bucketManager.Batch(operationStringList)
				if len(resultList) != len(operationStringList) {
					return recordList, data.ConvertError(e)
				}

				for i, r := range resultList {
					result := &OperationResult{
						Code:     r.Code,
						Hash:     r.Data.Hash,
						FSize:    r.Data.Fsize,
						PutTime:  r.Data.PutTime,
						MimeType: r.Data.MimeType,
						Type:     r.Data.Type,
						Error:    r.Data.Error,
					}
					record := &flow.WorkRecord{
						WorkInfo: operationWorkInfoList[i],
						Result:   result,
					}
					if !result.IsSuccess() {
						record.Err = data.NewError(result.Code, result.Error)
					}
					recordList = append(recordList, record)
				}
				return recordList, nil
			}), nil
		})).
		DoWorkListMaxCount(h.info.OperationCountPerRequest).
		SetOverseerEnable(h.info.EnableRecord).
		SetDBOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: h.emptyOperation(),
				},
				Result: &OperationResult{},
				Err:    nil,
			}
		}).
		SetLimit(flow.NewBlockLimit(h.info.WorkerCount*h.info.OperationCountPerRequest,
			flow.MaxLimitCount(h.info.WorkerCount*h.info.OperationCountPerRequest),
			flow.MinLimitCount(h.info.MinWorkerCount*h.info.OperationCountPerRequest),
			flow.IncreaseLimitCount(h.info.OperationCountPerRequest),
			flow.IncreaseLimitCountPeriod(time.Duration(h.info.WorkerCountIncreasePeriod)*time.Second))).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			result, _ := workRecord.Result.(*OperationResult)

			if workRecord.Err == nil && result != nil && result.IsValid() {
				return false, nil
			}

			if !h.info.RecordRedoWhileError {
				if result == nil {
					return false, data.NewEmptyError().AppendDescF("result:nil error:%s", workRecord.Err)
				} else {
					return false, data.NewEmptyError().AppendDescF("result:%s error:%s", result.ErrorDescription(), workRecord.Err)
				}
			}

			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDescF("result is invalid, %s", result.ErrorDescription())
			}
			return false, nil
		}).
		OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			operationResult, _ := result.(*OperationResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					log.InfoF("Skip line:%s because have done and success", work.Data)
					h.exporter.Success().Export(work.Data)
				} else {
					metric.AddFailureCount(1)
					errDesc := ""
					if operationResult != nil {
						errDesc = operationResult.ErrorDescription()
					}
					log.InfoF("Skip line:%s because have done and failure, %v%s", work.Data, err, errDesc)
					h.exporter.Fail().ExportF("%s%s-%s", work.Data, flow.ErrorSeparate, errDesc)
				}
			} else {
				metric.AddSkippedCount(1)

				operation, _ := work.Work.(Operation)
				h.onResult(work.Data, operation, &OperationResult{
					Code:  data.ErrorCodeUnknown,
					Error: fmt.Sprintf("%v", err),
				})
				log.InfoF("Skip line:%s because:%v", work.Data, err)
				if err != nil && err.Code == data.ErrorCodeLineHeader {
					h.exporter.Fail().Export(work.Data)
				} else {
					h.exporter.Fail().ExportF("%s%s-%v", work.Data, flow.ErrorSeparate, err)
				}
			}
		}).
		OnWorkSuccess(func(work *flow.WorkInfo, result flow.Result) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			operation, _ := work.Work.(Operation)
			operationResult, _ := result.(*OperationResult)
			if operationResult != nil && operationResult.IsSuccess() {
				metric.AddSuccessCount(1)
				h.exporter.Success().Export(work.Data)
			} else {
				metric.AddFailureCount(1)
				if operationResult == nil {
					h.exporter.Fail().ExportF("%s%s-no result", work.Data, flow.ErrorSeparate)
				} else {
					h.exporter.Fail().ExportF("%s%s[%d]%s", work.Data, flow.ErrorSeparate, operationResult.Code, operationResult.Error)
				}
			}
			h.onResult(work.Data, operation, operationResult)
		}).
		OnWorkFail(func(work *flow.WorkInfo, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.AddFailureCount(1)
			metric.PrintProgress("Batching:" + work.Data)
			h.exporter.Fail().ExportF("%s%s[%d]%s", work.Data, flow.ErrorSeparate, err.Code, err.Desc)

			operation, _ := work.Work.(Operation)
			h.onResult(work.Data, operation, &OperationResult{
				Code:  err.Code,
				Error: err.Desc,
			})
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	if !isArraySource {
		log.InfoF("job dir:%s, there is a cache related to this command in this folder, which will also be used next time the same command is executed. If you are sure that you don’t need it, you can delete this folder.", workspace.GetJobDir())

		// 数组源不输出结果
		resultPath := filepath.Join(workspace.GetJobDir(), ".result")
		if e := utils.MarshalToFile(resultPath, metric); e != nil {
			log.ErrorF("save batch result to path:%s error:%v", resultPath, e)
		} else {
			log.DebugF("save batch result to path:%s", resultPath)
		}

		log.Alert("--------------- Batch Result ---------------")
		log.AlertF("%20s%10d", "Total:", metric.TotalCount)
		log.AlertF("%20s%10d", "Success:", metric.SuccessCount)
		log.AlertF("%20s%10d", "Failure:", metric.FailureCount)
		log.AlertF("%20s%10d", "Skipped:", metric.SkippedCount)
		log.AlertF("%20s%10ds", "Duration:", metric.Duration)
		log.AlertF("--------------------------------------------")
	}
}
