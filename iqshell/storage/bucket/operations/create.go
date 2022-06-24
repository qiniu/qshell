package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type CreateInfo struct {
	RegionId string
	Bucket   string
	Private  bool
}

func (i *CreateInfo) Check() *data.CodeError {
	if len(i.RegionId) == 0 {
		return alert.CannotEmptyError("Region", "")
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

	if err := bucket.Create(bucket.CreateApiInfo{
		RegionId: info.RegionId,
		Bucket:   info.Bucket,
		Private:  info.Private,
	}); err != nil {
		log.ErrorF("bucket:%s-%s create error:%v", info.RegionId, info.Bucket, err)
	} else {
		log.AlertF("bucket:%s-%s create success", info.RegionId, info.Bucket)
	}
}
