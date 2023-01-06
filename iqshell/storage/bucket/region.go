package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

func Region(b string) (*storage.Zone, *data.CodeError) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	region, e := bucketManager.Zone(b)
	return region, data.ConvertError(e)
}

func CompleteBucketManagerRegion(bucketManager *storage.BucketManager, bucket string) *data.CodeError {
	if bucketManager == nil {
		return data.NewEmptyError().AppendDesc("bucketManager is empty")
	}

	if bucketManager.Cfg == nil {
		return data.NewEmptyError().AppendDesc("bucketManager.Cfg is empty")
	}

	if bucketManager.Cfg.Region != nil && bucketManager.Cfg.Zone != nil && len(bucketManager.Cfg.CentralRsHost) != 0 {
		return nil
	}

	region, e := storage.GetZone(bucketManager.Mac.AccessKey, bucket)
	if e != nil {
		return data.ConvertError(e)
	}
	bucketManager.Cfg.CentralRsHost = utils.RemoveUrlScheme(region.RsHost)
	bucketManager.Cfg.Region = region
	bucketManager.Cfg.Zone = region
	return nil
}
