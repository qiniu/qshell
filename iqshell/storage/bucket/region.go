package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func Region(b string) (*storage.Zone, *data.CodeError) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	region, e := bucketManager.Zone(b)
	return region, data.ConvertError(e)
}
