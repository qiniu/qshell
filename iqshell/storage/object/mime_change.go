package object

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type ChangeMimeApiInfo struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Mime   string `json:"mime"`
}

func (c *ChangeMimeApiInfo) GetBucket() string {
	return c.Bucket
}

func (c *ChangeMimeApiInfo) ToOperation() (string, *data.CodeError) {
	if len(c.Bucket) == 0 || len(c.Key) == 0 {
		return "", alert.CannotEmptyError("change mime operation bucket or key", "")
	}

	if len(c.Mime) == 0 {
		return "", alert.CannotEmptyError("change mime operation mime", "")
	}

	return storage.URIChangeMime(c.Bucket, c.Key, c.Mime), nil
}

func (c *ChangeMimeApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeMime|%s|%s|%s", c.Bucket, c.Key, c.Mime)
}

func ChangeMimeType(info *ChangeMimeApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
