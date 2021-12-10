package rs

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
)

type PrivateUrlApiInfo struct {
	PublicUrl string
	Deadline int64
}

func PrivateUrl(info PrivateUrlApiInfo) (finalUrl string, err error) {
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
