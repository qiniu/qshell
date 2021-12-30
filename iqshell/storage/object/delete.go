package object

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type DeleteApiInfo struct {
	Bucket    string
	Key       string
	AfterDays int
	Condition rs.OperationCondition
}

func (d DeleteApiInfo) ToOperation() (string, error) {
	if len(d.Bucket) == 0 || len(d.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("delete operation bucket or key", ""))
	}

	condition := rs.OperationConditionURI(d.Condition)
	if d.AfterDays > 0 {
		return storage.URIDeleteAfterDays(d.Bucket, d.Key, d.AfterDays) + condition, nil
	} else {
		return storage.URIDelete(d.Bucket, d.Key) + condition, nil
	}
}

func Delete(info DeleteApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}