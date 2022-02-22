package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type BatchDownloadInfo struct {
	GroupInfo group.Info
}

func (info *BatchDownloadInfo) Check() error {
	if info.GroupInfo.WorkCount < 1 || info.GroupInfo.WorkCount > 2000 {
		info.GroupInfo.WorkCount = 5
	}
	if err := info.GroupInfo.Check(); err != nil {
		return err
	}
	return nil
}

func BatchDownload(info BatchDownloadInfo) {
	downloadCfg := workspace.GetConfig().Download
	info.GroupInfo.InputFile = downloadCfg.KeyFile

	if err := downloadCfg.Check(); err != nil {
		log.ErrorF("download config check error:%v", err)
		return
	}

	downloadDomain := downloadCfg.DownloadDomain()
	if len(downloadDomain) == 0 {
		downloadDomain, _ = bucket.DomainOfBucket(downloadCfg.Bucket)
	}
	if len(downloadDomain) == 0 {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or io_host or bind domain to bucket", downloadCfg.Bucket)
		return
	}

	jobId := downloadCfg.JobId()
	cachePath := workspace.DownloadCachePath()
	dbPath := filepath.Join(cachePath, jobId, ".list")
	log.InfoF("download db dir:%s", dbPath)

	exporter, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.GroupInfo.SuccessExportFilePath,
		FailExportFilePath:     info.GroupInfo.FailExportFilePath,
		OverrideExportFilePath: info.GroupInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.Error(err)
		return
	}

	ds, err := newDownloadScanner(downloadCfg.KeyFile, info.GroupInfo.ItemSeparate, downloadCfg.Bucket, exporter)
	if err != nil {
		log.Error(err)
		return
	}

	timeStart := time.Now()
	var locker sync.Mutex
	var totalFileCount = ds.getFileLineCount()
	var currentFileCount int64
	var existsFileCount int64
	var updateFileCount int64
	var successFileCount int64
	var failureFileCount int64
	var skipBySuffixes int64

	hasPrefixes := len(downloadCfg.Prefix) > 0
	prefixes := strings.Split(downloadCfg.Prefix, ",")
	filterPrefix := func(name string) bool {
		if !hasPrefixes {
			return false
		}

		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				return false
			}
		}
		return true
	}
	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		apiInfo, hasMore := ds.scan()
		if apiInfo == nil {
			return
		}

		if filterPrefix(apiInfo.Key) {
			log.AlertF("Skip download `%s`, suffix filter not match", apiInfo.Key)
			locker.Lock()
			skipBySuffixes += 1
			locker.Unlock()
			return nil, true
		}

		return apiInfo, hasMore
	}).DoWork(func(work work.Work) (work.Result, error) {
		apiInfo := work.(*download.ApiInfo)
		apiInfo.Url = "" // downloadFile 时会自动创建
		apiInfo.Domain = downloadDomain
		apiInfo.ToFile = filepath.Join(downloadCfg.DestDir, apiInfo.Key)
		apiInfo.StatusDBPath = dbPath
		apiInfo.Referer = downloadCfg.Referer
		apiInfo.FileEncoding = downloadCfg.FileEncoding
		apiInfo.Bucket = downloadCfg.Bucket
		if !downloadCfg.CheckHash {
			apiInfo.FileHash = ""
		}

		locker.Lock()
		currentFileCount += 1
		locker.Unlock()

		if totalFileCount > 0 {
			log.AlertF("Downloading %s [%d/%d, %.1f%%] ...", apiInfo.Key, currentFileCount, totalFileCount,
				float32(currentFileCount)*100/float32(totalFileCount))
		} else {
			log.AlertF("Downloading %s [%d/-, -] ...", apiInfo.Key, currentFileCount)
		}

		file, err := downloadFile(DownloadInfo{
			ApiInfo:  *apiInfo,
			IsPublic: downloadCfg.Public,
		})

		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(*download.ApiInfo)
		res := result.(download.ApiResult)
		exporter.Success().ExportF("download success, [%s:%s] => %s", apiInfo.Bucket, apiInfo.Key, res.FileAbsPath)

		locker.Lock()
		if res.IsExist {
			existsFileCount += 1
		} else if res.IsUpdate {
			updateFileCount += 1
		} else {
			successFileCount += 1
		}
		locker.Unlock()
	}).OnWorkError(func(work work.Work, err error) {
		locker.Lock()
		failureFileCount += 1
		locker.Unlock()

		apiInfo := work.(*download.ApiInfo)
		exporter.Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* key fileSize fileHash and fileModifyTime */
			apiInfo.Key, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileHash, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()

	if totalFileCount == 0 {
		totalFileCount = skipBySuffixes + existsFileCount + successFileCount + updateFileCount + failureFileCount
	} else {
		skipBySuffixes = totalFileCount - existsFileCount - successFileCount - updateFileCount - failureFileCount
	}

	log.Alert("-------Download Result-------")
	log.AlertF("%10s%10d", "Total:", totalFileCount)
	log.AlertF("%10s%10d", "Skipped:", skipBySuffixes)
	log.AlertF("%10s%10d", "Exists:", existsFileCount)
	log.AlertF("%10s%10d", "Success:", successFileCount)
	log.AlertF("%10s%10d", "Update:", updateFileCount)
	log.AlertF("%10s%10d", "Failure:", failureFileCount)
	log.AlertF("%10s%15s", "Duration:", time.Since(timeStart))
	log.AlertF("-----------------------------")
	log.AlertF("See download log at path:%s", downloadCfg.LogFile)

	if failureFileCount > 0 {
		os.Exit(data.StatusError)
	}
}
