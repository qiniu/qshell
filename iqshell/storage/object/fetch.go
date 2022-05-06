package object

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type FetchApiInfo struct {
	Bucket  string
	Key     string
	FromUrl string
}

func (i *FetchApiInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", i.Bucket, i.Key, i.FromUrl)
}

type FetchResult storage.FetchRet

var _ flow.Result = (*FetchResult)(nil)

func (a *FetchResult) IsValid() bool {
	return len(a.Key) > 0 && len(a.MimeType) > 0 && len(a.Hash) > 0
}

func Fetch(info FetchApiInfo) (*FetchResult, *data.CodeError) {

	if len(info.Bucket) == 0 {
		return nil, alert.CannotEmptyError("bucket", "")
	}

	if len(info.FromUrl) == 0 {
		return nil, alert.CannotEmptyError("from url", "")
	}

	log.DebugF("fetch start: %s => [%s:%s]", info.FromUrl, info.Bucket, info.Key)
	bucketManager, e := bucket.GetBucketManager()
	if e != nil {
		return nil, e
	}

	var err error
	var result storage.FetchRet
	if len(info.Key) == 0 {
		result, err = bucketManager.FetchWithoutKey(info.FromUrl, info.Bucket)
	} else {
		result, err = bucketManager.Fetch(info.FromUrl, info.Bucket, info.Key)
	}
	log.DebugF("fetch   end: %s => [%s:%s]", info.FromUrl, info.Bucket, info.Key)
	return (*FetchResult)(&result), data.ConvertError(err)
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

func (result *AsyncFetchApiResult) IsValid() bool {
	return len(result.Id) > 0
}

func (result *AsyncFetchApiResult) String() string {
	return fmt.Sprintf(`{"id":"%s", "wait":%d}`, result.Id, result.Wait)
}

func AsyncFetch(info AsyncFetchApiInfo) (result *AsyncFetchApiResult, err *data.CodeError) {
	bm, err := bucket.GetBucketManager()
	if err != nil {
		return result, err
	}
	reqUrl, e := bm.ApiReqHost(info.Bucket)
	if e != nil {
		return result, data.ConvertError(e)
	}

	reqUrl += "/sisyphus/fetch"

	result = &AsyncFetchApiResult{}
	e = bm.Client.CredentialedCallWithJson(context.Background(), bm.Mac, auth.TokenQiniu, result, "POST", reqUrl, nil, info)
	return result, data.ConvertError(e)
}

func CheckAsyncFetchStatus(toBucket, id string) (ret AsyncFetchApiResult, err *data.CodeError) {
	bm, gErr := bucket.GetBucketManager()
	if gErr != nil {
		err = gErr
		return
	}

	reqUrl, aErr := bm.ApiReqHost(toBucket)
	if aErr != nil {
		err = data.ConvertError(aErr)
		return
	}

	mac, gErr := workspace.GetMac()
	if gErr != nil {
		err = gErr
		return
	}

	reqUrl += ("/sisyphus/fetch?id=" + id)
	ctx := auth.WithCredentialsType(workspace.GetContext(), mac, auth.TokenQiniu)
	err = data.ConvertError(bm.Client.Call(ctx, &ret, "GET", reqUrl, nil))
	return
}
