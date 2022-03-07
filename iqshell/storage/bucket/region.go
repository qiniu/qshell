package bucket

import (
	"github.com/qiniu/go-sdk/v7/storage"
)

//func Region(cfg *storage.Config, bucket string) (*storage.Region, error) {
//	if cfg.Region != nil {
//		return cfg.Region, nil
//	}
//
//	acc, err := workspace.GetAccount()
//	if err != nil {
//		return nil, errors.New("get account error:" + err.Error())
//	}
//
//	if len(acc.AccessKey) == 0 {
//		return nil, errors.New("can't get access key")
//	}
//
//	r, err := storage.GetRegion(acc.AccessKey, bucket)
//	if err != nil {
//		return nil, errors.New("get region error:" + err.Error())
//	}
//
//	return r, nil
//}

func Region(b string) (region *storage.Zone, err error) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	return bucketManager.Zone(b)
}
