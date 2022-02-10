package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Downloader struct {
	files

	Url                 string // 文件下载的 url 【必填】
	Domain              string // 文件下载的 domain 【选填】
	Referer             string // 请求 header 中的 Referer 【选填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件【选填】
}

func (d *Downloader) Download() (err error) {
	defer func() {
		if err != nil && d.RemoveFileWhenError {
			e := os.Remove(d.tempFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove temp file error:%v", e)
			}

			e = os.Remove(d.toFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove file error:%v", e)
			}
		}
	}()

	err = d.downloadFile()
	if err != nil {
		return err
	}

	err = d.rename()
	return err
}

func (d *Downloader) downloadFile() error {
	//new request
	log.DebugF("download start:%s", d.Url)
	req, err := http.NewRequest("GET", d.Url, nil)
	log.DebugF("download   end:%s error:%v", d.Url, err)
	if err != nil {
		return errors.New("New request failed by url:" + d.Url + " error:" + err.Error())
	}
	// set host
	if len(d.Domain) > 0 {
		req.Host = d.Domain
	}
	// 设置断点续传
	if d.fromBytes > 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", d.fromBytes))
	}
	// 配置 referer
	if len(d.Referer) > 0 {
		req.Header.Add("Referer", d.Referer)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Download failed by url:" + d.Url + " error:" + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.New("Download failed by url:" + d.Url + " status code:" + strconv.Itoa(resp.StatusCode))
	}

	var tempFileHandle *os.File
	if d.fromBytes > 0 {
		tempFileHandle, err = os.OpenFile(d.tempFile, os.O_APPEND|os.O_WRONLY, 0655)
	} else {
		tempFileHandle, err = os.Create(d.tempFile)
	}
	if err != nil {
		return errors.New("Open local temp file error:" + d.tempFile + " error:" + err.Error())
	}
	defer tempFileHandle.Close()

	_, err = io.Copy(tempFileHandle, resp.Body)
	if err != nil {
		return errors.New("Download failed by url:" + d.Url + " because save to local error:" + err.Error())
	}

	return nil
}

func (d *Downloader) rename() error {
	err := os.Rename(d.tempFile, d.toFile)
	if err != nil {
		return errors.New("Rename temp file to final file error" + err.Error())
	}
	return nil
}

func utf82GBK(text string) (string, error) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	return gbkEncoder.String(text)
}
