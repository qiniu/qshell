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
