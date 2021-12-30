package object

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
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

type AsyncFetchApiInfo storage.AsyncFetchParam
type AsyncFetchApiResult storage.AsyncFetchRet

func AsyncFetch(info AsyncFetchApiInfo) (AsyncFetchApiResult, error) {
	bm, err := bucket.GetBucketManager()
	if err != nil {
		return AsyncFetchApiResult{}, err
	}
	ret, err := bm.AsyncFetch(storage.AsyncFetchParam(info))
	return AsyncFetchApiResult(ret), err
}

func CheckAsyncFetchStatus(toBucket, id string) (ret AsyncFetchApiResult, err error) {
	bm, gErr := bucket.GetBucketManager()
	if gErr != nil {
		err = gErr
		return
	}

	reqUrl, aErr := bm.ApiReqHost(toBucket)
	if aErr != nil {
		err = aErr
		return
	}

	mac, gErr := workspace.GetMac()
	if gErr != nil {
		err = gErr
		return
	}

	reqUrl += ("/sisyphus/fetch?id=" + id)
	ctx := auth.WithCredentialsType(workspace.GetContext(), mac, auth.TokenQiniu)
	err = bm.Client.Call(ctx, &ret, "GET", reqUrl, nil)
	return
}
