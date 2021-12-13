package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type DeleteApiInfo struct {
	Bucket    string
	Key       string
	AfterDays int
}

func (d DeleteApiInfo) ToOperation() (string, error) {
	if len(d.Bucket) == 0 || len(d.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("delete operation bucket or key", ""))
	}

	if d.AfterDays > 0 {
		return storage.URIDeleteAfterDays(d.Bucket, d.Key, d.AfterDays), nil
	} else {
		return storage.URIDelete(d.Bucket, d.Key), nil
	}
}

type DeleteApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newDeleteApiResult(ret storage.BatchOpRet) DeleteApiResult {
	return DeleteApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func Delete(info DeleteApiInfo) (DeleteApiResult, error) {
	ret, err := batch.One(info)
	if err != nil {
		return DeleteApiResult{}, err
	}
	return newDeleteApiResult(ret), err
}

func BatchDelete(apiInfoChan <-chan DeleteApiInfo) (<-chan DeleteApiResult, <-chan error) {
	if len(apiInfoChan) == 0 {
		return nil, nil
	}

	batchInfoChan := make(chan batch.BatchOperation)
	go func() {
		for apiInfo := range apiInfoChan {
			batchInfoChan <- apiInfo
		}
	}()

	batchResultChan, errChan := batch.BatchWithChannel(batchInfoChan)

	apiResultChan := make(chan DeleteApiResult)
	go func() {
		for item := range batchResultChan {
			apiResultChan <- newDeleteApiResult(item)
		}
	}()

	return apiResultChan, errChan
}
