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
	"sync"
)

type DownloadActionInfo struct {
	Bucket                 string            `json:"bucket"`               // 文件所在 bucket 【必填】
	Key                    string            `json:"key"`                  // 文件被保存的 key 【必填】
	IsPublic               bool              `json:"-"`                    // 是否使用公有链接 【必填】
	HostProvider           host.Provider     `json:"-"`                    // 文件下载的 host, domain 可能为 ip, 需要搭配 host 使用 【选填】
	DestDir                string            `json:"-"`                    // 文件存储目标路径，目前是为了方便用户在批量下载时构建 ToFile 【此处选填】
	ToFile                 string            `json:"to_file"`              // 文件保存的路径 【必填】
	Referer                string            `json:"referer"`              // 请求 header 中的 Referer 【选填】
	FileEncoding           string            `json:"-"`                    // 文件编码方式 【选填】
	ServerFilePutTime      int64             `json:"server_file_put_time"` // 文件修改时间 【选填】
	ServerFileSize         int64             `json:"server_file_size"`     // 文件大小，有值则会检测文件大小 【选填】
	ServerFileHash         string            `json:"server_file_hash"`     // 文件 hash，有值则会检测 hash 【选填】
	DownloadFileSize       int64             `json:"download_file_size"`   // 下载的文件大小，下载整个文件时，等于 ServerFileSize；切片下载则为切片大小；有值则会检测文件大小【选填】
	CheckSize              bool              `json:"-"`                    // 是否检测文件大小 【选填】
	CheckHash              bool              `json:"-"`                    // 是否检测文件 hash 【选填】
	FromBytes              int64             `json:"-"`                    // 下载开始的位置，内部会缓存 【内部使用】
	ToBytes                int64             `json:"-"`                    // 下载的终止位置【内部使用】
	RemoveTempWhileError   bool              `json:"-"`                    // 当遇到错误时删除临时文件 【选填】
	UseGetFileApi          bool              `json:"-"`                    // 是否使用 get file api(私有云会使用)【选填】
	EnableSlice            bool              `json:"-"`                    // 大文件允许切片下载 【选填】
	SliceFileSizeThreshold int64             `json:"-"`                    // 允许切片下载，切片下载出发的文件大小阈值 【选填】
	SliceSize              int64             `json:"-"`                    // 允许切片下载，切片的大小 【选填】
	SliceConcurrentCount   int               `json:"-"`                    // 允许切片下载，并发下载切片的个数 【选填】
	Progress               progress.Progress `json:"-"`                    // 下载进度回调【选填】
}

func (i *DownloadActionInfo) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", i.Bucket, i.Key, i.ToFile)
}

type DownloadActionResult struct {
	FileModifyTime int64  `json:"file_modify_time"` // 下载后文件修改时间
	FileAbsPath    string `json:"file_abs_path"`    // 文件被保存的绝对路径
	IsUpdate       bool   `json:"is_update"`        // 是否为接续下载
	IsExist        bool   `json:"is_exist"`         // 是否为已存在
}

var _ flow.Result = (*DownloadActionResult)(nil)

func (a *DownloadActionResult) IsValid() bool {
	return len(a.FileAbsPath) > 0 && a.FileModifyTime > 0
}

