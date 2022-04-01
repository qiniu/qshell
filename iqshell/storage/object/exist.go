package object

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type ExistApiInfo struct {
	Bucket string
	Key    string
}

func Exist(info ExistApiInfo) (exists bool, err *data.CodeError) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return false, err
	}

	entry, sErr := bucketManager.Stat(info.Bucket, info.Key)
	if sErr != nil {
		if v, ok := sErr.(*storage.ErrorInfo); !ok {
			err = data.NewEmptyError().AppendDescF("check file exists error, %s", sErr.Error())
			return
		} else {
			if v.Code != 612 {
				err = data.NewEmptyError().AppendDescF("check file exists error, %s", v.Err)
				return
			} else {
				exists = false
				return
			}
		}
	}
	if entry.Hash != "" {
		exists = true
	}
	return
}
