package qshell

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"strings"
)

type FetchResult struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

func Fetch(mac *digest.Mac, remoteResUrl, bucket, key string) (fetchResult FetchResult, err error) {
	client := rs.New(mac)
	entry := bucket
	if key != "" {
		entry += ":" + key
	}
	fetchUri := fmt.Sprintf("/fetch/%s/to/%s",
		base64.URLEncoding.EncodeToString([]byte(remoteResUrl)),
		base64.URLEncoding.EncodeToString([]byte(entry)))
	err = client.Conn.Call(nil, &fetchResult, conf.IO_HOST+fetchUri)
	return
}

func Prefetch(mac *digest.Mac, bucket, key string) (err error) {
	client := rs.New(mac)
	prefetchUri := fmt.Sprintf("/prefetch/%s", base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)))
	err = client.Conn.Call(nil, nil, conf.IO_HOST+prefetchUri)
	return
}

func PrivateUrl(mac *digest.Mac, publicUrl string, deadline int64) string {
	h := hmac.New(sha1.New, mac.SecretKey)

	urlToSign := publicUrl
	if strings.Contains(publicUrl, "?") {
		urlToSign = fmt.Sprintf("%s&e=%d", urlToSign, deadline)
	} else {
		urlToSign = fmt.Sprintf("%s?e=%d", urlToSign, deadline)
	}
	h.Write([]byte(urlToSign))

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := mac.AccessKey + ":" + sign
	url := fmt.Sprintf("%s&token=%s", urlToSign, token)
	return url
}
