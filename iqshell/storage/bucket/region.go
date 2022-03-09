package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
)

func Region(b string) (region *storage.Zone, err error) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	return bucketManager.Zone(b)
}
