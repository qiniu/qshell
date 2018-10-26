package qshell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (m *BucketManager) M3u8FileList(bucket string, m3u8Key string) (slicesToDelete []EntryPath, err error) {
	dnLink, err := m.DownloadLink(bucket, m3u8Key)
	if err != nil {
		return
	}
	dnLink, err = m.PrivateUrl(dnLink, time.Now().Add(time.Second*3600).Unix())
	if err != nil {
		return
	}
	//get m3u8 file content
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		err = fmt.Errorf("new request for url %s error, %s", dnLink, reqErr)
		return
	}
	m3u8Resp, m3u8Err := http.DefaultClient.Do(m3u8Req)
	if m3u8Err != nil {
		err = fmt.Errorf("open url %s error, %s", dnLink, m3u8Err)
		return
	}
	defer m3u8Resp.Body.Close()
	if m3u8Resp.StatusCode != 200 {
		err = fmt.Errorf("download m3u8 file error, %s", m3u8Resp.Status)
		return
	}
	m3u8Bytes, readErr := ioutil.ReadAll(m3u8Resp.Body)
	if readErr != nil {
		err = fmt.Errorf("read m3u8 file content error, %s", readErr.Error())
		return
	}
	//check content
	if !strings.HasPrefix(string(m3u8Bytes), "#EXTM3U") {
		err = errors.New("invalid m3u8 file")
		return
	}
	slicesToDelete = make([]EntryPath, 0)
	bReader := bufio.NewScanner(bytes.NewReader(m3u8Bytes))
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if !strings.HasPrefix(line, "#") {
			var sliceKey string
			if strings.HasPrefix(line, "http://") ||
				strings.HasPrefix(line, "https://") {
				uri, pErr := url.Parse(line)
				if pErr != nil {
					fmt.Fprintln(os.Stderr, "invalid url,", line)
					continue
				}
				sliceKey = strings.TrimPrefix(uri.Path, "/")
			} else {
				sliceKey = strings.TrimPrefix(line, "/")
			}
			//append to delete list
			slicesToDelete = append(slicesToDelete, EntryPath{bucket, sliceKey})
		}
	}
	slicesToDelete = append(slicesToDelete, EntryPath{bucket, m3u8Key})
	return
}

func (m *BucketManager) DownloadLink(bucket, key string) (dnLink string, err error) {

	_, sErr := m.Stat(bucket, key)
	if sErr != nil {
		err = fmt.Errorf("stat m3u8 file error, %s", sErr)
		return
	}
	bucketDomains, bErr := m.DomainsOfBucket(bucket)
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
		if strings.HasSuffix(d, "qiniudn.com") ||
			strings.HasSuffix(d, "clouddn.com") ||
			strings.HasSuffix(d, "qiniucdn.com") {
			domain = d
			break
		}
	}
	//get first
	if domain == "" {
		domain = bucketDomains[0]
	}

	if domain == "" {
		return
	}
	dnLink = fmt.Sprintf("http://%s/%s", domain, key)
	return
}

//replace and upload
func (m *BucketManager) M3u8ReplaceDomain(bucket string, m3u8Key string, newDomain string) (err error) {
	dnLink, err := m.DownloadLink(bucket, m3u8Key)

	//create downoad link
	dnLink, err = m.PrivateUrl(dnLink, time.Now().Add(time.Second*3600).Unix())
	//get m3u8 file content
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		err = fmt.Errorf("new request for url %s error, %s", dnLink, reqErr)
		return
	}
	m3u8Resp, m3u8Err := http.DefaultClient.Do(m3u8Req)
	if m3u8Err != nil {
		err = fmt.Errorf("open url %s error, %s", dnLink, m3u8Err)
		return
	}
	defer m3u8Resp.Body.Close()
	if m3u8Resp.StatusCode != 200 {
		err = fmt.Errorf("download m3u8 file error, %s", m3u8Resp.Status)
		return
	}
	m3u8Bytes, readErr := ioutil.ReadAll(m3u8Resp.Body)
	if readErr != nil {
		err = fmt.Errorf("read m3u8 file content error, %s", readErr.Error())
		return
	}
	//check content
	if !strings.HasPrefix(string(m3u8Bytes), "#EXTM3U") {
		err = errors.New("invalid m3u8 file")
		return
	}

	newM3u8Lines := make([]string, 0, 200)
	var newLine string
	bReader := bufio.NewScanner(bytes.NewReader(m3u8Bytes))
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http://") ||
				strings.HasPrefix(line, "https://") {
				uri, pErr := url.Parse(line)
				if pErr != nil {
					fmt.Println("invalid url,", line)
					continue
				}

				if newDomain != "" {
					newLine = fmt.Sprintf("%s%s", newDomain, uri.Path)
				} else {
					newLine = uri.Path
				}
			} else {
				if newDomain != "" {
					newLine = fmt.Sprintf("%s%s", newDomain, line)
				} else {
					newLine = line
				}
			}
		} else {
			newLine = line
		}

		newM3u8Lines = append(newM3u8Lines, newLine)
	}

	//join and upload
	newM3u8Data := []byte(strings.Join(newM3u8Lines, "\n"))

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, m3u8Key),
	}
	upToken := putPolicy.UploadToken(m.GetMac())

	uploader := storage.NewFormUploader(nil)
	putRet := new(storage.PutRet)
	putErr := uploader.Put(nil, putRet, upToken, m3u8Key, bytes.NewReader(newM3u8Data), int64(len(newM3u8Data)), nil)

	if putErr != nil {
		err = putErr
		return
	}
	return
}
