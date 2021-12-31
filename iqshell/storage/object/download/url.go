package download

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"net/url"
	"strings"
	"time"
)


type PublicUrlApiInfo struct {
	BucketDomain string
	Key          string
	UseHttps     bool
}

// PublicUrl 返回公有空间的下载链接，不可以用于私有空间的下载
func PublicUrl(info PublicUrlApiInfo) (fileUrl string) {
	if info.UseHttps {
		fileUrl = fmt.Sprintf("https://%s/%s", info.BucketDomain, url.PathEscape(info.Key))
	} else {
		fileUrl = fmt.Sprintf("http://%s/%s", info.BucketDomain, url.PathEscape(info.Key))
	}
	return
}

// 私有下载链接

type PublicUrlToPrivateApiInfo struct {
	PublicUrl string
	Deadline  int64
}

// PublicUrlToPrivate 公转私
func PublicUrlToPrivate(info PublicUrlToPrivateApiInfo) (finalUrl string, err error) {
	if len(info.PublicUrl) == 0 {
		return "", errors.New(alert.CannotEmpty("url", ""))
	}

	if info.Deadline < 1 {
		return "", errors.New("deadline is invalid")
	}

	m, err := bucket.GetBucketManager()
	if err != nil {
		return "", err
	}

	srcUri, pErr := url.Parse(info.PublicUrl)
	if pErr != nil {
		err = pErr
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
	finalUrl = fmt.Sprintf("%s&token=%s", urlToSign, token)
	return
}

// 私有下载 url

type PrivateUrlApiInfo PublicUrlApiInfo

// PrivateUrl 返回私有空间的下载链接， 也可以用于公有空间的下载
func PrivateUrl(info PrivateUrlApiInfo) (fileUrl string) {
	publicUrl := PublicUrl(PublicUrlApiInfo(info))
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	privateUrl, _ := PublicUrlToPrivate(PublicUrlToPrivateApiInfo{
		PublicUrl: publicUrl,
		Deadline:  deadline,
	})
	fileUrl = privateUrl
	return
}
