package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

type CopyApiInfo struct {
	SourceBucket string
	SourceKey    string
	DestBucket   string
	DestKey      string
	Force        bool
}

func (m CopyApiInfo) ToOperation() (string, error) {
	if len(m.SourceBucket) == 0 || len(m.SourceKey) == 0 {
		return "", errors.New(alert.CannotEmpty("copy operation bucket or key of source and dest", ""))
	}

	return storage.URICopy(m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey, m.Force), nil
}

type CopyApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newCopyApiResult(ret storage.BatchOpRet) CopyApiResult {
	return CopyApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func Copy(info CopyApiInfo) (CopyApiResult, error) {
	ret, err := BatchOne(info)
	if err != nil {
		return CopyApiResult{}, err
	}
	return newCopyApiResult(ret), err
}

func BatchCopy(infos []CopyApiInfo) ([]CopyApiResult, error) {
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

	ret := make([]CopyApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newCopyApiResult(item))
	}

	return ret, err
}