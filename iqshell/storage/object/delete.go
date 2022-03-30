package object

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type DeleteApiInfo struct {
	Bucket          string
	Key             string
	DeleteAfterDays int
	Condition       batch.OperationCondition
}

func (d DeleteApiInfo) ToOperation() (string, error) {
	if len(d.Bucket) == 0 || len(d.Key) == 0 {
		return "", alert.CannotEmptyError("delete operation bucket or key", "")
	}

	condition := batch.OperationConditionURI(d.Condition)
	if d.DeleteAfterDays < 0 {
		return "", alert.Error("DeleteAfterDays can't be smaller than 0", "")
	} else if d.DeleteAfterDays == 0 {
		return storage.URIDelete(d.Bucket, d.Key) + condition, nil
	} else {
		return storage.URIDeleteAfterDays(d.Bucket, d.Key, d.DeleteAfterDays) + condition, nil
	}
}

func Delete(info DeleteApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}
