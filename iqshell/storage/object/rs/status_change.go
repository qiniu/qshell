package rs

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
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