package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/locker"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"path/filepath"
	"strings"
)

type BatchDownloadWithConfigInfo struct {
	flow.Info
	export.FileExporterConfig

	// 工作数据源
	InputFile    string // 工作数据源：文件
	ItemSeparate string // 工作数据源：每行元素按分隔符分的分隔符

	LocalDownloadConfig string
}

func (info *BatchDownloadWithConfigInfo) Check() *data.CodeError {
	if err := info.Info.Check(); err != nil {
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
		Info:               info.Info,
		FileExporterConfig: info.FileExporterConfig,
		InputFile:          info.InputFile,
		ItemSeparate:       info.ItemSeparate,
		DownloadCfg:        DefaultDownloadCfg(),
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
	flow.Info
	export.FileExporterConfig
	DownloadCfg

	// 工作数据源
	InputFile    string // 工作数据源：文件
	ItemSeparate string // 工作数据源：每行元素按分隔符分的分隔符
}

func (info *BatchDownloadInfo) Check() *data.CodeError {
	if info.WorkerCount < 1 || info.WorkerCount > 2000 {
		info.WorkerCount = 5
	}
	if err := info.Info.Check(); err != nil {
		return err
	}
	if err := info.DownloadCfg.Check(); err != nil {
		return err
	}
	return nil
}

func BatchDownload(cfg *iqshell.Config, info BatchDownloadInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		if len(info.RecordRoot) > 0 {
			return info.RecordRoot
		}
		return filepath.Join(cmdPath, info.JobId())
	}
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	// 配置 locker
	if e := locker.TryLock(); e != nil {
		log.ErrorF("batch job, %v", e)
		return
	}

	info.InputFile = info.KeyFile
	hostProvider := getDownloadHostProvider(workspace.GetConfig(), &info.DownloadCfg)
	if available, e := hostProvider.Available(); !available {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or bind domain to bucket; %v", info.Bucket, e)
		return
	}

	dbPath := filepath.Join(workspace.GetJobDir(), ".ldb")
	log.InfoF("download db dir:%s", dbPath)

	exporter, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:   info.SuccessExportFilePath,
		FailExportFilePath:      info.FailExportFilePath,
		OverwriteExportFilePath: info.OverwriteExportFilePath,
	})
	if err != nil {
		log.Error(err)
		return
	}

	metric := &Metric{}
	metric.Start()
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

	flow.New(info.Info).
		WorkProvider(NewWorkProvider(info.Bucket, info.InputFile, info.ItemSeparate)).
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

				metric.AddCurrentCount(1)

				if metric.TotalCount > 0 {
					log.AlertF("Downloading %s [%d/%d, %.1f%%] ...", apiInfo.Key, metric.CurrentCount, metric.TotalCount,
						float32(metric.CurrentCount)*100/float32(metric.TotalCount))
				} else {
					log.AlertF("Downloading %s [%d/-, -] ...", apiInfo.Key, metric.CurrentCount)
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
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		OnWorkSkip(func(workInfo *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddSkippedCount(1)

			log.Info(err.Error())
			exporter.Skip().Export(workInfo.Data)
		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			res := result.(*download.ApiResult)
			if res.IsExist {
				metric.AddExistCount(1)
			} else if res.IsUpdate {
				metric.AddUpdateCount(1)
			} else {
				metric.AddSuccessCount(1)
			}

			exporter.Success().Export(workInfo.Data)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddFailureCount(1)

			exporter.Fail().ExportF("%s%s%s", workInfo.Data, flow.ErrorSeparate, err)
			log.ErrorF("Download  Failed, %s error:%v", workInfo.Data, err)
		}).Build().Start()

	metric.End()
	if metric.TotalCount == 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.UpdateCount + metric.ExistCount + metric.SkippedCount
	}

	log.Alert("-------Download Result-------")
	log.AlertF("%10s%10d", "Total:", metric.TotalCount)
	log.AlertF("%10s%10d", "Skipped:", metric.SkippedCount)
	log.AlertF("%10s%10d", "Exists:", metric.ExistCount)
	log.AlertF("%10s%10d", "Success:", metric.SuccessCount)
	log.AlertF("%10s%10d", "Update:", metric.UpdateCount)
	log.AlertF("%10s%10d", "Failure:", metric.FailureCount)
	log.AlertF("%10s%10ds", "Duration:", metric.Duration)
	log.AlertF("-----------------------------")

	if workspace.GetConfig().Log.Enable() {
		log.AlertF("See download log at path:%s", workspace.GetConfig().Log.LogFile.Value())
	}

	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save download result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save download result to path:%s", resultPath)
	}

	if e := locker.TryUnlock(); e != nil {
		log.ErrorF("Download, %v", e)
		return
	}
}
