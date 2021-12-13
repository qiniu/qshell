package rs

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type ChangeStatusApiInfo struct {
	Bucket string
	Key    string
	Status int
}

func (c ChangeStatusApiInfo) ToOperation() (string, error) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("change status operation bucket or key", ""))
	}
	return fmt.Sprintf("/chstatus/%s/status/%c", storage.EncodedEntry(c.Bucket, c.Key), c.Status), nil
}

type ChangeStatusApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func newChangeStatusApiResult(ret storage.BatchOpRet) ChangeStatusApiResult {
	return ChangeStatusApiResult{
		Code:  ret.Code,
		Error: ret.Data.Error,
	}
}

func ChangeStatus(info ChangeStatusApiInfo) (ChangeStatusApiResult, error) {
	ret, err := batch.One(info)
	if err != nil {
		return ChangeStatusApiResult{}, err
	}
	return newChangeStatusApiResult(ret), err
}

func BatchChangeStatus(infos []ChangeStatusApiInfo) ([]ChangeStatusApiResult, error) {
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

	ret := make([]ChangeStatusApiResult, len(result))
	for _, item := range result {
		ret = append(ret, newChangeStatusApiResult(item))
	}

	return ret, err
}
