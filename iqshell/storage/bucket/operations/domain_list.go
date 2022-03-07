package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"os"
)

type ListDomainInfo struct {
	Bucket string
	Detail bool
}

func (info *ListDomainInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func ListDomains(cfg *iqshell.Config, info ListDomainInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	domains, err := bucket.AllDomainsOfBucket(info.Bucket)
	if err != nil {
		log.Error("Get domains error: ", err)
		os.Exit(data.StatusError)
	} else {
		if len(domains) == 0 {
			log.ErrorF("No domains found for bucket `%s`\n", info.Bucket)
		} else {
			if info.Detail {
				for _, domain := range domains {
					log.Alert(domain.DetailDescriptionString())
				}
			} else {
				for _, domain := range domains {
					log.Alert(domain.DescriptionString())
				}
			}
		}
	}
}
