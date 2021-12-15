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
	Condition OperationCondition
}

func (d DeleteApiInfo) ToOperation() (string, error) {
	if len(d.Bucket) == 0 || len(d.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("delete operation bucket or key", ""))
	}

	condition := OperationConditionURI(d.Condition)
	if d.AfterDays > 0 {
		return storage.URIDeleteAfterDays(d.Bucket, d.Key, d.AfterDays) + condition, nil
	} else {
		return storage.URIDelete(d.Bucket, d.Key) + condition, nil
	}
}

type DeleteApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
