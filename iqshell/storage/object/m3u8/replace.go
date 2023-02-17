package m3u8

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// replace and upload
type ReplaceDomainApiInfo struct {
	Bucket              string
	Key                 string
	NewDomain           string
	RemoveSparePreSlash bool
}

func ReplaceDomain(info ReplaceDomainApiInfo) *data.CodeError {
	dnLink, err := downloadLink(downloadLinkApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
	})

	//create download link
	if urlResult, e := download.PublicUrlToPrivate(download.PublicUrlToPrivateApiInfo{
		PublicUrl: dnLink,
		Deadline:  time.Now().Add(time.Second * 3600).Unix(),
	}); e != nil || urlResult == nil {
		return data.NewEmptyError().AppendDesc("public url to private").AppendError(err)
	} else {
		dnLink = urlResult.Url
	}

	//get m3u8 file content
	m3u8Req, reqErr := http.NewRequest("GET", dnLink, nil)
	if reqErr != nil {
		return data.NewEmptyError().AppendDescF("new request for url %s error, %s", dnLink, reqErr)
	}

	m3u8Resp, m3u8Err := http.DefaultClient.Do(m3u8Req)
	if m3u8Err != nil {
		return data.NewEmptyError().AppendDescF("open url %s error, %s", dnLink, m3u8Err)
	}
	defer m3u8Resp.Body.Close()

	if m3u8Resp.StatusCode != 200 {
		return data.NewEmptyError().AppendDescF("download m3u8 file error, %s", m3u8Resp.Status)
	}

	m3u8Bytes, readErr := ioutil.ReadAll(m3u8Resp.Body)
	if readErr != nil {
		return data.NewEmptyError().AppendDescF("read m3u8 file content error, %s", readErr.Error())
	}

	//check content
	if !strings.HasPrefix(string(m3u8Bytes), "#EXTM3U") {
		return data.NewEmptyError().AppendDesc("invalid m3u8 file")
	}

	newM3u8Lines := make([]string, 0, 200)
	bReader := bufio.NewScanner(bytes.NewReader(m3u8Bytes))
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		newLine := replaceTsNewDomain(line, info.NewDomain, info.RemoveSparePreSlash)
		log.DebugF("NewLine:%s", newLine)
		newM3u8Lines = append(newM3u8Lines, newLine)
	}

	//join and upload
	newM3u8Data := []byte(strings.Join(newM3u8Lines, "\n"))
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", info.Bucket, info.Key),
	}

	mac, err := workspace.GetMac()
	if err != nil {
		return err
	}

	upToken := putPolicy.UploadToken(mac)

	uploader := storage.NewFormUploader(nil)
	putRet := new(storage.PutRet)
	putExtra := storage.PutExtra{}
	putErr := uploader.Put(workspace.GetContext(), putRet, upToken, info.Key, bytes.NewReader(newM3u8Data), int64(len(newM3u8Data)), &putExtra)

	if putErr != nil {
		return data.ConvertError(putErr)
	}

	return nil
}

func replaceTsNewDomain(line string, newDomain string, removeSparePreSlash bool) (newLine string) {
	if strings.HasPrefix(line, "#") {
		return line
	}
	if strings.HasPrefix(line, "http://") ||
		strings.HasPrefix(line, "https://") {
		uri, pErr := url.Parse(line)
		if pErr != nil {
			fmt.Println("invalid url,", line)
			return line
		}
		line = uri.Path
	}
	if removeSparePreSlash {
		if strings.HasSuffix(newDomain, "/") || strings.HasPrefix(line, "/") {
			return strings.TrimRight(newDomain, "/") + "/" + strings.TrimLeft(line, "/")
		}
	}
	return newDomain + line
}
