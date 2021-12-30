package m3u8

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Slice struct {
	Bucket  string
	Key     string
	PutTime string
}

type SliceListApiInfo struct {
	Bucket string
	Key    string
}

func Slices(info SliceListApiInfo) ([]Slice, error) {
	dnLink, err := downloadLink(downloadLinkApiInfo(info))
	if err != nil {
		return nil, err
	}

	dnLink, err = object.PrivateUrl(object.PrivateUrlApiInfo{
		PublicUrl: dnLink,
		Deadline:  time.Now().Add(time.Second * 3600).Unix(),
	})
	if err != nil {
		return nil, err
	}

	//get m3u8 file content
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("new request for url %s error, %s", dnLink, reqErr)
	}

	m3u8Resp, m3u8Err := http.DefaultClient.Do(m3u8Req)
	if m3u8Err != nil {
		return nil, fmt.Errorf("open url %s error, %s", dnLink, m3u8Err)
	}

	defer m3u8Resp.Body.Close()
	if m3u8Resp.StatusCode != 200 {
		return nil, fmt.Errorf("download m3u8 file error, %s", m3u8Resp.Status)
	}

	m3u8Bytes, readErr := ioutil.ReadAll(m3u8Resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read m3u8 file content error, %s", readErr.Error())
	}

	//check content
	if !strings.HasPrefix(string(m3u8Bytes), "#EXTM3U") {
		return nil, errors.New("invalid m3u8 file")
	}

	slices := make([]Slice, 0)
	bReader := bufio.NewScanner(bytes.NewReader(m3u8Bytes))
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if !strings.HasPrefix(line, "#") {
			var sliceKey string
			if strings.HasPrefix(line, "http://") ||
				strings.HasPrefix(line, "https://") {
				uri, pErr := url.Parse(line)
				if pErr != nil {
					log.Warning("invalid url,", line)
					continue
				}
				sliceKey = strings.TrimPrefix(uri.Path, "/")
			} else {
				sliceKey = strings.TrimPrefix(line, "/")
			}
			slices = append(slices, Slice{Bucket: info.Bucket, Key: sliceKey})
		}
	}
	slices = append(slices, Slice{Bucket: info.Bucket, Key: info.Key})
	return slices, nil
}

type downloadLinkApiInfo SliceListApiInfo

func downloadLink(info downloadLinkApiInfo) (dnLink string, err error) {
	m, err := bucket.GetBucketManager()
	if err != nil {
		return "", err
	}

	_, sErr := m.Stat(info.Bucket, info.Key)
	if sErr != nil {
		err = fmt.Errorf("stat m3u8 file error, %s", sErr)
		return
	}

	bucketDomains, bErr := m.ListBucketDomains(info.Bucket)
	if bErr != nil {
		err = fmt.Errorf("get domain of bucket failed, %s", bErr.Error())
		return
	}

	if len(bucketDomains) == 0 {
		err = errors.New("no domain found for the bucket")
		return
	}

	var domain string
	for _, d := range bucketDomains {
		if strings.HasSuffix(d.Domain, "qiniudn.com") ||
			strings.HasSuffix(d.Domain, "clouddn.com") ||
			strings.HasSuffix(d.Domain, "qiniucdn.com") {
			domain = d.Domain
			break
		}
	}

	//get first
	if domain == "" && len(bucketDomains) > 0 {
		domain = bucketDomains[0].Domain
	}

	if domain == "" {
		return
	}

	dnLink = fmt.Sprintf("http://%s/%s", domain, info.Key)
	return
}
