package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type ChangeMimeApiInfo struct {
	Bucket string
	Key    string
	Mime   string
}

func (c ChangeMimeApiInfo) ToOperation() (string, error) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("change mime operation bucket or key", ""))
	}

	if len(c.Mime) == 0 {
		return "", errors.New(alert.CannotEmpty("change mime operation mime", ""))
	}

	return storage.URIChangeMime(c.Bucket, c.Key, c.Mime), nil
}

type ChangeMimeApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newChangeMimeApiResult(ret storage.BatchOpRet) ChangeMimeApiResult {
	return ChangeMimeApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func ChangeMime(info ChangeMimeApiInfo) (ChangeMimeApiResult, error) {
	ret, err := batch.One(info)
	if err != nil {
		return ChangeMimeApiResult{}, err
	}
	return newChangeMimeApiResult(ret), err
}

func BatchChangeMime(infos []ChangeMimeApiInfo) ([]ChangeMimeApiResult, error) {
	if len(infos) == 0 {
		return nil, nil
	}

	operations := make([]batch.BatchOperation, len(infos))
	for _, info := range infos {
		operations = append(operations, info)
	}

	result, err := batch.Batch(operations)
	if result == nil || len(result) == 0 {
		return nil, err
	}

	ret := make([]ChangeMimeApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newChangeMimeApiResult(item))
	}

	return ret, err
}
