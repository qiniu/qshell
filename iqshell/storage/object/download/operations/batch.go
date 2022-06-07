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
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"os"
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
		log.ErrorF("Download, %v", e)
		return
	}

	unlockHandler := func() {
		if e := locker.TryUnlock(); e != nil {
			log.ErrorF("Download, %v", e)
		}
	}
	workspace.AddCancelObserver(func(s os.Signal) {
		unlockHandler()
	})
	defer unlockHandler()

	info.InputFile = info.KeyFile
	hostProvider := getDownloadHostProvider(workspace.GetConfig(), &info.DownloadCfg)
	if available, e := hostProvider.Available(); !available {
		log.ErrorF("get download domain error: not find in config and can't get bucket(%s) domain, you can set cdn_domain or bind domain to bucket; %v", info.Bucket, e)
		return
	}

	dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
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

	hasSuffixes := len(info.Suffixes) > 0
	suffixes := strings.Split(info.Suffixes, ",")
	filterSuffixes := func(name string) bool {
		if !hasSuffixes {
			return false
		}

		for _, suffix := range suffixes {
			if strings.HasSuffix(name, suffix) {
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
				apiInfo.Referer = info.Referer
				apiInfo.FileEncoding = info.FileEncoding
				apiInfo.Bucket = info.Bucket
				apiInfo.RemoveTempWhileError = info.RemoveTempWhileError
				apiInfo.UseGetFileApi = info.GetFileApi
				if !info.CheckHash {
					apiInfo.FileHash = ""
				}

				metric.AddCurrentCount(1)
				metric.PrintProgress("Downloading " + apiInfo.Key)

				if file, e := downloadFile(apiInfo); e != nil {
					return nil, e
				} else {
					return file, nil
				}
			}), nil
		})).
		SetDBOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: &download.ApiInfo{},
				},
				Result: &download.ApiResult{},
				Err:    nil,
			}
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err != nil {
				return true, workRecord.Err
			}

			apiInfo, _ := workInfo.Work.(*download.ApiInfo)
			recordApiInfo, _ := workRecord.Work.(*download.ApiInfo)

			result, _ := workRecord.Result.(*download.ApiResult)
			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDesc("result is invalid")
			}

			// 本地文件和服务端文件均没有变化，则不需要重新下载
			if match, _ := utils.IsFileMatchFileModifyTime(apiInfo.ToFile, result.FileModifyTime);
				match && apiInfo.FileModifyTime == recordApiInfo.FileModifyTime {
				return false, nil
			}

			// 本地或服务端文件有变动，则先查 size，size 不同则需要重新下载， 相同再尝试检查 hash
			// 检测文件大小
			if _, cause = utils.IsFileMatchFileSize(apiInfo.ToFile, apiInfo.FileSize); err != nil ||
				apiInfo.FileSize != recordApiInfo.FileSize{
				return true, cause
			}

			if info.CheckHash {
				if _, mErr := object.Match(object.MatchApiInfo{
					Bucket:    apiInfo.Bucket,
					Key:       apiInfo.Key,
					FileHash:  apiInfo.FileHash,
					LocalFile: apiInfo.ToFile,
				}); mErr != nil {
					return true, mErr
				}
			}

			return false, nil
		}).
		ShouldSkip(func(workInfo *flow.WorkInfo) (skip bool, cause *data.CodeError) {
			apiInfo, _ := workInfo.Work.(*download.ApiInfo)
			if filterPrefix(apiInfo.Key) {
				return true, data.NewEmptyError().AppendDescF("Skip download `%s`, prefix filter not match", apiInfo.Key)
			}
			if filterSuffixes(apiInfo.Key) {
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
			res, _ := result.(*download.ApiResult)
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
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.UpdateCount + metric.ExistCount + metric.SkippedCount
	}

	resultPath := filepath.Join(workspace.GetJobDir(), ".result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save download result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save download result to path:%s", resultPath)
	}

	log.Info("-------Download Result-------")
	log.InfoF("%10s%10d", "Total:", metric.TotalCount)
	log.InfoF("%10s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%10s%10d", "Exists:", metric.ExistCount)
	log.InfoF("%10s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%10s%10d", "Update:", metric.UpdateCount)
	log.InfoF("%10s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%10s%10ds", "Duration:", metric.Duration)
	log.InfoF("-----------------------------")
	if workspace.GetConfig().Log.Enable() {
		log.InfoF("See download log at path:%s", workspace.GetConfig().Log.LogFile.Value())
	}
}
