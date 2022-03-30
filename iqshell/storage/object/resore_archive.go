package object

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type RestoreArchiveApiInfo struct {
	Bucket          string
	Key             string
	FreezeAfterDays int
}

func (r RestoreArchiveApiInfo) ToOperation() (string, error) {
	if len(r.Bucket) == 0 || len(r.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("Restore archive operation bucket or key", ""))
	}

	return storage.URIRestoreAr(r.Bucket, r.Key, r.FreezeAfterDays), nil
}

func (r RestoreArchiveApiInfo) WorkId() string {
	return fmt.Sprintf("RestoreArchive|%s|%s|%d", r.Bucket, r.Key, r.FreezeAfterDays)
}

func RestoreArchive(info RestoreArchiveApiInfo) (batch.OperationResult, error) {
	return batch.One(info)
}
