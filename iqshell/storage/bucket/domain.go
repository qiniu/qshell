package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
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
