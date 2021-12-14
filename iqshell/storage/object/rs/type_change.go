package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

type ChangeTypeApiInfo struct {
	Bucket string
	Key    string
	Type   int
}

func (c ChangeTypeApiInfo) ToOperation() (string, error) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("change type operation bucket or key", ""))
	}

	return storage.URIChangeType(c.Bucket, c.Key, c.Type), nil
}

type ChangeTypeApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newChangeTypeApiResult(ret storage.BatchOpRet) ChangeTypeApiResult {
	return ChangeTypeApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func ChangeType(info ChangeTypeApiInfo) (ChangeTypeApiResult, error) {
	ret, err := BatchOne(info)
	if err != nil {
		return ChangeTypeApiResult{}, err
	}
	return newChangeTypeApiResult(ret), err
}

func BatchChangeType(infos []ChangeTypeApiInfo) ([]ChangeTypeApiResult, error) {
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

	ret := make([]ChangeTypeApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newChangeTypeApiResult(item))
	}

	return ret, err
}
