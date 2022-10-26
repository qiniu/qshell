package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/progress"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"time"
)

type DownloadInfo struct {
	Bucket                 string // 文件被保存的 bucket
	Key                    string // 文件被保存的 key
	ToFile                 string // 文件保存的路径
	UseGetFileApi          bool   //
	IsPublic               bool   //
	CheckHash              bool   // 是否检测文件 hash
	Domain                 string // 下载的 domain
	RemoveTempWhileError   bool   //
	EnableSlice            bool   // 允许切片下载
	SliceFileSizeThreshold int64  // 允许切片下载，切片下载出发的文件大小阈值，考虑到不希望所有文件都使用切片下载的场景
	SliceSize              int64  // 允许切片下载，切片的大小
	SliceConcurrentCount   int    // 允许切片下载，并发下载切片的个数
}

func (info *DownloadInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func DownloadFile(cfg *iqshell.Config, info DownloadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	// 如果 ToFile 不存在则保存在当前文件录下，文件名为：key
	if len(info.ToFile) == 0 {
		info.ToFile = info.Key
	}

	fileStatus, err := object.Status(object.StatusApiInfo{
		Bucket:   info.Bucket,
		Key:      info.Key,
		NeedPart: false,
	})
	if err != nil {
		log.ErrorF("get file status error:%v", err)
		return
	}

	hostProvider := getDownloadHostProvider(workspace.GetConfig(), &DownloadCfg{
		IoHost:     info.Domain,
		CdnDomain:  info.Domain,
		Bucket:     info.Bucket,
		GetFileApi: info.UseGetFileApi,
	})
	if available, e := hostProvider.Available(); !available {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or io_host or bind domain to bucket, %v", info.Bucket, e)
		return
	}

	apiInfo := &download.ApiInfo{
		IsPublic:               info.IsPublic,
		HostProvider:           hostProvider,
		ToFile:                 info.ToFile,
		Referer:                "",
		FileEncoding:           "",
		Bucket:                 info.Bucket,
		Key:                    info.Key,
		ServerFilePutTime:      fileStatus.PutTime,
		ServerFileSize:         fileStatus.FSize,
		ServerFileHash:         fileStatus.Hash,
		CheckHash:              info.CheckHash,
		FromBytes:              0,
		RemoveTempWhileError:   info.RemoveTempWhileError,
		UseGetFileApi:          info.UseGetFileApi,
		EnableSlice:            info.EnableSlice,
		SliceSize:              info.SliceSize,
		SliceConcurrentCount:   info.SliceConcurrentCount,
		SliceFileSizeThreshold: info.SliceFileSizeThreshold,
		Progress:               progress.NewPrintProgress(" 进度"),
	}
	_, _ = downloadFile(apiInfo)
}

func downloadFile(info *download.ApiInfo) (*download.ApiResult, *data.CodeError) {
	log.InfoF("Download [%s:%s] => %s", info.Bucket, info.Key, info.ToFile)
	startTime := time.Now().UnixNano() / 1e6
	res, err := download.Download(info)
	if err != nil {
		log.ErrorF("Download  Failed, [%s:%s] => %s error:%v", info.Bucket, info.Key, info.ToFile, err)
		return res, err
	}

	fileStatus, sErr := os.Stat(res.FileAbsPath)
	if sErr != nil {
		log.ErrorF("Download  Failed, [%s:%s] => %s get file status error:%v", info.Bucket, info.Key, info.ToFile, err)
		return res, data.ConvertError(sErr)
	}
	if fileStatus == nil {
		log.ErrorF("Download  Failed, [%s:%s] => %s download speed: can't get file status", info.Bucket, info.Key, info.ToFile)
		return res, data.NewEmptyError().AppendDesc("can't get file status")
	}

	endTime := time.Now().UnixNano() / 1e6
	duration := float64(endTime-startTime) / 1000
	speed := "0KB/s"
	if duration > 0.1 {
		// 小于 100ms, 可能没有进行下载操作
		speed = fmt.Sprintf("%.2fKB/s", float64(fileStatus.Size())/duration/1024)
	}
	if res.IsExist {
		log.InfoF("Download Skip because file exist, [%s:%s] => %s", info.Bucket, info.Key, res.FileAbsPath)
	} else if res.IsUpdate {
		log.InfoF("Download update Success, [%s:%s] => %s duration:%.2fs speed:%s", info.Bucket, info.Key, res.FileAbsPath, duration, speed)
	} else {
		log.InfoF("Download Success, [%s:%s] => %s duration:%.2fs speed:%s", info.Bucket, info.Key, res.FileAbsPath, duration, speed)
	}

	return res, nil
}
