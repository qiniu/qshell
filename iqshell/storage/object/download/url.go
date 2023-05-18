package download

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"net/url"
	"strings"
	"time"
)

type UrlApiInfo struct {
	BucketDomain string
	Key          string
	UseHttps     bool
}

// PublicUrl 返回公有空间的下载链接，不可以用于私有空间的下载
func PublicUrl(info UrlApiInfo) (fileUrl string) {
	domain := utils.RemoveUrlScheme(info.BucketDomain)
	if info.UseHttps {
		fileUrl = fmt.Sprintf("https://%s/%s", domain, url.PathEscape(info.Key))
	} else {
		fileUrl = fmt.Sprintf("http://%s/%s", domain, url.PathEscape(info.Key))
	}
	return
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
	deadline := time.Now().Add(time.Minute * 24 * 30).Unix()
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

		// 源站域名需要签名
		if !info.IsPublicBucket || isIoSrcHost(info.Host, info.Key) {
			if u, e := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
				PublicUrl: urlString,
				Deadline:  5*60 + time.Now().Unix(),
			}); e != nil {
				return "", e
			} else {
				urlString = u.Url
			}
		}
	}
	return urlString, nil
}

// CreateSrcDownloadDomainWithBucket bucket 源站下载域名
func CreateSrcDownloadDomainWithBucket(cfg *config.Config, bucketName string) ([]string, *data.CodeError) {

	hosts := make([]string, 0, 0)
	if cfg != nil && len(cfg.GetIoSrcHost()) > 0 {
		hosts = append(hosts, cfg.GetIoSrcHost())
	}

	serverCfgHost, err := getSrcDownloadDomainWithBucket(bucketName)
	if err != nil {
		log.WarningF("get io src host for bucket:%s error:%v", bucketName, err)
	}
	if len(serverCfgHost) > 0 {
		hosts = append(hosts, serverCfgHost)
	}
	return hosts, nil
}

func isIoSrcHost(host string, bucketName string) bool {
	host = utils.RemoveUrlScheme(host)
	if len(host) == 0 {
		return false
	}

	customEndpoint := ""
	if workspace.GetConfig() != nil {
		customEndpoint = workspace.GetConfig().GetIoSrcHost()
	}
	if len(customEndpoint) > 0 {
		return strings.Contains(host, customEndpoint)
	}

	srcDownloadDomain, _ := getSrcDownloadDomainWithBucket(bucketName)
	if len(srcDownloadDomain) == 0 {
		return false
	}

	return strings.Contains(host, srcDownloadDomain)
}

func getSrcDownloadDomainWithBucket(bucketName string) (string, *data.CodeError) {
	region, err := bucket.Region(bucketName)
	if err != nil {
		return "", err
	}

	if len(region.IoSrcHost) == 0 {
		return "", data.NewEmptyError().AppendDesc("io src is empty")
	}

	return region.IoSrcHost, nil
}
