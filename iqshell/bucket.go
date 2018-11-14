package iqshell

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GetRet struct {
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
	Expiry   int64  `json:"expires"`
	Version  string `json:"version"`
}

type EntryPath struct {
	Bucket string
	Key    string
}

type ChgmEntryPath struct {
	EntryPath
	MimeType string
}

type ChtypeEntryPath struct {
	EntryPath
	FileType int
}

type DeleteAfterDaysEntryPath struct {
	EntryPath
	DeleteAfterDays int
}

type RenameEntryPath MoveEntryPath

type MoveEntryPath struct {
	SrcEntry EntryPath
	DstEntry EntryPath
	Force    bool
}

type CopyEntryPath MoveEntryPath

type BucketManager struct {
	*storage.BucketManager
}

type BucketDomainsRet []struct {
	Domain string `json:"domain"`
	Tbl    string `json:"tbl"`
	Owner  int    `json:"owner"`
}

func (m *BucketManager) DomainsOfBucket(bucket string) (domains []string, err error) {
	ctx := context.WithValue(context.TODO(), "mac", m.Mac)
	var reqHost string

	scheme := "http://"
	if m.Cfg.UseHTTPS {
		scheme = "https://"
	}

	reqHost = fmt.Sprintf("%s%s", scheme, viper.GetString("hosts.api_host"))
	reqURL := fmt.Sprintf("%s/v7/domain/list?tbl=%v", reqHost, bucket)
	headers := http.Header{}
	ret := new(BucketDomainsRet)
	cErr := m.Client.Call(ctx, ret, "POST", reqURL, headers)
	if cErr != nil {
		err = cErr
		return
	}
	for _, d := range *ret {
		domains = append(domains, d.Domain)
	}
	return

}