// Download 下载一个文件，从 Url 下载保存至 ToFile
func Download(info *DownloadActionInfo) (res *DownloadActionResult, err *data.CodeError) {
	if len(info.ToFile) == 0 {
		err = data.NewEmptyError().AppendDesc("the filename saved after downloading is empty")
		return
	}

	f, err := createDownloadFiles(info.ToFile, info.FileEncoding)
	if err != nil {
		return
	}

	res = &DownloadActionResult{
		FileAbsPath: f.toAbsFile,
	}

	// 以 '/' 结尾，不管大小是否为 0 ，均视为文件夹
	if strings.HasSuffix(info.Key, "/") {
		if info.ServerFileSize > 0 {
			return nil, data.NewEmptyError().AppendDescF("[%s:%s] should be a folder, but its size isn't 0:%d", info.Bucket, info.Key, info.ServerFileSize)
		}

		res.IsExist, _ = utils.ExistDir(f.toAbsFile)
		if !res.IsExist {
			err = utils.CreateDirIfNotExist(f.toAbsFile)
		}
		res.FileModifyTime, _ = utils.LocalFileModify(f.toAbsFile)
		return res, err
	}

	// 文件存在则检查文件状态
	checkMode := -1
	if info.CheckHash {
		checkMode = object.MatchCheckModeFileHash
	} else if info.CheckSize {
		checkMode = object.MatchCheckModeFileSize
	}

	// 读不到 status 按不存在该文件处理
	fileStatus, _ := os.Stat(f.toAbsFile)
	tempFileStatus, _ := os.Stat(f.tempFile)
	if fileStatus != nil {
		// 文件已下载，检测文件内容
		if checkMode < 0 {
			// 文件已存在，无论文件是什么均认为是预期
			res.IsExist = true
			return res, nil
		} else {
			checkResult, mErr := object.Match(object.MatchApiInfo{
				Bucket:         info.Bucket,
				Key:            info.Key,
				LocalFile:      f.toAbsFile,
				CheckMode:      checkMode,
				ServerFileHash: info.ServerFileHash,
				ServerFileSize: info.DownloadFileSize,
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
				if fileModifyTime, fErr := utils.LocalFileModify(f.toAbsFile); fErr != nil {
					log.WarningF("Get file ModifyTime error:%v", fErr)
				} else {
					res.FileModifyTime = fileModifyTime
				}
				return
			}
		}
	} else if tempFileStatus != nil && tempFileStatus.Size() > 0 {
		// 文件已下载了一部分，需要继续下载
		res.IsUpdate = true

		// 下载了一半， 先检查文件是否已经改变，改变则移除临时文件，重新下载
		status, sErr := object.Status(object.StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: false,
		})
		if sErr != nil {
			return nil, data.NewEmptyError().AppendDescF("download part, get file status error:%v", sErr)
		}
		if info.ServerFileSize != status.FSize || (len(info.ServerFileHash) > 0 && info.ServerFileHash != status.Hash) {
			f.fromBytes = 0
			log.DebugF("download part, remove download file because file doesn't match")
			if rErr := os.Remove(f.tempFile); rErr != nil {
				log.ErrorF("download part, remove download file error:%s", rErr)
			}
		}
	}

	// 检查 fromBytes 和 fileSize
	if (info.ServerFileSize+f.fromBytes) > 0 && f.fromBytes >= info.ServerFileSize {
		errorDesc := "download, check fromBytes error: fromBytes bigger than file size, should remove temp file and retry."
		log.Error(errorDesc)
		_ = f.cleanTempFile()
		return nil, data.NewEmptyError().AppendDesc(errorDesc)
	}

	// 下载
	err = download(f, info)
	if err != nil {
		return
	}

	if fStatus, sErr := os.Stat(f.toAbsFile); sErr != nil {
		return res, data.NewEmptyError().AppendDesc("get file stat error after download").AppendError(sErr)
	} else {
		res.FileModifyTime = fStatus.ModTime().Unix()
	}

	// 检查下载后的数据是否符合预期
	if checkMode >= 0 {
		checkResult, mErr := object.Match(object.MatchApiInfo{
			Bucket:         info.Bucket,
			Key:            info.Key,
			LocalFile:      f.toAbsFile,
			CheckMode:      checkMode,
			ServerFileHash: info.ServerFileHash,
			ServerFileSize: info.DownloadFileSize,
		})
		if mErr != nil || (checkResult != nil && !checkResult.Match) {
			log.DebugF("after download, remove download file because file doesn't match")
			if rErr := os.Remove(f.toAbsFile); rErr != nil {
				log.ErrorF("after download, remove download file error:%s", rErr)
			}
			return res, data.NewEmptyError().AppendDesc("check error after download").AppendError(mErr)
		}
	}

	return res, nil
}

func download(fInfo *fileInfo, info *DownloadActionInfo) (err *data.CodeError) {
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
	err = downloadTempFile(fInfo, info)
	if err != nil {
		return err
	}

	err = renameTempFile(fInfo)
	return err
}

func downloadTempFile(fInfo *fileInfo, info *DownloadActionInfo) (err *data.CodeError) {
	useHttps := workspace.GetConfig().IsUseHttps()
	for times := 0; times < 6; times++ {
		dl, cErr := createDownloader(info)
		if cErr != nil {
			return data.NewEmptyError().AppendDesc("Download create downloader error:" + cErr.Error())
		}
		log.DebugF("Download[%d] [%s:%s] => %s", times, info.Bucket, info.Key, info.ToFile)

		h, pErr := info.HostProvider.Provide()
		if h == nil || pErr != nil {
			err = data.NewEmptyError().AppendDescF("no available host:%+v", pErr)
			log.DebugF("Stop download [%s:%s] => %s, because no available host", info.Bucket, info.Key, info.ToFile)
			break
		}

		hostString := h.GetServer()
		urlString, cErr := createDownloadUrlWithHost(h, info, useHttps)
		if len(hostString) == 0 || cErr != nil {
			err = data.NewEmptyError().AppendDescF("create download url error:%+v", cErr)
			log.DebugF("Stop download [%s:%s] => %s, because %+v", info.Bucket, info.Key, info.ToFile, err)
			break
		}

		err = downloadTempFileWithDownloader(dl, fInfo, &DownloadApiInfo{
			Url:            urlString,
			Host:           hostString,
			Referer:        info.Referer,
			RangeFromBytes: fInfo.fromBytes,
			RangeToBytes:   0,
			CheckSize:      info.CheckSize,
			FileSize:       info.ServerFileSize,
			CheckHash:      info.CheckHash,
			FileHash:       info.ServerFileHash,
			Progress:       info.Progress,
		})
		if err == nil {
			break
		}

		log.DebugF("Download[%d] [%s:%s] => %s, err:%+v", times, info.Bucket, info.Key, info.ToFile, err)
		if (err.Code > 399 && err.Code < 500) ||
			err.Code == 612 || err.Code == 631 {
			log.DebugF("Stop download [%s:%s] => %s, because [%+v]", info.Bucket, info.Key, info.ToFile, err)
			break
		}

		info.HostProvider.Freeze(h)
	}
	return err
}

