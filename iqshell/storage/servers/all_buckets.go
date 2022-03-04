package servers

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

// AllBuckets List list 所有 bucket
func AllBuckets(shared bool) (buckets []string, err error) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}
	return bucketManager.Buckets(shared)
}

func BucketRegion(b string) (region *storage.Zone, err error) {
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	return bucketManager.Zone(b)
}