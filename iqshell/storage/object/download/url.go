package download

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
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
	publicUrl := PublicUrl(UrlApiInfo(info))
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
