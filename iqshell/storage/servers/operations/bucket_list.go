package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/servers"
)

type ListInfo servers.ListApiInfo

func (info *ListInfo) Check() *data.CodeError {
	if info.Limit <= 0 {
		info.Limit = 50
	}
	return nil
}

// List 列举所有 bucket
func List(cfg *iqshell.Config, info ListInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	log.AlertF("%s", servers.BucketInfoDetailDescriptionStringFormat())
	servers.AllBuckets(servers.ListApiInfo(info), func(bucket *servers.BucketInfo, err *data.CodeError) {
		if err != nil {
			data.SetCmdStatusError()
			log.ErrorF("Get buckets error: %v", err)
			return
		}

		if bucket == nil {
			return
		}

		if info.Detail {
			log.AlertF("%s", bucket.DetailDescriptionString())
		} else {
			log.AlertF("%s", bucket.DescriptionString())
		}
	})
}
