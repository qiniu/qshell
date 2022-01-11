package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Downloader struct {
	Url                 string // 文件下载的 url 【必填】
	ToFile              string // 文件保存的路径 【必填】
	Domain              string // 文件下载的 domain 【选填】
	Referer             string // 请求 header 中的 Referer 【选填】
	FileEncoding        string // 文件编码方式 【选填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件【选填】

	fileDir   string // 保存文件的路径，从 ToFile 解析
	tempFile  string // 临时保存的文件路径 ToFile + .tmp
	fromBytes int64  // 下载开始位置，检查本地 tempFile 文件，读取已下载文件长度
}

func (d *Downloader) Download() (err error) {
	defer func() {
		if err != nil && d.RemoveFileWhenError {
			e := os.Remove(d.tempFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove temp file error:%v", e)
			}

			e = os.Remove(d.ToFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove file error:%v", e)
			}
		}
	}()

	err = d.check()
	if err != nil {
		return err
	}

	err = d.prepare()
	if err != nil {
		return err
	}

	err = d.downloadFile()
	if err != nil {
		return err
	}

	err = d.rename()
	return err
}

func (d *Downloader) check() error {
	if len(d.Url) == 0 {
		return errors.New("download url can't be empty")
	}
	if len(d.ToFile) == 0 {
		return errors.New("the filename saved after downloading is empty")
	}
	return nil
}

func (d *Downloader) prepare() (err error) {
	// 文件路径
	d.ToFile, err = filepath.Abs(d.ToFile)
	if err != nil {
		err = errors.New("get save file abs path error:" + err.Error())
		return
	}

	if strings.ToLower(d.FileEncoding) == "gbk" {
		d.ToFile, err = utf82GBK(d.ToFile)
		if err != nil {
			err = errors.New("gbk file path:" + d.ToFile + " error:" + err.Error())
			return
		}
	}

	d.fileDir = filepath.Dir(d.ToFile)
	d.tempFile = fmt.Sprintf("%s.tmp", d.ToFile)

	err = os.MkdirAll(d.fileDir, 0775)
	if err != nil {
		return errors.New("MkdirAll failed for " + d.fileDir + " error:" + err.Error())
	}

	tempFileStatus, err := os.Stat(d.tempFile)
	if err != nil && os.IsNotExist(err) {
		d.fromBytes = 0
		return nil
	}

	if tempFileStatus != nil && !tempFileStatus.IsDir() {
		d.fromBytes = tempFileStatus.Size()
	}

	return nil
}

func (d *Downloader) downloadFile() error {
	//new request
	req, err := http.NewRequest("GET", d.Url, nil)
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
	err := os.Rename(d.tempFile, d.ToFile)
	if err != nil {
		return errors.New("Rename temp file to final file error" + err.Error())
	}
	return nil
}

func utf82GBK(text string) (string, error) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	return gbkEncoder.String(text)
}
