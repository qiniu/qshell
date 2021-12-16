package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type FetchApiInfo struct {
	Bucket  string
	Key     string
	FromUrl string
}

type FetchResult storage.FetchRet

func Fetch(info FetchApiInfo) (FetchResult, error) {
	if len(info.Bucket) == 0 {
		return FetchResult{}, errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.Key) == 0 {
		key, err := utils.KeyFromUrl(info.FromUrl)
		if err != nil || len(key) == 0 {
			return FetchResult{}, errors.New("get key from url failed:" + err.Error())
		}
		info.Key = key
	}

	if len(info.FromUrl) == 0 {
		return FetchResult{}, errors.New(alert.CannotEmpty("from url", ""))
	}

	log.InfoF("fetch info: bucket:%s key:%s fromUrl:%s", info.Bucket, info.Key, info.FromUrl)

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return FetchResult{}, err
	}
	fetchResult, err := bucketManager.Fetch(info.FromUrl, info.Bucket, info.Key)
	return FetchResult(fetchResult), err
}
