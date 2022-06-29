package object

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type DeleteApiInfo struct {
	Bucket          string                   `json:"bucket"`
	Key             string                   `json:"key"`
	DeleteAfterDays int                      `json:"delete_after_days"`
	Condition       batch.OperationCondition `json:"condition"`
}

func (d *DeleteApiInfo) ToOperation() (string, *data.CodeError) {
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

func (d *DeleteApiInfo) WorkId() string {
	return fmt.Sprintf("Delete|%s|%s", d.Bucket, d.Key)
}

func Delete(info *DeleteApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
