package bucket

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"strings"
)

// AllDomainsOfBucket 获取一个存储空间绑定的CDN域名
func AllDomainsOfBucket(bucket string) (domains []string, err error) {
	bucketManager, gErr := GetBucketManager()
	if gErr != nil {
		return nil, gErr
	}

	infos, err := bucketManager.ListBucketDomains(bucket)
	if err != nil {
		if e, ok := err.(*storage.ErrorInfo); ok {
			if e.Code != 404 {
				return
			}
			err = nil
		} else {
			return
		}
	}

	for _, d := range infos {
		domains = append(domains, d.Domain)
	}
	return
}

func DomainOfBucket(bucket string) (domain string, err error) {
	//get domains of bucket
	domainsOfBucket, gErr := AllDomainsOfBucket(bucket)
	if gErr != nil {
		err = fmt.Errorf("Get domains of bucket error: %v", gErr)
		return
	}

	if len(domainsOfBucket) == 0 {
		err = fmt.Errorf("No domains found for bucket: %s", bucket)
		return
	}

	for _, d := range domainsOfBucket {
		if !strings.HasPrefix(d, ".") {
			domain = d
			break
		}
	}
	return
}
