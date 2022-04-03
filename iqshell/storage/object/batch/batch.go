package batch

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
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

	MaxOperationCountPerRequest int
}

func (info *Info) Check() *data.CodeError {
	if err := info.Info.Check(); err != nil {
		return err
	}
	info.Force = true

	if info.MinItemsCount < 1 {
		info.MinItemsCount = 1
	}

	if info.MaxOperationCountPerRequest <= 0 ||
		info.MaxOperationCountPerRequest > defaultOperationCountPerRequest {
		info.MaxOperationCountPerRequest = defaultOperationCountPerRequest
	}

	return nil
}

type Handler interface {
	ItemsToOperation(func(items []string) (operation Operation, err *data.CodeError)) Handler
	OnResult(func(operationInfo string, operation Operation, result *OperationResult)) Handler
	OnError(func(err *data.CodeError)) Handler
	Start()
}

func NewHandler(info Info) Handler {
	return &handler{
		info: &info,
	}
}

type handler struct {
	info                  *Info
	operationItemsCreator func(items []string) (operation Operation, err *data.CodeError)
	onError               func(err *data.CodeError)
	onResult              func(operationInfo string, operation Operation, result *OperationResult)
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
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		h.onError(err)
		return
	}

	workBuilder := flow.New(h.info.Info)
	var workerBuilder *flow.WorkerProvideBuilder
	if h.info.WorkList != nil && len(h.info.WorkList) > 0 {
		workerBuilder = workBuilder.WorkProviderWithArray(h.info.WorkList)
	} else {
		log.DebugF("forceFlag: %v, overwriteFlag: %v, worker: %v, inputFile: %q, bsuccessFname: %q, bfailureFname: %q, sep: %q",
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

	workerBuilder.
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewWorker(func(workInfoList []*flow.WorkInfo) ([]*flow.WorkRecord, *data.CodeError) {

				recordList := make([]*flow.WorkRecord, 0, len(workInfoList))
				operationStrings := make([]string, 0, len(workInfoList))
				for _, workInfo := range workInfoList {
					recordList = append(recordList, &flow.WorkRecord{
						WorkInfo: workInfo,
					})

					if operation, ok := workInfo.Work.(Operation); !ok {
						return nil, alert.Error("batch WorkerProvider, operation type conv error", "")
					} else {
						if operationString, e := operation.ToOperation(); e != nil {
							return nil, alert.Error("batch WorkerProvider, ToOperation error:"+e.Error(), "")
						} else {
							operationStrings = append(operationStrings, operationString)
						}
					}
				}

				resultList, e := bucketManager.Batch(operationStrings)
				if len(resultList) != len(operationStrings) {
					return recordList, data.ConvertError(e)
				}

				for i, r := range resultList {
					operationResult := &OperationResult{
						Code:     r.Code,
						Hash:     r.Data.Hash,
						FSize:    r.Data.Fsize,
						PutTime:  r.Data.PutTime,
						MimeType: r.Data.MimeType,
						Type:     r.Data.Type,
						Error:    r.Data.Error,
					}
					recordList[i].Result = operationResult
				}
				return recordList, nil
			}), nil
		})).
		DoWorkListMaxCount(h.info.MaxOperationCountPerRequest).
		OnWorkSkip(func(work *flow.WorkInfo, err *data.CodeError) {

		}).
		OnWorkSuccess(func(work *flow.WorkInfo, result flow.Result) {
			operation, _ := work.Work.(Operation)
			operationResult, _ := result.(*OperationResult)
			h.onResult(work.Data, operation, operationResult)
		}).
		OnWorkFail(func(work *flow.WorkInfo, err *data.CodeError) {
			operation, _ := work.Work.(Operation)
			h.onResult(work.Data, operation, &OperationResult{
				Code:  data.ErrorCodeUnknown,
				Error: err.Error(),
			})
		}).Build().Start()
}
