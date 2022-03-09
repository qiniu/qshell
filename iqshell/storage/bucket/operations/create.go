package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type CreateInfo bucket.CreateApiInfo

func (i *CreateInfo) Check() error {
	if len(i.RegionId) == 0 {
		return alert.CannotEmptyError("RegionId", "")
	}
	if len(i.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func Create(cfg *iqshell.Config, info CreateInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if err := bucket.Create(bucket.CreateApiInfo(info)); err != nil {
		log.ErrorF("bucket:%s-%s create error:%v", info.RegionId, info.Bucket, err)
	} else {
		log.AlertF("bucket:%s-%s create success", info.RegionId, info.Bucket)
	}
}
