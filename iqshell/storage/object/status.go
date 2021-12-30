package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type StatusApiInfo struct {
	Bucket string
	Key    string
}

func (s StatusApiInfo) ToOperation() (string, error) {
	if len(s.Bucket) == 0 || len(s.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("status operation bucket or key", ""))
	}
	return storage.URIStat(s.Bucket, s.Key), nil
}

func Status(info StatusApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}

// ChangeStatusApiInfo 修改 status
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

func ChangeStatus(info ChangeStatusApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}