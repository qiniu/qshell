package qshell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	fio "qiniu/api.v6/io"
	"qiniu/api.v6/rs"
	"qiniu/log"
	"qiniu/rpc"
	"strings"
	"time"
)

type BucketDomain []string

func M3u8FileList(mac *digest.Mac, bucket string, m3u8Key string) (slicesToDelete []rs.EntryPath, err error) {
	client := rs.NewMac(mac)
	//check m3u8 file exists
	_, sErr := client.Stat(nil, bucket, m3u8Key)
	if sErr != nil {
		if v, ok := sErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("stat m3u8 file error, %s", v.Err)
		} else {
			err = fmt.Errorf("stat m3u8 file error, %s", sErr)
		}
		return
	}
	//get domain list of bucket
	bucketDomainUrl := fmt.Sprintf("%s/v6/domain/list", conf.API_HOST)
	bucketDomainData := map[string][]string{
		"tbl": []string{bucket},
	}
	bucketDomains := BucketDomain{}
	bErr := client.Conn.CallWithForm(nil, &bucketDomains, bucketDomainUrl, bucketDomainData)
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
		err = errors.New("no valid domain found for the bucket")
		return
	}
	//create downoad link
	dnLink := fmt.Sprintf("http://%s/%s", domain, m3u8Key)
	dnLink = PrivateUrl(mac, dnLink, time.Now().Add(time.Second*3600).Unix())
	//get m3u8 file content
	dnLink = strings.Replace(dnLink, fmt.Sprintf("http://%s", domain), conf.IO_HOST, -1)
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		err = fmt.Errorf("new request for url %s error, %s", dnLink, reqErr)
		return
	}
	m3u8Req.Host = domain
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
					log.Errorf("invalid url, %s", line)
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

//replace and upload
func M3u8ReplaceDomain(mac *digest.Mac, bucket string, m3u8Key string, newDomain string) (err error) {
	client := rs.NewMac(mac)
	//check m3u8 file exists
	_, sErr := client.Stat(nil, bucket, m3u8Key)
	if sErr != nil {
		if v, ok := sErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("stat m3u8 file error, %s", v.Err)
		} else {
			err = fmt.Errorf("stat m3u8 file error, %s", sErr)
		}
		return
	}
	//get domain list of bucket
	bucketDomainUrl := fmt.Sprintf("%s/v6/domain/list", conf.API_HOST)
	bucketDomainData := map[string][]string{
		"tbl": []string{bucket},
	}
	bucketDomains := BucketDomain{}
	bErr := client.Conn.CallWithForm(nil, &bucketDomains, bucketDomainUrl, bucketDomainData)
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
		err = errors.New("no valid domain found for the bucket")
		return
	}
	//create downoad link
	dnLink := fmt.Sprintf("http://%s/%s", domain, m3u8Key)
	dnLink = PrivateUrl(mac, dnLink, time.Now().Add(time.Second*3600).Unix())
	//get m3u8 file content
	dnLink = strings.Replace(dnLink, fmt.Sprintf("http://%s", domain), conf.IO_HOST, -1)
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		err = fmt.Errorf("new request for url %s error, %s", dnLink, reqErr)
		return
	}
	m3u8Req.Host = domain
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
					log.Errorf("invalid url, %s", line)
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

	putPolicy := rs.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, m3u8Key),
	}
	upToken := putPolicy.Token(mac)

	putClient := rpc.NewClient("")
	putErr := fio.Put2(putClient, nil, nil, upToken, m3u8Key, bytes.NewReader(newM3u8Data),
		int64(len(newM3u8Data)), nil)
	if putErr != nil {
		err = putErr
		return
	}
	return
}
