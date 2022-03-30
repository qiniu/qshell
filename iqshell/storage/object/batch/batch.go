package batch

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type Info struct {
	group.Info

	MaxOperationCountPerRequest int
	OperationOverseer           work.Overseer
}

func (info *Info) Check() error {
	if err := info.Info.Check(); err != nil {
		return err
	}

	if info.MaxOperationCountPerRequest <= 0 ||
		info.MaxOperationCountPerRequest > defaultOperationCountPerRequest {
		info.MaxOperationCountPerRequest = defaultOperationCountPerRequest
	}

	return nil
}

type Flow interface {
	ReadOperation(func() (operation Operation, complete bool)) Flow
	OnResult(func(operation Operation, result OperationResult)) Flow
	OnError(func(err error)) Flow
	Start()
}

func NewFlow(info Info) Flow {
	return &flow{
		info: &info,
	}
}

type flow struct {
	info          *Info
	readOperation func() (operation Operation, hasMore bool)
	onError       func(err error)
	onResult      func(operation Operation, result OperationResult)
}

func (f *flow) ReadOperation(reader func() (operation Operation, complete bool)) Flow {
	f.readOperation = reader
	return f
}

func (f *flow) OnResult(handler func(operation Operation, result OperationResult)) Flow {
	f.onResult = handler
	return f
}

func (f *flow) OnError(handler func(err error)) Flow {
	f.onError = handler
	return f
}

func (f *flow) Start() {
	if f.readOperation == nil {
		log.Error(errors.New(alert.CannotEmpty("operation reader", "")))
		return
	}

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		f.onError(err)
		return
	}

	work.NewFlowHandler(f.info.Info.FlowInfo).ReadWork(func() (work work.Work, hasMore bool) {
		task := &batchOperations{
			operations:       make([]Operation, 0, 0),
			operationStrings: make([]string, 0, 0),
		}

		for {
			operation, hasMore := f.readOperation()
			if operation == nil {
				if hasMore {
					log.Debug("batch task producer: operation invalid: value is nil")
					continue
				} else {
					log.Debug("batch task producer: read operation complete")
					break
				}
			}

			if f.info.OperationOverseer != nil {
				if f.info.OperationOverseer.HasDone(operation) {
					continue
				}
				f.info.OperationOverseer.WillWork(operation)
			}

			operationString, err := operation.ToOperation()
			if err != nil {
				f.onResult(operation, OperationResult{
					Code:  -999,
					Error: err.Error(),
				})
				continue
			}
			task.operations = append(task.operations, operation)
			task.operationStrings = append(task.operationStrings, operationString)
		}
		log.DebugF("batch task producer: produce one task: task count:%d", len(task.operations))
		if len(task.operations) == 0 {
			return nil, hasMore
		}

		return task, hasMore
	}).DoWork(func(work work.Work) (work.Result, error) {
		task := work.(*batchOperations)
		return bucketManager.Batch(task.operationStrings)
	}).OnWorkError(func(work work.Work, err error) {
		// 当出现此错误时，所有的均失败了
		task := work.(*batchOperations)
		for _, operation := range task.operations {
			operationResult := OperationResult{
				Code:  data.ErrorCodeUnknown,
				Error: err.Error(),
			}
			if f.info.OperationOverseer != nil {
				f.info.OperationOverseer.WorkDone(operation, operationResult, nil)
			}
			f.onResult(operation, operationResult)
		}

		if err != nil {
			log.DebugF("batch task consumer: batch error:%v", err)
			f.onError(err)
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		task := work.(*batchOperations)
		results := result.([]storage.BatchOpRet)
		log.Debug("batch task consumer: batch success")
		for i, r := range results {
			operationResult := OperationResult{
				Code:     r.Code,
				Hash:     r.Data.Hash,
				FSize:    r.Data.Fsize,
				PutTime:  r.Data.PutTime,
				MimeType: r.Data.MimeType,
				Type:     r.Data.Type,
				Error:    r.Data.Error,
			}
			if f.info.OperationOverseer != nil {
				f.info.OperationOverseer.WorkDone(task.operations[i], operationResult, nil)
			}
			f.onResult(task.operations[i], operationResult)
		}
	}).OnWorksComplete(func() {
		log.Debug("batch: end")
	}).Start()
}

type batchOperations struct {
	operations       []Operation
	operationStrings []string
}

func (b batchOperations) WorkId() string {
	return ""
}