func downloadTempFileWithDownloader(dl downloader, fInfo *fileInfo, info *DownloadApiInfo) *data.CodeError {
	// 下载之前先验证文件大小，两点考虑：
	// 1. 文件可能被更换，非预期文件
	// 2. 文件开启瘦身等，预期文件
	// 上面两点无法区分，但必须能让用户可以下载预期文件
	// 不检测文件信息时，使用下载的文件信息作为标准，可以保证下载成功
	// 此方案有个问题：如果获取文件信息之后，下载之前文件改变了，下载仍会失败（此情景概率极低，且用户重新下载即可）。
	if file, err := utils.GetNetworkFileInfo(info.Url); err != nil {
		return err
	} else if info.CheckHash && info.FileHash != file.Hash {
		return data.NewEmptyError().AppendDescF("file hash doesn't match, %s but except:%s", file.Hash, info.FileHash)
	} else if info.CheckSize && info.FileSize != file.Size {
		return data.NewEmptyError().AppendDescF("file size doesn't match, %d but except:%d", file.Size, info.FileSize)
	} else {
		info.FileSize = file.Size
		info.FileHash = file.Hash
	}

	response, err := dl.Download(info)
	if err != nil {
		return data.NewEmptyError().AppendDesc(" Download error:" + err.Error())
	}
	if response == nil {
		return data.NewEmptyError().AppendDesc(" Download error: response empty")
	}
	if response.StatusCode/100 != 2 {
		return data.NewError(response.StatusCode, "").AppendDescF(" Download error: %v", response)
	}
	if response.Body == nil {
		return data.NewEmptyError().AppendDesc(" Download error: response body empty")
	}

	if response != nil && response.Body != nil {
		if info.Progress != nil {
			info.Progress.SetFileSize(response.ContentLength + info.RangeFromBytes)
			info.Progress.SendSize(info.RangeFromBytes)
			info.Progress.Start()
		}
		defer response.Body.Close()
	}

	var fErr error
	var tempFileHandle *os.File
	isExist, _ := utils.ExistFile(fInfo.tempFile)
	if isExist {
		tempFileHandle, fErr = os.OpenFile(fInfo.tempFile, os.O_APPEND|os.O_WRONLY, 0655)
		log.DebugF("download %s => %s from:%d", info.Url, fInfo.toFile, info.RangeFromBytes)
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

func renameTempFile(fInfo *fileInfo) *data.CodeError {
	err := os.Rename(fInfo.tempFile, fInfo.toAbsFile)
	if err != nil {
		return data.NewEmptyError().AppendDescF(" Rename temp file to final file error:%v", err.Error())
	}
	return nil
}

type downloader interface {
	Download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError)
}

func createDownloader(info *DownloadActionInfo) (downloader, *data.CodeError) {
	// 使用切片，并发至少为 2，至少要能切两片，达到切片阈值
	if info.EnableSlice &&
		info.SliceConcurrentCount > 1 &&
		info.ServerFileSize > info.SliceSize &&
		info.ServerFileSize > info.SliceFileSizeThreshold {
		return &sliceDownloader{
			SliceSize:              info.SliceSize,
			FileHash:               info.ServerFileHash,
			FileEncoding:           info.FileEncoding,
			ToFile:                 info.ToFile,
			slicesDir:              "",
			ConcurrentCount:        info.SliceConcurrentCount,
			totalSliceCount:        0,
			slices:                 nil,
			downloadError:          nil,
			currentReadSliceIndex:  0,
			currentReadSliceOffset: 0,
			locker:                 sync.Mutex{},
		}, nil
	} else {
		return &downloaderFile{}, nil
	}
}

func utf82GBK(text string) (string, *data.CodeError) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	d, err := gbkEncoder.String(text)
	return d, data.ConvertError(err)
}
