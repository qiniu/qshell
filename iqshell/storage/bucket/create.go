package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

type CreateApiInfo struct {
	RegionId string
	Bucket   string
}

func Create(info CreateApiInfo) *data.CodeError {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return err
	}
	e := bucketManager.CreateBucket(info.Bucket, storage.RegionID(info.RegionId))
	return data.ConvertError(e)
}
