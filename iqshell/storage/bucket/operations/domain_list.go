package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"os"
)

type ListDomainInfo struct {
	Bucket string
}

func ListDomains(info ListDomainInfo) {
	domains, err := bucket.AllDomainsOfBucket(info.Bucket)
	if err != nil {
		log.Error("Get domains error: ", err)
		os.Exit(data.StatusError)
	} else {
		if len(domains) == 0 {
			log.ErrorF("No domains found for bucket `%s`\n", info.Bucket)
		} else {
			for _, domain := range domains {
				log.Alert(domain)
			}
		}
	}
}
