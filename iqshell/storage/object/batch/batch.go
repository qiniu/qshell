package batch

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

const (
	MaxOperationCountPerRequest = 1000
)

type BatchOperation interface {
	ToOperation() (string, error)
}

func One(operation BatchOperation) (storage.BatchOpRet, error) {
	ret, err := Batch([]BatchOperation{operation})
	if err != nil || len(ret) == 0 {
		return storage.BatchOpRet{}, err
	}
	return ret[0], err
}

// Batch operations 长度最大有限制
func Batch(operations []BatchOperation) ([]storage.BatchOpRet, error) {
	bm, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	operationStrings := make([]string, len(operations))
	for _, operation := range operations {
		operationString, err := operation.ToOperation()
		if err != nil {
			log.Warning(err)
			continue
		}
		operationStrings = append(operationStrings, operationString)
	}
	return bm.Batch(operationStrings)
}

func BatchWithChannel(operations <-chan BatchOperation) (<-chan storage.BatchOpRet, <-chan error) {

	errChan := make(chan error)
	resultChan := make(chan storage.BatchOpRet)

	go func() {
		operationArray := make([]BatchOperation, 0, MaxOperationCountPerRequest)
		for operation := range operations {
			if len(operationArray) < MaxOperationCountPerRequest {
				operationArray = append(operationArray, operation)
			} else {
				results, err := Batch(operationArray)
				for _, result := range results {
					resultChan <- result
				}
				if err != nil {
					errChan <- err
					break
				}
			}
		}

		// 关闭通道
		close(errChan)
		close(resultChan)
	}()

	return resultChan, errChan
}
