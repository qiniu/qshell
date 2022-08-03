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
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
	"strings"
)

type ApiInfo struct {
	Bucket               string            `json:"bucket"`               // 文件所在 bucket 【必填】
	Key                  string            `json:"key"`                  // 文件被保存的 key 【必填】
	IsPublic             bool              `json:"-"`                    // 是否使用公有链接 【必填】
	HostProvider         host.Provider     `json:"-"`                    // 文件下载的 host, domain 可能为 ip, 需要搭配 host 使用 【选填】
	DestDir              string            `json:"-"`                    // 文件存储目标路径，目前是为了方便用户在批量下载时构建 ToFile 【此处选填】
	ToFile               string            `json:"to_file"`              // 文件保存的路径 【必填】
	Referer              string            `json:"referer"`              // 请求 header 中的 Referer 【选填】
	FileEncoding         string            `json:"-"`                    // 文件编码方式 【选填】
	ServerFilePutTime    int64             `json:"server_file_put_time"` // 文件修改时间 【选填】
	ServerFileSize       int64             `json:"server_file_size"`     // 文件大小，有值则会检测文件大小 【选填】
	ServerFileHash       string            `json:"server_file_hash"`     // 文件 hash，有值则会检测 hash 【选填】
	FromBytes            int64             `json:"-"`                    // 下载开始的位置，内部会缓存 【内部使用】
	ToBytes              int64             `json:"-"`                    // 下载的终止位置【内部使用】
	RemoveTempWhileError bool              `json:"-"`                    // 当遇到错误时删除临时文件 【选填】
	UseGetFileApi        bool              `json:"-"`                    // 是否使用 get file api(私有云会使用)【选填】
	BigFileEnableSlice   bool              `json:"-"`                    // 大文件允许切片下载，大于 300M 【选填】
	Progress             progress.Progress `json:"-"`                    // 下载进度回调【选填】
}

func (i *ApiInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", i.Bucket, i.Key, i.ToFile)
}

type ApiResult struct {
	FileModifyTime int64  `json:"file_modify_time"` // 下载后文件修改时间
	FileAbsPath    string `json:"file_abs_path"`    // 文件被保存的绝对路径
	IsUpdate       bool   `json:"is_update"`        // 是否为接续下载
	IsExist        bool   `json:"is_exist"`         // 是否为已存在
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

	// 以 '/' 结尾，不管大小是否为 0 ，均视为文件夹
	if strings.HasSuffix(info.Key, "/") {
		res.IsExist, _ = utils.ExistDir(f.toAbsFile)
		if !res.IsExist {
			err = utils.CreateDirIfNotExist(f.toAbsFile)
		}
		res.FileModifyTime, _ = utils.FileModify(f.toAbsFile)
		return res, err
	}

	// 文件存在则检查文件状态
	checkMode := object.MatchCheckModeFileSize
	if len(info.ServerFileHash) > 0 {
		checkMode = object.MatchCheckModeFileHash
	}

	// 读不到 status 按不存在该文件处理
	fileStatus, _ := os.Stat(f.toAbsFile)
	tempFileStatus, _ := os.Stat(f.tempFile)
	if tempFileStatus != nil && tempFileStatus.Size() > 0 {
		// 文件已下载了一部分，需要继续下载
		res.IsUpdate = true
	} else if fileStatus != nil {
		// 文件已下载，检测文件内容
		checkResult, mErr := object.Match(object.MatchApiInfo{
			Bucket:         info.Bucket,
			Key:            info.Key,
			LocalFile:      f.toAbsFile,
			CheckMode:      checkMode,
			ServerFileHash: info.ServerFileHash,
			ServerFileSize: info.ServerFileSize,
		})
		if mErr != nil {
			f.fromBytes = 0
			log.DebugF("check error before download:%v", mErr)
		}
		if checkResult != nil {
			res.IsExist = checkResult.Exist
		}

		if mErr == nil && checkResult.Match {
			// 文件已下载，并在文件匹配，不再下载
			if fileModifyTime, fErr := utils.FileModify(f.toAbsFile); fErr != nil {
				log.WarningF("Get file ModifyTime error:%v", fErr)
			} else {
				res.FileModifyTime = fileModifyTime
			}
			return
		}
	}

	// 下载
	err = download(f, info)
	if err != nil {
		return
	}

	res.FileModifyTime, err = utils.FileModify(f.toAbsFile)
	if err != nil {
		return
	}

	// 检查下载后的数据是否符合预期
	checkResult, mErr := object.Match(object.MatchApiInfo{
		Bucket:         info.Bucket,
		Key:            info.Key,
		LocalFile:      f.toAbsFile,
		CheckMode:      checkMode,
		ServerFileHash: info.ServerFileHash,
		ServerFileSize: info.ServerFileSize,
	})
	if mErr != nil || (checkResult != nil && !checkResult.Match) {
		return res, data.NewEmptyError().AppendDesc("check error after download").AppendError(mErr)
	}

	return res, nil
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
		return data.NewEmptyError().AppendDesc("Download create downloader error:" + err.Error())
	}

	var response *http.Response
	for times := 0; times < 6; times++ {
		if available, _ := info.HostProvider.Available(); !available {
			log.DebugF("Stop download [%s:%s] => %s, because no available host", info.Bucket, info.Key, info.ToFile)
			break
		}

		response, err = dl.Download(info)
		if err == nil && response != nil && response.StatusCode/100 == 2 {
			break
		}

		if response != nil {
			if (response.StatusCode > 399 && response.StatusCode < 500) ||
				response.StatusCode == 612 || response.StatusCode == 631 {
				log.DebugF("Stop download [%s:%s] => %s, because [%s]", info.Bucket, info.Key, info.ToFile, response.Status)
				break
			}
		}
	}

	if response != nil && response.Body != nil {
		if info.Progress != nil {
			info.Progress.SetFileSize(response.ContentLength + info.FromBytes)
			info.Progress.SendSize(info.FromBytes)
			info.Progress.Start()
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
	if response.Body == nil {
		return data.NewEmptyError().AppendDesc(" Download error: response body empty")
	}

	var fErr error
	var tempFileHandle *os.File
	isExist, _ := utils.ExistFile(fInfo.tempFile)
	if isExist {
		tempFileHandle, fErr = os.OpenFile(fInfo.tempFile, os.O_APPEND|os.O_WRONLY, 0655)
		log.DebugF("download [%s:%s] => %s from:%d", info.Bucket, info.Key, info.ToFile, info.FromBytes)
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
	err := os.Rename(fInfo.tempFile, fInfo.toAbsFile)
	if err != nil {
		return data.NewEmptyError().AppendDescF(" Rename temp file to final file error:%v", err.Error())
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
		if info.BigFileEnableSlice && info.ServerFileSize > 40*utils.MB {
			return &sliceDownloader{}, nil
		} else {
			return &getDownloader{useHttps: userHttps}, nil
		}
	}
}

func utf82GBK(text string) (string, *data.CodeError) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	d, err := gbkEncoder.String(text)
	return d, data.ConvertError(err)
}
