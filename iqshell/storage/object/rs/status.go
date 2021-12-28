package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

type StatusApiInfo struct {
	Bucket string
	Key    string
}

var _ BatchOperation = (*StatusApiInfo)(nil)

func (s StatusApiInfo) ToOperation() (string, error) {
	if len(s.Bucket) == 0 || len(s.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("status operation bucket or key", ""))
	}
	return storage.URIStat(s.Bucket, s.Key), nil
}
