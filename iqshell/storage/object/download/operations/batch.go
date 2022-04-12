package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type BatchDownloadWithConfigInfo struct {
	BatchInfo           batch.Info
	LocalDownloadConfig string
}

func (info *BatchDownloadWithConfigInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}
	return nil
}

func BatchDownloadWithConfig(cfg *iqshell.Config, info BatchDownloadWithConfigInfo) {
	if iqshell.ShowDocumentIfNeeded(cfg) {
		return
	}

	if !iqshell.Check(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}) {
		return
	}

	downloadInfo := BatchDownloadInfo{
		BatchInfo:   info.BatchInfo,
		DownloadCfg: DefaultDownloadCfg(),
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
	BatchInfo batch.Info
	DownloadCfg
}

func (info *BatchDownloadInfo) Check() *data.CodeError {
	if info.BatchInfo.WorkerCount < 1 || info.BatchInfo.WorkerCount > 2000 {
		info.BatchInfo.WorkerCount = 5
	}
	if err := info.BatchInfo.Check(); err != nil {
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
		BeforeLoadFileLog: func() {
			if len(info.RecordRoot) == 0 {
				info.RecordRoot = downloadCachePath(workspace.GetConfig(), &info.DownloadCfg)
			}
			if data.Empty(cfg.CmdCfg.Log.LogFile) {
				workspace.GetConfig().Log.LogFile = data.NewString(filepath.Join(info.RecordRoot, "log.txt"))
			}
		},
	}); !shouldContinue {
		return
	}

	log.InfoF("record root: %s", info.RecordRoot)

	info.BatchInfo.InputFile = info.KeyFile
	hostProvider := getDownloadHostProvider(workspace.GetConfig(), &info.DownloadCfg)
	if available, e := hostProvider.Available(); !available {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or bind domain to bucket; %v", info.Bucket, e)
		return
	}

	dbPath := filepath.Join(info.RecordRoot, ".ldb")
	log.InfoF("download db dir:%s", dbPath)

	exporter, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:   info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:      info.BatchInfo.FailExportFilePath,
		OverwriteExportFilePath: info.BatchInfo.OverwriteExportFilePath,
	})
	if err != nil {
		log.Error(err)
		return
	}

	timeStart := time.Now()
	var locker sync.Mutex
	var totalFileCount int64
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

	flow.New(info.BatchInfo.Info).
		WorkProvider(NewWorkProvider(info.Bucket, info.BatchInfo.InputFile, info.BatchInfo.ItemSeparate)).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				apiInfo := workInfo.Work.(*download.ApiInfo)
				apiInfo.HostProvider = hostProvider
				apiInfo.ToFile = filepath.Join(info.DestDir, apiInfo.Key)
				apiInfo.StatusDBPath = dbPath
				apiInfo.Referer = info.Referer
				apiInfo.FileEncoding = info.FileEncoding
				apiInfo.Bucket = info.Bucket
				apiInfo.RemoveTempWhileError = info.RemoveTempWhileError
				apiInfo.UseGetFileApi = info.GetFileApi
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

				if file, e := downloadFile(apiInfo); e != nil {
					return nil, e
				} else {
					return file, nil
				}
			}), nil
		})).
		ShouldSkip(func(workInfo *flow.WorkInfo) (skip bool, cause *data.CodeError) {
			apiInfo := workInfo.Work.(*download.ApiInfo)
			if filterPrefix(apiInfo.Key) {
				return true, data.NewEmptyError().AppendDescF("Skip download `%s`, suffix filter not match", apiInfo.Key)
			}
			return false, nil
		}).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			totalFileCount = flow.WorkProvider.WorkTotalCount()
			return nil
		}).
		OnWorkSkip(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			locker.Lock()
			skipBySuffixes += 1
			locker.Unlock()

			log.Info(err.Error())
			exporter.Skip().Export(workInfo.Data)
		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			res := result.(download.ApiResult)
			locker.Lock()
			if res.IsExist {
				existsFileCount += 1
			} else if res.IsUpdate {
				updateFileCount += 1
			} else {
				successFileCount += 1
			}
			locker.Unlock()

			exporter.Success().Export(workInfo.Data)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			locker.Lock()
			failureFileCount += 1
			locker.Unlock()

			exporter.Fail().ExportF("%s%s%s", workInfo.Data, flow.ErrorSeparate, err)
			log.ErrorF("Download  Failed, %s error:%v", workInfo.Data, err)
		}).Build().Start()

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

	if workspace.GetConfig().Log.Enable() {
		log.AlertF("See download log at path:%s", workspace.GetConfig().Log.LogFile.Value())
	}
}
