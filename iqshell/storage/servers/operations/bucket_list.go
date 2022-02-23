package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/servers"
	"os"
)

type ListInfo struct {
	Shared bool
}

func (info *ListInfo) Check() error {
	return nil
}

// List list 所有 bucket
func List(cfg *iqshell.Config, info ListInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	buckets, err := servers.AllBuckets(info.Shared)
	if err != nil {
		log.ErrorF("Get buckets error: %v", err)
		os.Exit(data.StatusError)
	} else if len(buckets) == 0 {
		log.Warning("No buckets found")
		return
	}

	for _, b := range buckets {
		log.Alert(b)
	}
}
