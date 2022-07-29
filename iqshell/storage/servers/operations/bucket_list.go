package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/servers"
)

type ListInfo struct {
	servers.ListApiInfo

	Detail bool
}

func (info *ListInfo) Check() *data.CodeError {
	return nil
}

// List list 所有 bucket
func List(cfg *iqshell.Config, info ListInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	buckets, err := servers.AllBuckets(info.ListApiInfo)
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Get buckets error: %v", err)
		return
	} else if len(buckets) == 0 {
		log.Warning("No buckets found")
		return
	}

	if info.Detail {
		for _, b := range buckets {
			log.AlertF("%s", b.DetailDescriptionString())
		}
	} else {
		for _, b := range buckets {
			log.AlertF("%s", b.DescriptionString())
		}
	}

}
