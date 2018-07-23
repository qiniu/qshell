package qshell

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/conf"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"net/url"
	"strings"
)

type FetchResult struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

type ChgmEntryPath struct {
	Bucket   string
	Key      string
	MimeType string
}

type ChtypeEntryPath struct {
	Bucket   string
	Key      string
	FileType int
}

type DeleteAfterDaysEntryPath struct {
	Bucket          string
	Key             string
	DeleteAfterDays int
}

type RenameEntryPath struct {
	Bucket string
	OldKey string
	NewKey string
}

type MoveEntryPath struct {
	SrcBucket  string
	DestBucket string
	SrcKey     string
	DestKey    string
}

type CopyEntryPath struct {
	SrcBucket  string
	DestBucket string
	SrcKey     string
	DestKey    string
}

type BatchItemRet struct {
	Code int              `json:"code"`
	Data BatchItemRetData `json:"data"`
}

type BatchItemRetData struct {
	Fsize    int    `json:"fsize,omitempty"`
	Hash     string `json:"hash,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	PutTime  int64  `json:"putTime,omitempty"`
	Error    string `json:"error,omitempty"`
	FileType int    `json:"type"`
}

func Fetch(mac *digest.Mac, remoteResUrl, bucket, key string) (fetchResult FetchResult, err error) {
	client := rs.NewMac(mac)
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
	client := rs.NewMac(mac)
	prefetchUri := fmt.Sprintf("/prefetch/%s", base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)))
	err = client.Conn.Call(nil, nil, conf.IO_HOST+prefetchUri)
	return
}

func PrivateUrl(mac *digest.Mac, publicUrl string, deadline int64) (finalUrl string, err error) {
	srcUri, pErr := url.Parse(publicUrl)
	if pErr != nil {
		err = pErr
		return
	}

	h := hmac.New(sha1.New, mac.SecretKey)

	urlToSign := srcUri.String()
	if strings.Contains(publicUrl, "?") {
		urlToSign = fmt.Sprintf("%s&e=%d", urlToSign, deadline)
	} else {
		urlToSign = fmt.Sprintf("%s?e=%d", urlToSign, deadline)
	}
	h.Write([]byte(urlToSign))

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := mac.AccessKey + ":" + sign
	finalUrl = fmt.Sprintf("%s&token=%s", urlToSign, token)
	return
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

func BatchStat(client rs.Client, entries []rs.EntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIStat(e.Bucket, e.Key)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchChgm(client rs.Client, entries []ChgmEntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIChangeMime(e.Bucket, e.Key, e.MimeType)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchChtype(client rs.Client, entries []ChtypeEntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIChangeType(e.Bucket, e.Key, e.FileType)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchDeleteAfterDays(client rs.Client, entries []DeleteAfterDaysEntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIDeleteAfterDays(e.Bucket, e.Key, e.DeleteAfterDays)
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

func BatchRename(client rs.Client, entries []RenameEntryPath, force bool) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIMove(e.Bucket, e.OldKey, e.Bucket, e.NewKey, force)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchMove(client rs.Client, entries []MoveEntryPath, force bool) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URIMove(e.SrcBucket, e.SrcKey, e.DestBucket, e.DestKey, force)
	}
	err = client.Batch(nil, &ret, b)
	return
}

func BatchCopy(client rs.Client, entries []CopyEntryPath, force bool) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = rs.URICopy(e.SrcBucket, e.SrcKey, e.DestBucket, e.DestKey, force)
	}
	err = client.Batch(nil, &ret, b)
	return
}