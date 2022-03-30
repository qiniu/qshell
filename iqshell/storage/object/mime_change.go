package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type ChangeMimeApiInfo struct {
	Bucket string
	Key    string
	Mime   string
}

func (c ChangeMimeApiInfo) ToOperation() (string, error) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("change mime operation bucket or key", ""))
	}

	if len(c.Mime) == 0 {
		return "", errors.New(alert.CannotEmpty("change mime operation mime", ""))
	}

	return storage.URIChangeMime(c.Bucket, c.Key, c.Mime), nil
}

func (c ChangeMimeApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeMime|%s|%s|%s", c.Bucket, c.Key, c.Mime)
}

func ChangeMimeType(info ChangeMimeApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}
