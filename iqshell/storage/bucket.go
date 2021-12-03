package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/conf"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Get 接口返回的结构
type GetRet struct {
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
	Expiry   int64  `json:"expires"`
	Version  string `json:"version"`
}

type EntryPath struct {
	Bucket  string
	Key     string
	PutTime string
}

// 改变文件mime需要的信息
type ChgmEntryPath struct {
	EntryPath
	MimeType string
}

// 改变文件存储类型需要的信息
type ChtypeEntryPath struct {
	EntryPath
	FileType int
}

// 设置deleteAfterDays需要的参数
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

// 获取一个存储空间绑定的CDN域名
func (m *BucketManager) DomainsOfBucket(bucket string) (domains []string, err error) {
	infos, err := m.ListBucketDomains(bucket)
	if err != nil {
		if e, ok := err.(*storage.ErrorInfo); ok {
			if e.Code != 404 {
				return
			}
			err = nil
		} else {
			return
		}
	}
	for _, d := range infos {
		domains = append(domains, d.Domain)
	}
	return
}

// 返回公有空间的下载链接，不可以用于私有空间的下载
func (m *BucketManager) MakePublicDownloadLink(domainOfBucket, fileKey string, useHttps bool) (fileUrl string) {
	if useHttps {
		fileUrl = fmt.Sprintf("https://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	} else {
		fileUrl = fmt.Sprintf("http://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	}
	return
}

// 返回私有空间的下载链接， 也可以用于公有空间的下载
func (m *BucketManager) MakePrivateDownloadLink(domainOfBucket, fileKey string, useHttps bool) (fileUrl string) {
	var publicUrl string
	if useHttps {
		publicUrl = fmt.Sprintf("https://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	} else {
		publicUrl = fmt.Sprintf("http://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	}
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	privateUrl, _ := m.PrivateUrl(publicUrl, deadline)
	fileUrl = privateUrl
	return
}

// 返回私有空间的下载链接， 也可以用于公有空间的下载
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

// 从存储空间下载文件（不需要绑定CDN域名）
func (m *BucketManager) Get(bucket, key string, destFile string) (err error) {
	entryUri := strings.Join([]string{bucket, key}, ":")

	var (
		reqHost string
		reqErr  error
	)
	reqHost = config.RsHost()
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
	url := strings.Join([]string{reqHost, "get", utils.Encode(entryUri)}, "/")

	var data GetRet

	ctx := auth.WithCredentials(context.Background(), m.Mac)
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

	for {
		f, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err == nil {
			defer f.Close()
			if _, err := io.Copy(f, resp.Body); err != nil {
				fmt.Fprintf(os.Stderr, "Qget: err: %s\n", err)
			}
			break
		} else if os.IsNotExist(err) {
			destDir := filepath.Dir(destFile)
			if err := os.MkdirAll(destDir, 0700); err != nil {
				fmt.Fprintf(os.Stderr, "Qget: err: %s\n", err)
				break
			}
		} else {
			fmt.Fprintf(os.Stderr, "Qget: err: %s\n", err)
			break
		}
	}

	return
}

func (m *BucketManager) CheckExists(bucket, key string) (exists bool, err error) {
	entry, sErr := m.Stat(bucket, key)
	if sErr != nil {
		if v, ok := sErr.(*storage.ErrorInfo); !ok {
			err = fmt.Errorf("Check file exists error, %s", sErr.Error())
			return
		} else {
			if v.Code != 612 {
				err = fmt.Errorf("Check file exists error, %s", v.Err)
				return
			} else {
				exists = false
				return
			}
		}
	}
	if entry.Hash != "" {
		exists = true
	}
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

func batchURICondition(entry EntryPath) string {
	cond := ""
	if entry.PutTime != "" {
		cond += "putTime=" + entry.PutTime
	}
	if cond == "" {
		return ""
	}
	return fmt.Sprintf("/cond/%s", base64.URLEncoding.EncodeToString([]byte(cond)))
}
func batchURIDelete(entry EntryPath) string {
	return fmt.Sprintf("/delete/%s%s", storage.EncodedEntry(entry.Bucket, entry.Key), batchURICondition(entry))
}
func (m *BucketManager) BatchDelete(entries []EntryPath) (ret []storage.BatchOpRet, err error) {
	ops := make([]string, 0, len(entries))
	for _, entry := range entries {
		ops = append(ops, batchURIDelete(entry))
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

func (m *BucketManager) CheckAsyncFetchStatus(bucket, id string) (ret storage.AsyncFetchRet, err error) {

	reqUrl, err := m.ApiReqHost(bucket)
	if err != nil {
		return
	}

	reqUrl += ("/sisyphus/fetch?id=" + id)

	ctx := auth.WithCredentialsType(context.Background(), m.Mac, auth.TokenQiniu)
	err = m.Client.Call(ctx, &ret, "GET", reqUrl, nil)
	return
}

// 禁用七牛存储中的对象
func (m *BucketManager) ChStatus(bucket, key string, forbidden bool) (err error) {
	ctx := auth.WithCredentials(context.Background(), m.Mac)
	reqHost, reqErr := m.RsReqHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	var status int
	if forbidden {
		status = 1
	} else {
		status = 0
	}
	reqURL := fmt.Sprintf("%s%s", reqHost, fmt.Sprintf("/chstatus/%s/status/%d", storage.EncodedEntry(bucket, key), status))
	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_FORM)
	err = m.Client.Call(ctx, nil, "POST", reqURL, headers)
	return

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

// GetBucketManager 返回一个BucketManager 指针
func GetBucketManagerWithConfig(cfg *storage.Config) *BucketManager {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		fmt.Fprintf(os.Stderr, "GetBucketManager: %v\n", gErr)
		os.Exit(1)
	}
	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	return NewBucketManager(mac, cfg)
}

func GetBucketManager() *BucketManager {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		fmt.Fprintf(os.Stderr, "GetBucketManager: %v\n", gErr)
		os.Exit(1)
	}
	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cfg := storage.Config{
		UpHost:        config.UpHost(),
		IoHost:        config.IoHost(),
		RsHost:        config.RsHost(),
		ApiHost:       config.ApiHost(),
		RsfHost:       config.RsfHost(),
		CentralRsHost: config.RsHost(),
	}
	storage.UcHost = config.UcHost()
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