func (m *BucketManager) MakePrivateDownloadLink(domainOfBucket, fileKey string) (fileUrl string) {

	publicUrl := fmt.Sprintf("http://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	privateUrl, _ := m.PrivateUrl(publicUrl, deadline)
	fileUrl = privateUrl
	return
}

func (m *BucketManager) PrivateUrl(publicUrl string, deadline int64) (finalUrl string, err error) {
	srcUri, pErr := url.Parse(publicUrl)
	if pErr != nil {
		err = pErr
		return
	}

	h := hmac.New(sha1.New, m.Mac.SecretKey)

	urlToSign := srcUri.String()
	if strings.Contains(publicUrl, "?") {
		urlToSign = fmt.Sprintf("%s&e=%d", urlToSign, deadline)
	} else {
		urlToSign = fmt.Sprintf("%s?e=%d", urlToSign, deadline)
	}
	h.Write([]byte(urlToSign))

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := m.Mac.AccessKey + ":" + sign
	finalUrl = fmt.Sprintf("%s&token=%s", urlToSign, token)
	return
}

func (m *BucketManager) GetMac() *qbox.Mac {
	return m.Mac
}

func (m *BucketManager) rsHost(bucket string) (rsHost string, err error) {
	zone, err := m.Zone(bucket)
	if err != nil {
		return
	}

	rsHost = zone.GetRsHost(m.Cfg.UseHTTPS)
	return
}

func (m *BucketManager) Get(bucket, key string, destFile string) (err error) {
	entryUri := strings.Join([]string{bucket, key}, ":")

	var (
		reqHost string
		reqErr  error
	)
	reqHost = viper.GetString("hosts.rs_host")
	if reqHost == "" {
		reqHost, reqErr = m.rsHost(bucket)
		if reqErr != nil {
			err = reqErr
			return
		}
	}
	if !strings.HasPrefix(reqHost, "http") {
		reqHost = "http://" + reqHost
	}
	url := strings.Join([]string{reqHost, "get", Encode(entryUri)}, "/")

	var data GetRet

	ctx := context.WithValue(context.TODO(), "mac", m.Mac)
	headers := http.Header{}

	err = storage.DefaultClient.Call(ctx, &data, "GET", url, headers)
	if err != nil {
		return
	}
	resp, err := storage.DefaultClient.DoRequest(context.Background(), "GET", data.URL, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		body, rerr := ioutil.ReadAll(resp.Body)
		if rerr != nil {
			return rerr
		}
		fmt.Fprintf(os.Stderr, "Qget: http respcode: %d, respbody: %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}
	if strings.ContainsRune(destFile, os.PathSeparator) {
		destFile = filepath.Base(destFile)
	}
	f, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	return
}

func (m *BucketManager) Saveas(publicUrl, saveBucket, saveKey string) (string, error) {
	uri, parseErr := url.Parse(publicUrl)
	if parseErr != nil {
		return "", parseErr
	}
	baseUrl := uri.Host + uri.RequestURI()
	saveEntry := saveBucket + ":" + saveKey
	encodedSaveEntry := base64.URLEncoding.EncodeToString([]byte(saveEntry))
	baseUrl += "|saveas/" + encodedSaveEntry
	mac := m.GetMac()
	h := hmac.New(sha1.New, mac.SecretKey)
	h.Write([]byte(baseUrl))
	sign := h.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return publicUrl + "|saveas/" + encodedSaveEntry + "/sign/" + mac.AccessKey + ":" + encodedSign, nil
}

func (m *BucketManager) BatchStat(entries []EntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIStat(entry.Bucket, entry.Key))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchDelete(entries []EntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIDelete(entry.Bucket, entry.Key))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchCopy(entries []CopyEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URICopy(entry.SrcEntry.Bucket, entry.SrcEntry.Key, entry.DstEntry.Bucket, entry.DstEntry.Key, entry.Force))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchMove(entries []MoveEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIMove(entry.SrcEntry.Bucket, entry.SrcEntry.Key, entry.DstEntry.Bucket, entry.DstEntry.Key, entry.Force))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchRename(entries []RenameEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIMove(entry.SrcEntry.Bucket, entry.SrcEntry.Key, entry.DstEntry.Bucket, entry.DstEntry.Key, entry.Force))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchChgm(entries []ChgmEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIChangeMime(entry.Bucket, entry.Key, entry.MimeType))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchChtype(entries []ChtypeEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIChangeType(entry.Bucket, entry.Key, entry.FileType))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchDeleteAfterDays(entries []DeleteAfterDaysEntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, storage.URIDeleteAfterDays(entry.Bucket, entry.Key, entry.DeleteAfterDays))
	}
	return m.Batch(ops)
}

func (m *BucketManager) BatchSign(urls []string, deadline int64) (ret []string, err error) {
	for _, url := range urls {
		finalUrl, pErr := m.PrivateUrl(url, deadline)
		if pErr != nil {
			err = pErr
			return
		}
		ret = append(ret, finalUrl)
	}
	return
}

// NewBucketManager 用来构建一个新的资源管理对象
func NewBucketManager(mac *qbox.Mac, cfg *storage.Config) *BucketManager {
	bm := storage.NewBucketManager(mac, cfg)

	return &BucketManager{
		BucketManager: bm,
	}
}

// NewBucketManagerEx 用来构建一个新的资源管理对象
func NewBucketManagerEx(mac *qbox.Mac, cfg *storage.Config, client *storage.Client) *BucketManager {
	bm := storage.NewBucketManagerEx(mac, cfg, client)
	return &BucketManager{
		BucketManager: bm,
	}
}

func GetBucketManager() *BucketManager {
	account, gErr := GetAccount()
	if gErr != nil {
		fmt.Fprintf(os.Stderr, "GetBucketManager: %v\n", gErr)
		os.Exit(1)
	}
	mac := qbox.NewMac(account.AccessKey, account.SecretKey)
	cfg := storage.Config{
		RsHost:  viper.GetString("hosts.rs_host"),
		ApiHost: viper.GetString("hosts.api_host"),
		RsfHost: viper.GetString("hosts.rsf_host"),
	}
	return NewBucketManager(mac, &cfg)
}

func GetUpHost(cfg *storage.Config, ak, bucket string) (upHost string, err error) {

	var zone *storage.Zone
	if cfg.Zone != nil {
		zone = cfg.Zone
	} else {
		if v, zoneErr := storage.GetZone(ak, bucket); zoneErr != nil {
			err = zoneErr
			return
		} else {
			zone = v
		}
	}

	scheme := "http://"
	if cfg.UseHTTPS {
		scheme = "https://"
	}

	host := zone.SrcUpHosts[0]
	if cfg.UseCdnDomains {
		host = zone.CdnUpHosts[0]
	}

	upHost = fmt.Sprintf("%s%s", scheme, host)
	return
}
