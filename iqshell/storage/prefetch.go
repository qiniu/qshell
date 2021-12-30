package storage

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type PrefetchApiInfo struct {
	Bucket string
	Key    string
}

func Prefetch(info PrefetchApiInfo) error {
	if len(info.Bucket) == 0 {
		return errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.Key) == 0 {
		return errors.New(alert.CannotEmpty("key", ""))
	}

	m, err := bucket.GetBucketManager()
	if err != nil {
		return err
	}

	return m.Prefetch(info.Bucket, info.Key)
}
