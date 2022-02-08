package batch

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type Info struct {
	group.Info

	MaxOperationCountPerRequest int
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

	if f.info.MaxOperationCountPerRequest <= 0 ||
		f.info.MaxOperationCountPerRequest > defaultOperationCountPerRequest {
		f.info.MaxOperationCountPerRequest = defaultOperationCountPerRequest
	}

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		f.onError(err)
		return
	}

	work.NewFlowHandler(f.info.Info.Info).ReadWork(func() (work work.Work, hasMore bool) {
		task := &batchOperations{
			operations:       make([]Operation, 0, 0),
			operationStrings: make([]string, 0, 0),
		}

		for {
			if workspace.IsCmdInterrupt() {
				break
			}
			operation, hasMore := f.readOperation()
			if operation == nil {
				if hasMore {
					log.Debug("batch task producer: operation invalid")
					continue
				} else {
					log.Debug("batch task producer: read operation complete")
					break
				}
			}
			operationString, err := operation.ToOperation()
			if err != nil {
				log.Debug("batch task producer: parse operation error")
				log.Warning(err)
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
		if err != nil {
			log.DebugF("batch task consumer: batch error:%v", err)
			f.onError(err)
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		task := work.(*batchOperations)
		results := result.([]storage.BatchOpRet)
		log.Debug("batch task consumer: batch success")
		for i, r := range results {
			f.onResult(
				task.operations[i],
				OperationResult{
					Code:     r.Code,
					Hash:     r.Data.Hash,
					FSize:    r.Data.Fsize,
					PutTime:  r.Data.PutTime,
					MimeType: r.Data.MimeType,
					Type:     r.Data.Type,
					Error:    r.Data.Error,
				})
		}
	}).OnWorksComplete(func() {
		log.Debug("batch: end")
	}).Start()
}

type batchOperations struct {
	operations       []Operation
	operationStrings []string
}
