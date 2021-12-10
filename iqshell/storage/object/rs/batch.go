package rs

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type BatchOperation interface {
	ToOperation() (string, error)
}

func BatchOne(operation BatchOperation) (storage.BatchOpRet, error) {
	ret, err := Batch([]BatchOperation{operation})
	if err != nil || len(ret) == 0 {
		return storage.BatchOpRet{}, err
	}
	return ret[0], err
}

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
