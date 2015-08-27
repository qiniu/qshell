package qshell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BucketDomain []string

func M3u8FileList(mac *digest.Mac, bucket string, m3u8Key string, isPrivate bool) (slicesToDelete []rs.EntryPath, err error) {
	client := rs.New(mac)
	//check m3u8 file exists
	_, sErr := client.Stat(nil, bucket, m3u8Key)
	if sErr != nil {
		err = errors.New(fmt.Sprintf("stat m3u8 file error, %s", sErr.Error()))
		return
	}
	//get domain list of bucket
	bucketDomainUrl := "http://api.qiniu.com/v6/domain/list"
	bucketDomainData := map[string][]string{
		"tbl": []string{bucket},
	}
	bucketDomains := BucketDomain{}
	bErr := client.Conn.CallWithForm(nil, &bucketDomains, bucketDomainUrl, bucketDomainData)
	if bErr != nil {
		err = errors.New(fmt.Sprintf("get domain of bucket failed due to, %s", bErr.Error()))
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
		err = errors.New("no valid domain found for the bucket")
		return
	}
	//create downoad link
	dnLink := fmt.Sprintf("http://%s/%s", domain, m3u8Key)
	if isPrivate {
		dnLink = PrivateUrl(mac, dnLink, time.Now().Add(time.Second*3600).Unix())
	}
	//get m3u8 file content
	m3u8Resp, m3u8Err := http.Get(dnLink)
	if m3u8Err != nil {
		err = errors.New(fmt.Sprintf("open url %s error due to, %s", dnLink, m3u8Err))
		return
	}
	defer m3u8Resp.Body.Close()
	if m3u8Resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("download file error due to, %s", m3u8Resp.Status))
		return
	}
	m3u8Bytes, readErr := ioutil.ReadAll(m3u8Resp.Body)
	if readErr != nil {
		err = errors.New(fmt.Sprintf("read m3u8 file content error due to, %s", readErr.Error()))
		return
	}
	//check content
	if !strings.HasPrefix(string(m3u8Bytes), "#EXTM3U") {
		err = errors.New("invalid m3u8 file")
		return
	}
	slicesToDelete = make([]rs.EntryPath, 0)
	bReader := bufio.NewScanner(bytes.NewReader(m3u8Bytes))
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if !strings.HasPrefix(line, "#") {
			var sliceKey string
			if strings.HasPrefix(line, "http://") ||
				strings.HasPrefix(line, "https://") {
				uri, pErr := url.Parse(line)
				if pErr != nil {
					log.Error(fmt.Sprintf("invalid url, %s", line))
					continue
				}
				sliceKey = strings.TrimPrefix(uri.Path, "/")
			} else {
				sliceKey = strings.TrimPrefix(line, "/")
			}
			//append to delete list
			slicesToDelete = append(slicesToDelete, rs.EntryPath{bucket, sliceKey})
		}
	}
	slicesToDelete = append(slicesToDelete, rs.EntryPath{bucket, m3u8Key})
	return
}
