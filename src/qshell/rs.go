package qshell

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"net/url"
	"strings"
)

type FetchResult struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

type ChgmEntryPath struct {
	Bucket   string
	Key      string
	MimeType string
}

type BatchItemRet struct {
	Code int              `json:"code"`
	Data BatchItemRetData `json:"data"`
}

type BatchItemRetData struct {
	Error string `json:"error,omitempty"`
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

func Saveas(mac *digest.Mac, publicUrl string, saveBucket string, saveKey string) (string, error) {
	uri, parseErr := url.Parse(publicUrl)
	if parseErr != nil {
		return "", parseErr
	}
	baseUrl := uri.Host + uri.RequestURI()
	saveEntry := saveBucket + ":" + saveKey
	encodedSaveEntry := base64.URLEncoding.EncodeToString([]byte(saveEntry))
	baseUrl += "|saveas/" + encodedSaveEntry
	h := hmac.New(sha1.New, mac.SecretKey)
	h.Write([]byte(baseUrl))
	sign := h.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return publicUrl + "|saveas/" + encodedSaveEntry + "/sign/" + mac.AccessKey + ":" + encodedSign, nil
}

func BatchChgm(client rs.Client, entries []ChgmEntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIChangeMime(e.Bucket, e.Key, e.MimeType)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchDelete(client rs.Client, entries []rs.EntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIDelete(e.Bucket, e.Key)
	}
	err = client.Batch(nil, &ret, b)
	return
}
