package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type BatchDownloadWithConfigInfo struct {
	GroupInfo           group.Info
	LocalDownloadConfig string
}

func (info *BatchDownloadWithConfigInfo) Check() error {
	if err := info.GroupInfo.Check(); err != nil {
		return err
	}
	return nil
}

func BatchDownloadWithConfig(cfg *iqshell.Config, info BatchDownloadWithConfigInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	downloadInfo := BatchDownloadInfo{
		GroupInfo:   info.GroupInfo,
		DownloadCfg: DownloadCfg{},
	}
	if err := utils.UnMarshalFromFile(info.LocalDownloadConfig, &downloadInfo.DownloadCfg); err != nil {
		log.ErrorF("UnMarshal: read download config error:%v config file:%s", info.LocalDownloadConfig, err)
		return
	}
	if err := utils.UnMarshalFromFile(info.LocalDownloadConfig, cfg.CmdCfg.Log); err != nil {
		log.ErrorF("UnMarshal: read log setting error:%v config file:%s", info.LocalDownloadConfig, err)
		return
	}
	BatchDownload(cfg, downloadInfo)
}

type BatchDownloadInfo struct {
	GroupInfo group.Info
	DownloadCfg
}

func (info *BatchDownloadInfo) Check() error {
	if info.GroupInfo.WorkerCount < 1 || info.GroupInfo.WorkerCount > 2000 {
		info.GroupInfo.WorkerCount = 5
	}
	if err := info.GroupInfo.Check(); err != nil {
		return err
	}
	if err := info.DownloadCfg.Check(); err != nil {
		return err
	}
	return nil
}

func BatchDownload(cfg *iqshell.Config, info BatchDownloadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if data.NotEmpty(workspace.GetConfig().Log.LogFile) {
		log.AlertF("Writing upload log to file:%s \n\n", workspace.GetConfig().Log.LogFile.Value())
	} else {
		log.Debug("log file not set \n\n")
	}

	info.GroupInfo.InputFile = info.KeyFile
	downloadDomain, downloadHost := getDownloadDomainAndHost(workspace.GetConfig(), &info.DownloadCfg)
	if len(downloadDomain) == 0 && len(downloadHost) == 0 {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or io_host or bind domain to bucket", info.Bucket)
		return
	}

	log.DebugF("Download Domain:%s", downloadDomain)
	log.DebugF("Download Domain:%s", downloadHost)

	if len(info.RecordRoot) == 0 {
		info.RecordRoot = workspace.GetWorkspace()
	}
	jobId := info.DownloadCfg.JobId()
	cachePath := downloadCachePath(workspace.GetConfig(), &info.DownloadCfg)
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

	ds, err := newDownloadScanner(info.KeyFile, info.GroupInfo.ItemSeparate, info.Bucket, exporter)
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

	hasPrefixes := len(info.Prefix) > 0
	prefixes := strings.Split(info.Prefix, ",")
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
	work.NewFlowHandler(info.GroupInfo.FlowInfo).ReadWork(func() (work work.Work, hasMore bool) {
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
		apiInfo.Domain = downloadDomain
		apiInfo.Host = downloadHost
		apiInfo.ToFile = filepath.Join(info.DestDir, apiInfo.Key)
		apiInfo.StatusDBPath = dbPath
		apiInfo.Referer = info.Referer
		apiInfo.FileEncoding = info.FileEncoding
		apiInfo.Bucket = info.Bucket
		apiInfo.UserGetFileApi = info.GetFileApi
		if !info.CheckHash {
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

		file, err := downloadFile(apiInfo)

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

	if data.Empty(workspace.GetConfig().Log.LogFile) {
		log.AlertF("See download log at path:%s", workspace.GetConfig().Log.LogFile.Value())
	}

	if failureFileCount > 0 {
		os.Exit(data.StatusError)
	}
}
