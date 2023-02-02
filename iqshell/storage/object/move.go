package object

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type MoveApiInfo struct {
	SourceBucket string `json:"source_bucket"`
	SourceKey    string `json:"source_key"`
	DestBucket   string `json:"dest_bucket"`
	DestKey      string `json:"dest_key"`
	Force        bool   `json:"force"`
}

func (m *MoveApiInfo) GetBucket() string {
	return m.SourceBucket
}

func (m *MoveApiInfo) ToOperation() (string, *data.CodeError) {
	if len(m.SourceBucket) == 0 || len(m.SourceKey) == 0 || len(m.DestBucket) == 0 || len(m.DestKey) == 0 {
		return "", alert.CannotEmptyError("move operation bucket or key of source and dest", "")
	}

	return storage.URIMove(m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey, m.Force), nil
}

func (m *MoveApiInfo) WorkId() string {
	return fmt.Sprintf("Move|%s|%s|%s|%s", m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey)
}

func Move(info *MoveApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
