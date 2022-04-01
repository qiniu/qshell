package storage

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type PrefetchApiInfo struct {
	Bucket string
	Key    string
}

func Prefetch(info PrefetchApiInfo) *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("bucket", "")
	}

	if len(info.Key) == 0 {
		return alert.CannotEmptyError("key", "")
	}

	m, err := bucket.GetBucketManager()
	if err != nil {
		return err
	}

	return data.ConvertError(m.Prefetch(info.Bucket, info.Key))
}
