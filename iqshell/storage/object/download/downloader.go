package download

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/host"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
	"strconv"
)

type ApiInfo struct {
	Bucket               string            // 文件所在 bucket 【必填】
	Key                  string            // 文件被保存的 key 【必填】
	IsPublic             bool              // 是否使用共有链接 【必填】
	HostProvider         host.Provider     // 文件下载的 host, domain 可能为 ip, 需要搭配 host 使用 【选填】
	ToFile               string            // 文件保存的路径 【必填】
	Referer              string            // 请求 header 中的 Referer 【选填】
	FileEncoding         string            // 文件编码方式 【选填】
	FileModifyTime       int64             // 文件修改时间 【选填】
	FileSize             int64             // 文件大小，有值则会检测文件大小 【选填】
	FileHash             string            // 文件 hash，有值则会检测 hash 【选填】
	FromBytes            int64             // 下载开始的位置，内部会缓存 【内部使用】
	RemoveTempWhileError bool              // 当遇到错误时删除临时文件 【选填】
	UseGetFileApi        bool              // 是否使用 get file api(私有云会使用)【选填】
	Progress             progress.Progress // 下载进度回调【选填】
}

func (i *ApiInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", i.Bucket, i.Key, i.ToFile)
}

type ApiResult struct {
	FileModifyTime int64  // 下载后文件修改时间
	FileAbsPath    string // 文件被保存的绝对路径
	IsUpdate       bool   // 是否为接续下载
	IsExist        bool   // 是否为已存在
}

var _ flow.Result = (*ApiResult)(nil)

func (a *ApiResult) IsValid() bool {
	return len(a.FileAbsPath) > 0 && a.FileModifyTime > 0
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info *ApiInfo) (res *ApiResult, err *data.CodeError) {
	if len(info.ToFile) == 0 {
		err = data.NewEmptyError().AppendDesc("the filename saved after downloading is empty")
		return
	}

	f, err := createDownloadFiles(info.ToFile, info.FileEncoding)
	if err != nil {
		return
	}

	res = &ApiResult{
		FileAbsPath: f.toAbsFile,
	}

	// 文件存在则检查文件状态
	fileStatus, sErr := os.Stat(f.toAbsFile)
	tempFileStatus, tempErr := os.Stat(f.tempFile)
	if sErr == nil || os.IsExist(err) || tempErr == nil || os.IsExist(tempErr) {
		if tempFileStatus != nil && tempFileStatus.Size() > 0 {
			// 文件是否已下载了一部分，需要继续下载
			res.IsUpdate = true
		}

		if fileStatus != nil {
			// 文件已下载，检测文件内容
			if checkErr := (&FileChecker{
				File:                f.toAbsFile,
				Bucket:              info.Bucket,
				Key:                 info.Key,
				FileHash:            info.FileHash,
				FileSize:            info.FileSize,
				RemoveFileWhenError: true,
			}).IsFileMatch(); checkErr == nil {
				// 文件是否已下载完成，如果完成跳过下载阶段，直接验证
				res.IsExist = true
				return
			} else {
				f.fromBytes = 0
				log.DebugF("check error before download:%v", checkErr)
			}
		}
	}

	// 下载
	err = download(f, info)
	if err != nil {
		return
	}

	info.FileModifyTime, err = utils.FileModify(f.toAbsFile)
	if err != nil {
		return
	}

	// 检查下载后的数据是否符合预期
	if checkErr := (&FileChecker{
		File:                f.toAbsFile,
		Bucket:              info.Bucket,
		Key:                 info.Key,
		FileHash:            info.FileHash,
		FileSize:            info.FileSize,
		RemoveFileWhenError: true,
	}).IsFileMatch(); checkErr != nil {
		log.DebugF("check error after download:%v", checkErr)
	}
	return
}

func download(fInfo *fileInfo, info *ApiInfo) (err *data.CodeError) {
	defer func() {
		if info.RemoveTempWhileError && err != nil {
			e := os.Remove(fInfo.tempFile)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download: remove temp file error:%v", e)
			} else {
				log.DebugF("download: remove temp file success:%s", fInfo.tempFile)
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

func downloadFile(fInfo *fileInfo, info *ApiInfo) *data.CodeError {
	dl, err := createDownloader(info)
	if err != nil {
		return data.NewEmptyError().AppendDesc(" Download create downloader error:" + err.Error())
	}

	var response *http.Response
	for i := 0; i < 6; i++ {
		response, err = dl.Download(info)
		if err == nil || !utils.IsHostUnavailableError(err) {
			break
		}
	}

	if response != nil && response.Body != nil {
		if info.Progress != nil {
			size := response.Header.Get("Content-Length")
			if sizeInt, err := strconv.ParseInt(size, 10, 64); err == nil {
				info.Progress.SetFileSize(sizeInt + info.FromBytes)
				info.Progress.SendSize(info.FromBytes)
				info.Progress.Start()
			}
		}
		defer response.Body.Close()
	}

	if err != nil {
		return data.NewEmptyError().AppendDesc(" Download error:" + err.Error())
	}
	if response == nil {
		return data.NewEmptyError().AppendDesc(" Download error: response empty")
	}
	if response.StatusCode/100 != 2 {
		return data.NewEmptyError().AppendDescF(" Download error: %v", response)
	}
	defer response.Body.Close()

	var fErr error
	var tempFileHandle *os.File
	if info.FromBytes > 0 {
		tempFileHandle, fErr = os.OpenFile(fInfo.tempFile, os.O_APPEND|os.O_WRONLY, 0655)
		log.InfoF("download [%s:%s] => %s from:%d", info.Bucket, info.Key, info.ToFile, info.FromBytes)
	} else {
		tempFileHandle, fErr = os.Create(fInfo.tempFile)
	}
	if fErr != nil {
		return data.NewEmptyError().AppendDesc(" Open local temp file error:" + fInfo.tempFile + " error:" + fErr.Error())
	}
	defer tempFileHandle.Close()

	if info.Progress != nil {
		_, fErr = io.Copy(tempFileHandle, io.TeeReader(response.Body, info.Progress))
		if fErr == nil {
			info.Progress.End()
		}
	} else {
		_, fErr = io.Copy(tempFileHandle, response.Body)
	}
	if fErr != nil {
		return data.NewEmptyError().AppendDescF(" Download error:%v", fErr)
	}

	return nil
}

func renameTempFile(fInfo *fileInfo, info *ApiInfo) *data.CodeError {
	err := os.Rename(fInfo.tempFile, fInfo.toFile)
	if err != nil {
		return data.NewEmptyError().AppendDesc(" Rename temp file to final file error" + err.Error())
	}
	return nil
}

type downloader interface {
	Download(info *ApiInfo) (response *http.Response, err *data.CodeError)
}

func createDownloader(info *ApiInfo) (downloader, *data.CodeError) {
	userHttps := workspace.GetConfig().IsUseHttps()
	if info.UseGetFileApi {
		mac, err := workspace.GetMac()
		if err != nil {
			return nil, data.NewEmptyError().AppendDescF("download get mac error:%v", mac)
		}
		return &getFileApiDownloader{
			useHttps: userHttps,
			mac:      mac,
		}, nil
	} else {
		return &getDownloader{useHttps: userHttps}, nil
	}
}

func utf82GBK(text string) (string, *data.CodeError) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	d, err := gbkEncoder.String(text)
	return d, data.ConvertError(err)
}
