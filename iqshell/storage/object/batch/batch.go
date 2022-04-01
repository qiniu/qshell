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

	Force     bool // 无需验证即可 batch 操作，类似于验证码验证
	Overwrite bool // 是否覆盖

	// 工作数据源
	WorkList     []flow.Work // 工作数据源：列表
	InputFile    string      // 工作数据源：文件
	ItemSeparate string      // 工作数据源：每行元素按分隔符分的分隔符
	EnableStdin  bool        // 工作数据源：stdin, 当 InputFile 不存在时使用 stdin

	MaxOperationCountPerRequest int
}

func (info *Info) Check() *data.CodeError {
	if err := info.Info.Check(); err != nil {
		return err
	}

	if info.MaxOperationCountPerRequest <= 0 ||
		info.MaxOperationCountPerRequest > defaultOperationCountPerRequest {
		info.MaxOperationCountPerRequest = defaultOperationCountPerRequest
	}

	return nil
}

type Handler interface {
	ItemsToOperation(func(items []string) (operation Operation, err *data.CodeError)) Handler
	OnResult(func(operation Operation, result *OperationResult)) Handler
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
	onResult              func(operation Operation, result *OperationResult)
}

func (h *handler) ItemsToOperation(reader func(items []string) (operation Operation, err *data.CodeError)) Handler {
	h.operationItemsCreator = reader
	return h
}

func (h *handler) OnResult(handler func(operation Operation, result *OperationResult)) Handler {
	h.onResult = handler
	return h
}

func (h *handler) OnError(handler func(err *data.CodeError)) Handler {
	h.onError = handler
	return h
}

func (h *handler) Start() {
	if h.operationItemsCreator == nil {
		log.Error(data.NewEmptyError().AppendDesc(alert.CannotEmpty("operation reader", "")))
		return
	}

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		h.onError(err)
		return
	}

	f := &flow.Flow{}

	// 配置 work provider
	if h.info.WorkList != nil && len(h.info.WorkList) > 0 {
		if provider, e := flow.NewArrayWorkProvider(h.info.WorkList); e != nil {
			return
		} else {
			f.WorkProvider = provider
		}
	} else {
		workCreator := flow.NewLineSeparateWorkCreator(h.info.ItemSeparate, func(items []string) (work flow.Work, err *data.CodeError) {
			return h.operationItemsCreator(items)
		})
		if provider, e := flow.NewWorkProviderOfFile(h.info.InputFile, h.info.EnableStdin, workCreator); e != nil {
			return
		} else {
			f.WorkProvider = provider
		}
	}

	// 配置 work 打包
	f.WorkPacker = flow.NewWorkPacker(h.info.MaxOperationCountPerRequest)

	// 配置 worker provider
	f.WorkerProvider = flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
		// worker creator
		return flow.NewWorker(func(work flow.Work) (flow.Result, *data.CodeError) {

			// batch 使用了 work package
			workPackage, ok := work.(*flow.WorkPackage)
			if !ok {
				return nil, alert.Error("batch WorkerProvider, WorkPackage type conv error", "")
			}

			operationStrings := make([]string, 0, len(workPackage.WorkRecords))
			for _, record := range workPackage.WorkRecords {
				if operation, ok := record.Work.(Operation); !ok {
					return nil, alert.Error("batch WorkerProvider, operation type conv error", "")
				} else {
					if operationString, err := operation.ToOperation(); err != nil {
						return nil, alert.Error("batch WorkerProvider, ToOperation error:"+err.Error(), "")
					} else {
						operationStrings = append(operationStrings, operationString)
					}
				}
			}

			if result, e := bucketManager.Batch(operationStrings); e != nil {
				for _, r := range workPackage.WorkRecords {
					r.Err = data.ConvertError(e)
				}
			} else {
				for i, r := range result {
					operationResult := &OperationResult{
						Code:     r.Code,
						Hash:     r.Data.Hash,
						FSize:    r.Data.Fsize,
						PutTime:  r.Data.PutTime,
						MimeType: r.Data.MimeType,
						Type:     r.Data.Type,
						Error:    r.Data.Error,
					}
					workPackage.WorkRecords[i].Result = operationResult
				}
			}
			return workPackage, nil
		}), nil
	})

	// 配置时间监听
	f.EventListener = flow.EventListener{
		WillWorkFunc:   nil,
		OnWorkSkipFunc: nil,
		OnWorkSuccessFunc: func(work flow.Work, result flow.Result) {
			operation, _ := work.(Operation)
			operationResult, _ := result.(*OperationResult)
			h.onResult(operation, operationResult)
		},
		OnWorkFailFunc: func(work flow.Work, err *data.CodeError) {
			operation, _ := work.(Operation)
			h.onResult(operation, &OperationResult{
				Code:  data.ErrorCodeUnknown,
				Error: err.Error(),
			})
		},
	}

	// 开始
	f.Start()
}
