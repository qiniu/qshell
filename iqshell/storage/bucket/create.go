package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
)

type CreateApiInfo struct {
	RegionId string
	Bucket   string
}

func Create(info CreateApiInfo) error {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return err
	}
	return bucketManager.CreateBucket(info.Bucket, storage.RegionID(info.RegionId))
}
