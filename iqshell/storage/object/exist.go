package object

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type ExistApiInfo struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func Exist(info ExistApiInfo) (exists bool, err *data.CodeError) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return false, err
	}

	if len(info.Bucket) == 0 {
		return false, alert.CannotEmptyError("Bucket", "").HeaderInsertDesc("Check Exist")
	}
	if len(info.Key) == 0 {
		return false, alert.CannotEmptyError("Key", "").HeaderInsertDesc("Check Exist")
	}

	entry, sErr := bucketManager.Stat(info.Bucket, info.Key)
	if sErr != nil {
		if v, ok := sErr.(*storage.ErrorInfo); !ok {
			return false, data.NewEmptyError().AppendDescF("check file exists error, %s", sErr.Error())
		} else {
			if v.Code != 612 {
				return true, nil
			} else {
				return false, data.NewEmptyError().AppendDescF("check file exists error, %s", v.Err)
			}
		}
	}

	log.DebugF("Check [%s:%s] Exist, FileHash:%s PutTime:%d", info.Bucket, info.Key, entry.Hash, entry.PutTime)
	if len(entry.Hash) == 0 {
		return false, nil
	}

	return true, nil
}
