package rs

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"sync"
)

const (
	MaxOperationCountPerRequest = 1000
)

type OperationCondition struct {
	PutTime string
}

func OperationConditionURI(condition OperationCondition) string {
	cond := ""
	if condition.PutTime != "" {
		cond += "putTime=" + condition.PutTime
	}
	if cond == "" {
		return ""
	}
	return fmt.Sprintf("/cond/%s", base64.URLEncoding.EncodeToString([]byte(cond)))
}

type BatchOperation interface {
	ToOperation() (string, error)
}

type OperationResult struct {
	Code     int
	Hash     string
	FSize    int64
	PutTime  int64
	MimeType string
	Type     int
	Error    string
}

type BatchHandler interface {
	WorkCount() int
	ReadOperation() (operation BatchOperation, complete bool)
	HandlerResult(operation BatchOperation, result OperationResult)
	HandlerError(err error)
}

func BatchOne(operation BatchOperation) (OperationResult, error) {
	ret, err := Batch([]BatchOperation{operation})
	if err != nil || len(ret) == 0 {
		return OperationResult{}, err
	}

	return ret[0], nil
}

// Batch operations 长度最大有限制
func Batch(operations []BatchOperation) ([]OperationResult, error) {
	bm, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	ret := make([]OperationResult, 0, len(operations))
	operationStrings := make([]string, MaxOperationCountPerRequest)
	for i, operation := range operations {
		if workspace.IsCmdInterrupt() {
			break
		}

		if operationString, err := operation.ToOperation(); err == nil {
			operationStrings = append(operationStrings, operationString)
		} else {
			log.Warning(err)
		}

		if i == len(operations) || len(operationStrings) >= MaxOperationCountPerRequest {
			operationStrings = make([]string, MaxOperationCountPerRequest)
			results, bEerr := bm.Batch(operationStrings)
			for _, result := range results {
				ret = append(ret, OperationResult{
					Code:     result.Code,
					Hash:     result.Data.Hash,
					FSize:    result.Data.Fsize,
					PutTime:  result.Data.PutTime,
					MimeType: result.Data.MimeType,
					Type:     result.Data.Type,
					Error:    result.Data.Error,
				})
			}

			if bEerr != nil {
				err = bEerr
				break
			}
		}
	}

	return ret, err
}

func BatchWithHandler(handler BatchHandler) {
	if handler == nil {
		return
	}

	bm, err := bucket.GetBucketManager()
	if err != nil {
		handler.HandlerError(err)
		return
	}

	workCount := handler.WorkCount()
	if workCount == 0 {
		workCount = 1
	}

	log.Debug("batch: start")
	log.DebugF("work count: %d", workCount)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(workCount + 1)
	batchOperationChan := make(chan batchTask)
	go func() {
		batchTaskProduct(bm, handler, batchOperationChan)
		waitGroup.Done()
	}()

	for i := 0; i < workCount; i++ {
		go func() {
			batchTaskConsume(handler, batchOperationChan)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	log.Debug("batch: end")
}

type batchTask struct {
	manager          *storage.BucketManager
	operations       []BatchOperation
	operationStrings []string
}

func batchTaskProduct(bucketManager *storage.BucketManager, handler BatchHandler, batchTaskChan chan<- batchTask) {
	task := batchTask{
		manager:          bucketManager,
		operations:       make([]BatchOperation, 0, MaxOperationCountPerRequest),
		operationStrings: make([]string, 0, MaxOperationCountPerRequest),
	}

	log.Debug("batch task producer: start")
	for {
		if workspace.IsCmdInterrupt() {
			break
		}

		operation, complete := handler.ReadOperation()
		if complete {
			log.Debug("batch task producer: read operation complete")
			batchTaskChan <- task
			break
		}

		if operation == nil {
			log.Debug("batch task producer: operation invalid")
			continue
		}

		operationString, err := operation.ToOperation()
		if err != nil {
			log.Debug("batch task producer: parse operation error")
			log.Warning(err)
			continue
		}

		task.operations = append(task.operations, operation)
		task.operationStrings = append(task.operationStrings, operationString)

		if len(task.operationStrings) >= MaxOperationCountPerRequest {
			batchTaskChan <- task
			task = batchTask{
				operations:       make([]BatchOperation, 0, MaxOperationCountPerRequest),
				operationStrings: make([]string, 0, MaxOperationCountPerRequest),
			}
			log.Debug("batch task producer: produce one task: task count:%d", len(task.operations))
		}
	}
	close(batchTaskChan)
	log.Debug("batch task producer: end")
}

func batchTaskConsume(handler BatchHandler, batchTaskChan <-chan batchTask) {
	log.Debug("batch task consumer: start")
	for task := range batchTaskChan {
		if workspace.IsCmdInterrupt() {
			break
		}

		log.Debug("batch task consumer: get new task")
		if len(task.operations) == 0 {
			log.Debug("batch task consumer: task invalid")
			continue
		}

		results, err := task.manager.Batch(task.operationStrings)
		if err != nil {
			log.Debug("batch task consumer: batch error")
			handler.HandlerError(err)
		}

		log.Debug("batch task consumer: batch success")
		for i, result := range results {
			handler.HandlerResult(
				task.operations[i],
				OperationResult{
					Code:     result.Code,
					Hash:     result.Data.Hash,
					FSize:    result.Data.Fsize,
					PutTime:  result.Data.PutTime,
					MimeType: result.Data.MimeType,
					Type:     result.Data.Type,
					Error:    result.Data.Error,
				})
		}
	}
	log.Debug("batch task consumer: end")
}
