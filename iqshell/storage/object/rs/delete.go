package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
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
		Code:     ret.Code,
		Error:    ret.Data.Error,
	}
}

func Delete(info DeleteApiInfo) (DeleteApiResult, error) {
	ret, err := BatchOne(info)
	if err != nil {
		return DeleteApiResult{}, err
	}
	return newDeleteApiResult(ret), err
}

func BatchDelete(infos []DeleteApiInfo) ([]DeleteApiResult, error) {
	if len(infos) == 0 {
		return nil, nil
	}

	operations := make([]BatchOperation, len(infos))
	for _, info := range infos {
		operations = append(operations, info)
	}

	result, err := Batch(operations)
	if result == nil || len(result) == 0 {
		return nil, err
	}

	ret := make([]DeleteApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newDeleteApiResult(item))
	}

	return ret, err
}
