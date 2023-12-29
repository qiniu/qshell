package download

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type UrlApiInfo struct {
	BucketDomain string
	Key          string
	UseHttps     bool
}

// PublicUrl 返回公有空间的下载链接，不可以用于私有空间的下载
func PublicUrl(info UrlApiInfo) (fileUrl string) {
	domain := utils.Endpoint(info.UseHttps, info.BucketDomain)
	return fmt.Sprintf("%s/%s", domain, url.PathEscape(info.Key))
}

// PublicUrlToPrivateApiInfo 私有下载链接
type PublicUrlToPrivateApiInfo struct {
	PublicUrl string
	Deadline  int64
}

type PublicUrlToPrivateApiResult struct {
	Url string
}

var _ flow.Result = (*PublicUrlToPrivateApiResult)(nil)

func (p *PublicUrlToPrivateApiResult) IsValid() bool {
	return len(p.Url) > 0
}

// PublicUrlToPrivate 公转私
func PublicUrlToPrivate(info PublicUrlToPrivateApiInfo) (result *PublicUrlToPrivateApiResult, err *data.CodeError) {
	if len(info.PublicUrl) == 0 {
		return nil, alert.CannotEmptyError("url", "")
	}

	if info.Deadline < 1 {
		return nil, data.NewEmptyError().AppendDesc("deadline is invalid")
	}

	m, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	srcUri, pErr := url.Parse(info.PublicUrl)
	if pErr != nil {
		err = data.ConvertError(pErr)
		return
	}

	h := hmac.New(sha1.New, m.Mac.SecretKey)

	urlToSign := srcUri.String()
	if strings.Contains(info.PublicUrl, "?") {
		urlToSign = fmt.Sprintf("%s&e=%d", urlToSign, info.Deadline)
	} else {
		urlToSign = fmt.Sprintf("%s?e=%d", urlToSign, info.Deadline)
	}
	h.Write([]byte(urlToSign))

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := m.Mac.AccessKey + ":" + sign
	return &PublicUrlToPrivateApiResult{
		Url: fmt.Sprintf("%s&token=%s", urlToSign, token),
	}, nil
}

// PrivateUrl 返回私有空间的下载链接， 也可以用于公有空间的下载
func PrivateUrl(info UrlApiInfo) (fileUrl string) {
	publicUrl := PublicUrl(info)
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	result, _ := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
		PublicUrl: publicUrl,
		Deadline:  deadline,
	})
	if result != nil {
		fileUrl = result.Url
	}
	return
}

// 下载 Url
func createDownloadUrl(info *DownloadApiInfo) (string, *data.CodeError) {
	urlString := ""
	useHttps := workspace.GetConfig().IsUseHttps()

	// 构造下载 url
	if info.UseGetFileApi {
		mac, err := workspace.GetMac()
		if err != nil {
			return "", data.NewEmptyError().AppendDescF("download get mac error:%v", mac)
		}
		urlString = utils.Endpoint(useHttps, info.Host)
		urlString = strings.Join([]string{urlString, "getfile", mac.AccessKey, info.Bucket, url.PathEscape(info.Key)}, "/")
	} else {
		urlString = PublicUrl(UrlApiInfo{
			BucketDomain: info.Host,
			Key:          info.Key,
			UseHttps:     useHttps,
		})

		isIoSrc, iErr := isIoSrcHost(info.Host, info.Bucket)
		if iErr != nil {
			log.WarningF("check host(bucket:%s) is src host error:%v", info.Bucket, iErr)
		}
		// 源站域名需要签名
		if !info.IsPublicBucket || isIoSrc {
			// 是源站域名，但也可能是 私有空间，所以此处逻辑不能对调
			expire := 60 * time.Minute
			if isIoSrc {
				expire = 3 * time.Minute
			}

			if u, e := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
				PublicUrl: urlString,
				Deadline:  time.Now().Add(expire).Unix(),
			}); e != nil {
				return "", e
			} else {
				urlString = u.Url
			}
		}
	}
	return urlString, nil
}

func isIoSrcHost(host string, bucketName string) (bool, *data.CodeError) {
	host = utils.RemoveUrlScheme(host)
	if len(host) == 0 {
		return false, nil
	}

	srcDownloadDomain, err := GetBucketIoSrcDomain(bucketName)
	if err != nil {
		return false, err
	}
	return host == srcDownloadDomain, nil
}

func GetBucketIoSrcDomain(bucketName string) (string, *data.CodeError) {
	region, err := bucket.Region(bucketName)
	if err != nil {
		return "", err
	}

	if len(region.IoSrcHost) == 0 {
		return "", data.NewEmptyError().AppendDesc("io src is empty")
	}

	return region.IoSrcHost, nil
}
