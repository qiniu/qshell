package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
)

type ApiInfo struct {
	IsPublic       bool   // 是否使用共有链接 【必填】
	Domain         string // 文件下载的 domain 【必填】
	ToFile         string // 文件保存的路径 【必填】
	StatusDBPath   string // 下载状态缓存的 db 路径 【选填】
	Referer        string // 请求 header 中的 Referer 【选填】
	FileEncoding   string // 文件编码方式 【选填】
	Bucket         string // 文件所在 bucket，用于验证 hash 【选填】
	Key            string // 文件被保存的 key，用于验证 hash 【选填】
	FileModifyTime int64  // 文件修改时间 【选填】
	FileSize       int64  // 文件大小，有值则会检测文件大小 【选填】
	FileHash       string // 文件 hash，有值则会检测 hash 【选填】
	FromBytes      int64  // 下载开始的位置，内部会缓存 【选填】
	UserGetFileApi bool   // 是否使用 get file api(私有云会使用)
}

type ApiResult struct {
	FileAbsPath string // 文件被保存的绝对路径
	IsUpdate    bool   // 是否为接续下载
	IsExist     bool   // 是否为已存在
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info ApiInfo) (res ApiResult, err error) {
	if len(info.ToFile) == 0 {
		err = errors.New("the filename saved after downloading is empty")
		return
	}

	f, err := createDownloadFiles(info.ToFile, info.FileEncoding)
	if err != nil {
		return
	}

	// 检查文件是否已存在，如果存在是否符合预期
	dbChecker := &dbHandler{
		DBFilePath:           info.StatusDBPath,
		FilePath:             f.toAbsFile,
		FileHash:             info.FileHash,
		FileSize:             info.FileSize,
		FileServerUpdateTime: info.FileModifyTime,
	}
	err = dbChecker.init()
	if err != nil {
		return
	}

	shouldDownload := true

	// 文件存在则检查文件状态
	fileStatus, err := os.Stat(f.toAbsFile)
	tempFileStatus, tempErr := os.Stat(f.tempFile)
	if err == nil || os.IsExist(err) || tempErr == nil || os.IsExist(tempErr) {
		// 中间文件 和 最终文件 中任意一个存在
		if cErr := dbChecker.checkInfoOfDB(); cErr != nil {
			// 检查服务端文件是否变更，如果变更则清除
			log.WarningF("Local file `%s` exist for key `%s`, but not match:%v", f.toAbsFile, info.Key, cErr)
			if e := f.clean(); e != nil {
				log.WarningF("Local file `%s` exist for key `%s`, clean error:%v", f.toAbsFile, info.Key, e)
			}
			if sErr := dbChecker.saveInfoToDB(); sErr != nil {
				log.WarningF("Local file `%s` exist for key `%s`, save info to db error:%v", f.toAbsFile, info.Key, sErr)
			}
		}
		if tempFileStatus != nil && tempFileStatus.Size() > 0 {
			// 文件是否已下载了一部分，需要继续下载
			res.IsUpdate = true
		}
		if fileStatus != nil && fileStatus.Size() == info.FileSize {
			// 文件是否已下载完成，如果完成跳过下载阶段，直接验证
			res.IsExist = true
			shouldDownload = false
		}
	}

	res.FileAbsPath = f.toAbsFile

	// 下载
	if shouldDownload {
		err = download(f, info)
		if err != nil {
			return
		}
		err = dbChecker.saveInfoToDB()
		if err != nil {
			err = fmt.Errorf("download info save to db error:%v key:%s localFile:%s", err, f.toAbsFile, info.Key)
			return
		}
	}

	// 检查下载后的数据是否符合预期
	err = (&LocalFileInfo{
		File:                f.toAbsFile,
		Bucket:              info.Bucket,
		Key:                 info.Key,
		FileHash:            info.FileHash,
		FileSize:            info.FileSize,
		RemoveFileWhenError: true,
	}).CheckDownloadFile()
	if err != nil {
		return
	}

	return
}

func download(fInfo fileInfo, info ApiInfo) (err error) {
	defer func() {
		if err != nil {
			e := os.Remove(fInfo.tempFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove temp file error:%v", e)
			}

			e = os.Remove(fInfo.toFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove file error:%v", e)
			}
		}
	}()

	info.FromBytes = fInfo.fromBytes
	err = downloadFile(fInfo, info)
	if err != nil {
		return err
	}

	err = renameTempFile(fInfo, info)
	return err
}

func downloadFile(fInfo fileInfo, info ApiInfo) error {
	dl, err := createDownloader(info)
	if err != nil {
		return errors.New(" Download create downloader error:" + err.Error())
	}

	response, err := dl.Download(info)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return errors.New(" Download error:" + err.Error())
	}
	if response == nil {
		return errors.New(" Download error: response empty")
	}
	if response.StatusCode/100 != 2 {
		return fmt.Errorf(" Download error: %v", response)
	}

	var tempFileHandle *os.File
	if info.FromBytes > 0 {
		tempFileHandle, err = os.OpenFile(fInfo.tempFile, os.O_APPEND|os.O_WRONLY, 0655)
	} else {
		tempFileHandle, err = os.Create(fInfo.tempFile)
	}
	if err != nil {
		return errors.New(" Open local temp file error:" + fInfo.tempFile + " error:" + err.Error())
	}
	defer tempFileHandle.Close()

	_, err = io.Copy(tempFileHandle, response.Body)
	if err != nil {
		return fmt.Errorf(" Download error:%v", err)
	}

	return nil
}

func renameTempFile(fInfo fileInfo, info ApiInfo) error {
	err := os.Rename(fInfo.tempFile, fInfo.toFile)
	if err != nil {
		return errors.New(" Rename temp file to final file error" + err.Error())
	}
	return nil
}

type downloader interface {
	Download(info ApiInfo) (response *http.Response, err error)
}

func createDownloader(info ApiInfo) (downloader, error) {
	userHttps := workspace.GetConfig().IsUseHttps()
	if info.UserGetFileApi {
		mac, err := workspace.GetMac()
		if err != nil {
			return nil, fmt.Errorf("download get mac error:%v", mac)
		}
		return &getFileApiDownloader{
			useHttps:   userHttps,
			mac:        mac,
		}, nil
	} else {
		return &getDownloader{useHttps: userHttps}, nil
	}
}

func utf82GBK(text string) (string, error) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	return gbkEncoder.String(text)
}