package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type GetBucketInfo struct {
	Bucket string
}

func (i *GetBucketInfo) Check() *data.CodeError {
	if len(i.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func GetBucket(cfg *iqshell.Config, info GetBucketInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if bucketInfo, err := bucket.GetBucketInfo(bucket.GetBucketApiInfo{
		Bucket: info.Bucket,
	}); err != nil {
		log.ErrorF("get bucket(%s) info error:%v", info.Bucket, err)
	} else {
		log.AlertF("%-10s:%s", "Bucket", info.Bucket)
		log.AlertF("%-10s:%s", "RegionID", bucketInfo.Region)
		log.AlertF("%-10s:%v", "Private", bucketInfo.Private > 0)
	}
}
