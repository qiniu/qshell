package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

type CopyApiInfo struct {
	SourceBucket string
	SourceKey    string
	DestBucket   string
	DestKey      string
	Force        bool
}

func (m CopyApiInfo) ToOperation() (string, error) {
	if len(m.SourceBucket) == 0 || len(m.SourceKey) == 0 || len(m.DestBucket) == 0 || len(m.DestKey) == 0 {
		return "", errors.New(alert.CannotEmpty("copy operation bucket or key of source and dest", ""))
	}

	return storage.URICopy(m.SourceBucket, m.SourceKey, m.DestBucket, m.DestKey, m.Force), nil
}