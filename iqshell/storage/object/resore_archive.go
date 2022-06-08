package object

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type RestoreArchiveApiInfo struct {
	Bucket          string `json:"bucket"`
	Key             string `json:"key"`
	FreezeAfterDays int    `json:"freeze_after_days"`
}

func (r *RestoreArchiveApiInfo) ToOperation() (string, *data.CodeError) {
	if len(r.Bucket) == 0 || len(r.Key) == 0 {
		return "", alert.CannotEmptyError("Restore archive operation bucket or key", "")
	}

	return storage.URIRestoreAr(r.Bucket, r.Key, r.FreezeAfterDays), nil
}

func (r *RestoreArchiveApiInfo) WorkId() string {
	return fmt.Sprintf("RestoreArchive|%s|%s|%d", r.Bucket, r.Key, r.FreezeAfterDays)
}

func RestoreArchive(info *RestoreArchiveApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
