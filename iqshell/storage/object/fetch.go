package object

import (
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type FetchApiInfo struct {
	Bucket  string
	Key     string
	FromUrl string
}

type FetchResult = storage.FetchRet

func Fetch(info FetchApiInfo) (result FetchResult, err error) {
	if len(info.Bucket) == 0 {
		return result, errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.FromUrl) == 0 {
		return result, errors.New(alert.CannotEmpty("from url", ""))
	}

	log.DebugF("fetch start: %s => [%s|%s]", info.FromUrl, info.Bucket, info.Key)
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return result, err
	}

	if len(info.Key) == 0 {
		result, err = bucketManager.FetchWithoutKey(info.FromUrl, info.Bucket)
	} else {
		result, err = bucketManager.Fetch(info.FromUrl, info.Bucket, info.Key)
	}
	log.DebugF("fetch   end: %s => [%s|%s]", info.FromUrl, info.Bucket, info.Key)
	return result, err
}

type AsyncFetchApiInfo struct {
	Url              string `json:"url"`
	Host             string `json:"host,omitempty"`
	Bucket           string `json:"bucket"`
	Key              string `json:"key,omitempty"`
	Md5              string `json:"md5,omitempty"`
	Etag             string `json:"etag,omitempty"`
	CallbackURL      string `json:"callbackurl,omitempty"`
	CallbackBody     string `json:"callbackbody,omitempty"`
	CallbackBodyType string `json:"callbackbodytype,omitempty"`
	FileType         int    `json:"file_type,omitempty"`
	IgnoreSameKey    bool   `json:"ignore_same_key"` // false: 如果空间中已经存在同名文件则放弃本次抓取(仅对比 Key，不校验文件内容), true: 有同名会抓取
}

type AsyncFetchApiResult struct {
	Id   string `json:"id"`
	Wait int    `json:"wait"`
}

func (result AsyncFetchApiResult) String() string {
	return fmt.Sprintf(`{"id":"%s", "wait":%d}`, result.Id, result.Wait)
}

func AsyncFetch(info AsyncFetchApiInfo) (result AsyncFetchApiResult, err error) {
	bm, err := bucket.GetBucketManager()
	if err != nil {
		return result, err
	}
	reqUrl, err := bm.ApiReqHost(info.Bucket)
	if err != nil {
		return result, err
	}

	reqUrl += "/sisyphus/fetch"

	err = bm.Client.CredentialedCallWithJson(context.Background(), bm.Mac, auth.TokenQiniu, &result, "POST", reqUrl, nil, info)
	return result, err
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
