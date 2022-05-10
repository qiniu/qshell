package object

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type CopyApiInfo struct {
	SourceBucket string
	SourceKey    string
	DestBucket   string
	DestKey      string
	Force        bool
}

func (m *CopyApiInfo) ToOperation() (string, *data.CodeError) {
	if len(m.SourceBucket) == 0 || len(m.SourceKey) == 0 || len(m.DestBucket) == 0 || len(m.DestKey) == 0 {
		return "", alert.CannotEmptyError("copy operation bucket or key of source and dest", "")
	}

	return storage.URICopy(m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey, m.Force), nil
}

func (m *CopyApiInfo) WorkId() string {
	return fmt.Sprintf("Copy|%s|%s|%s|%s", m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey)
}

func Copy(info *CopyApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
