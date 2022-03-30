package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
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

func (c ChangeTypeApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeStatus|%s|%s|%d", c.Bucket, c.Key, c.Type)
}

func ChangeType(info ChangeTypeApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}
