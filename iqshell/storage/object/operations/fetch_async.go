package operations

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

type CheckAsyncFetchStatusInfo struct {
	Bucket string
	Id     string
}

func (info *CheckAsyncFetchStatusInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Id) == 0 {
		return alert.CannotEmptyError("Id", "")
	}
	return nil
}

func CheckAsyncFetchStatus(cfg *iqshell.Config, info CheckAsyncFetchStatusInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	ret, err := object.CheckAsyncFetchStatus(info.Bucket, info.Id)
	if err != nil {
		log.ErrorF("CheckAsyncFetchStatus error: %v", err)
	} else {
		log.Alert(ret)
	}
}

type BatchAsyncFetchInfo struct {
	BatchInfo               batch.Info
	Bucket                  string // fetch 的目的 bucket
	Host                    string // 从指定URL下载时指定的HOST
	CallbackUrl             string // 抓取成功的回调地址
	CallbackBody            string //
	CallbackBodyType        string //
	CallbackHost            string // 回调时使用的HOST
	FileType                int    // 文件存储类型， 0 标准存储， 1 低频存储
	Overwrite               bool   //
	DisableCheckFetchResult bool   // 不检测是否 fetch 成功
}

func (info *BatchAsyncFetchInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func BatchAsyncFetch(cfg *iqshell.Config, info BatchAsyncFetchInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", cfg.CmdCfg.CmdId, info.Bucket, info.BatchInfo.InputFile))
		return filepath.Join(cmdPath, jobId)
	}
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	info.BatchInfo.Force = true
	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}

	fetchResultChan := make(chan flow.Work, info.BatchInfo.Info.WorkerCount*10)

	wait := &sync.WaitGroup{}
	wait.Add(2)

	// fetch
	go func() {
		batchAsyncFetch(cfg, info, exporter, fetchResultChan)
		close(fetchResultChan)
		wait.Done()
	}()

	// check
	go func() {
		batchAsyncFetchCheck(cfg, info, exporter, fetchResultChan)
		wait.Done()
	}()

	wait.Wait()
}

func batchAsyncFetch(cfg *iqshell.Config, info BatchAsyncFetchInfo,
	exporter *export.FileExporter, fetchResultChan chan<- flow.Work) {

	metric := &batch.Metric{}
	metric.Start()

	dbPath := filepath.Join(workspace.GetJobDir(), "fetch.recorder")
	if info.BatchInfo.EnableRecord {
		log.DebugF("batch async fetch recorder:%s", dbPath)
	} else {
		log.Debug("batch async fetch recorder:Not Enable")
	}

	flow.New(info.BatchInfo.Info).
		WorkProviderWithFile(info.BatchInfo.InputFile,
			info.BatchInfo.EnableStdin,
			flow.NewItemsWorkCreator(info.BatchInfo.ItemSeparate, 1, func(items []string) (work flow.Work, err *data.CodeError) {
				var size uint64 = 0
				fromUrl := items[0]
				if len(items) > 1 {
					s, pErr := strconv.ParseUint(items[1], 10, 64)
					if pErr != nil {
						return nil, alert.Error("parse size error:"+pErr.Error(), "")
					}
					size = s
				}

				saveKey := ""
				if len(items) > 2 && len(items[2]) > 0 {
					saveKey = items[2]
				} else {
					key, pErr := utils.KeyFromUrl(fromUrl)
					if pErr != nil || len(key) == 0 {
						return nil, alert.Error("get key form url error:"+pErr.Error()+" check url style", "")
					}
					saveKey = key
				}

				return &asyncFetchItem{
					fileSize: size,
					info: object.AsyncFetchApiInfo{
						Url:              fromUrl,
						Host:             info.Host,
						Bucket:           info.Bucket,
						Key:              saveKey,
						Md5:              "", // 设置了该值，抓取的过程使用文件md5值进行校验, 校验失败不存在七牛空间
						Etag:             "", // 设置了该值， 抓取的过程中使用etag进行校验，失败不保存在存储空间中
						CallbackURL:      info.CallbackUrl,
						CallbackBody:     info.CallbackBody,
						CallbackBodyType: info.CallbackBodyType,
						FileType:         info.FileType,
						IgnoreSameKey:    !info.Overwrite, // 此处需要翻转逻辑
					},
				}, nil
			})).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in, _ := workInfo.Work.(*asyncFetchItem)

				metric.PrintProgress(fmt.Sprintf("Fetching, %s => [%s:%s]", in.info.Url, in.info.Bucket, in.info.Key))

				result, e := object.AsyncFetch(in.info)
				return &asyncFetchResult{
					Bucket:   in.info.Bucket,
					Key:      in.info.Key,
					Url:      in.info.Url,
					FileSize: in.fileSize,
					Info:     result,
				}, e
			}), nil
		})).
		SetOverseerEnable(info.BatchInfo.EnableRecord).
		SetDBOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: &asyncFetchItem{},
				},
				Result: &asyncFetchResult{},
				Err:    nil,
			}
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err == nil {
				return false, nil
			}
			if !info.BatchInfo.RecordRedoWhileError {
				return false, workRecord.Err
			}

			result, _ := workRecord.Result.(*object.FetchResult)
			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDesc("result is invalid")
			}
			return false, nil
		}).
		FlowWillStartFunc(func(flow *flow.Flow) (err *data.CodeError) {
			metric.AddTotalCount(flow.WorkProvider.WorkTotalCount())
			return nil
		}).
		OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if result != nil && result.IsValid() {
					metric.AddSuccessCount(1)
					if info.DisableCheckFetchResult {
						exporter.Success().ExportF("%s", work.Data)
					}
					log.InfoF("Fetch skip line:%s because have done and success", work.Data)
					// 成功的任务需要添加到队列中等待检查（有些任务可能未来得及检查用户取消，有些检查时失败，重新检查可能成功）
					if res, ok := result.(*asyncFetchResult); ok {
						fetchResultChan <- res
					}
				} else {
					metric.AddFailureCount(1)
					// 不进行检查，需要导出失败条目
					exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
					log.InfoF("Fetch skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				// 不进行检查，需要导出失败条目
				exporter.Fail().ExportF("%s%s%v", work.Data, flow.ErrorSeparate, err)
				log.InfoF("Fetch skip line:%s because:%v", work.Data, err)
			}
		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			metric.AddSuccessCount(1)
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			in := workInfo.Work.(*asyncFetchItem)
			res := result.(*asyncFetchResult)
			if info.DisableCheckFetchResult {
				exporter.Success().ExportF("%s", workInfo.Data)
			}
			fetchResultChan <- res
			log.InfoF("Fetch Response, '%s' => [%s:%s] id:%s wait:%d",
				in.info.Url, in.info.Bucket, in.info.Key, res.Info.Id, res.Info.Wait)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddFailureCount(1)
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			if in, ok := workInfo.Work.(*asyncFetchItem); ok {
				exporter.Fail().ExportF("%s%s%v", in.info.Url, flow.ErrorSeparate, err)
				log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error: %v", in.info.Url, in.info.Bucket, in.info.Key, err)
			} else {
				// 不进行检查，需要导出失败条目
				exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
				log.ErrorF("Fetch Failed, %s, Error: %v", workInfo.Data, err)
			}
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	resultPath := filepath.Join(workspace.GetJobDir(), "fetch.result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save batch async fetch result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save batch async fetch result to path:%s", resultPath)
	}

	log.Info("")
	log.Info("--------------- Batch Fetch Result ---------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("--------------------------------------------------")
}

func batchAsyncFetchCheck(cfg *iqshell.Config, info BatchAsyncFetchInfo,
	exporter *export.FileExporter, fetchResultChan <-chan flow.Work) {
	if info.DisableCheckFetchResult {
		log.DebugF("batch async fetch check: disable")
		for r := range fetchResultChan {
			r.WorkId()
		}
		return
	}

	metric := &batch.Metric{}
	metric.Start()

	dbPath := filepath.Join(workspace.GetJobDir(), "check.recorder")
	if info.BatchInfo.EnableRecord {
		log.DebugF("batch async fetch check recorder:%s", dbPath)
	} else {
		log.Debug("batch async fetch check recorder:Not Enable")
	}

	flow.New(info.BatchInfo.Info).
		WorkProviderWithChan(fetchResultChan).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*asyncFetchResult)

				metric.AddCurrentCount(1)
				metric.PrintProgress(fmt.Sprintf("Checking, %s => [%s:%s]", in.Url, in.Bucket, in.Key))

				checkTimes := 0
				maxDuration := asyncFetchCheckMaxDuration(in.FileSize)
				minDuration := 2
				checkStartTime := time.Now().Add(time.Duration(minDuration) * time.Second)
				checkEndTime := time.Now().Add(time.Duration(maxDuration) * time.Second)
				for {
					current := time.Now()
					if current.After(checkStartTime) {
						checkTimes += 1
						ret, cErr := object.CheckAsyncFetchStatus(in.Bucket, in.Info.Id)
						log.DebugF("batch async fetch check, bucket:%s key:%s id:%s wait:%d", in.Bucket, in.Key, in.Key, ret.Wait)
						if cErr != nil {
							log.ErrorF("CheckAsyncFetchStatus: %v", cErr)
						} else if ret.Wait == -1 { // 视频抓取过一次，有可能成功了，有可能失败了
							if exist, err := object.Exist(object.ExistApiInfo{
								Bucket: in.Bucket,
								Key:    in.Key,
							}); exist {
								return in, nil
							} else {
								log.ErrorF("Check Stat:%s error:%v ID:%s", in.Key, err, in.Info.Id)
							}
						}
					}

					if checkTimes == 0 || current.Before(checkEndTime) {
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}
				return nil, data.NewEmptyError().AppendDesc("can't find object in bucket")
			}), nil
		})).
		SetOverseerEnable(info.BatchInfo.EnableRecord).
		SetDBOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: &asyncFetchResult{},
				},
				Result: &asyncFetchResult{},
				Err:    nil,
			}
		}).
		ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
			if workRecord.Err == nil {
				return false, nil
			}

			if !info.BatchInfo.RecordRedoWhileError {
				return false, workRecord.Err
			}

			result, _ := workRecord.Result.(*asyncFetchResult)
			if result == nil {
				return true, data.NewEmptyError().AppendDesc("no result found")
			}
			if !result.IsValid() {
				return true, data.NewEmptyError().AppendDesc("result is invalid")
			}
			return false, nil
		}).
		OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
			metric.AddCurrentCount(1)
			metric.PrintProgress("Batching:" + work.Data)

			in := work.Work.(*asyncFetchResult)
			operationResult, _ := result.(*asyncFetchResult)
			if err != nil && err.Code == data.ErrorCodeAlreadyDone {
				if operationResult != nil && operationResult.IsValid() {
					metric.AddSuccessCount(1)
					exporter.Success().ExportF("%s\t%s", in.Url, in.Key)
					log.InfoF("Check skip line:%s because have done and success", work.Data)
				} else {
					metric.AddFailureCount(1)
					exporter.Fail().ExportF("%s%s%v", in.Url, flow.ErrorSeparate, err)
					log.InfoF("Check skip line:%s because have done and failure, %v", work.Data, err)
				}
			} else {
				metric.AddSkippedCount(1)
				exporter.Fail().ExportF("%s%s%v", in.Url, flow.ErrorSeparate, err)
				log.InfoF("Check skip line:%s because:%v", work.Data, err)
			}
		}).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			metric.AddSuccessCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			in := workInfo.Work.(*asyncFetchResult)
			exporter.Success().ExportF("%s\t%s", in.Url, in.Key)
			log.InfoF("Fetch Success, %s => [%s:%s]", in.Url, in.Bucket, in.Key)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			metric.AddFailureCount(1)
			metric.PrintProgress("Batching:" + workInfo.Data)

			if in, ok := workInfo.Work.(*asyncFetchResult); ok {
				exporter.Fail().ExportF("%s%s%v", in.Url, flow.ErrorSeparate, err)
				log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error: %v", in.Url, in.Bucket, in.Key, err)
			} else {
				exporter.Fail().ExportF("%s\t%d\t%s", in.Url, in.FileSize, in.Key)
				log.ErrorF("Fetch Failed, %s => [%s:%s], ID:%s", in.Url, in.Bucket, in.Key, in.Info.Id)
			}
		}).Build().Start()

	metric.End()
	if metric.TotalCount <= 0 {
		metric.TotalCount = metric.SuccessCount + metric.FailureCount + metric.SkippedCount
	}

	// 输出结果
	resultPath := filepath.Join(workspace.GetJobDir(), "check.result")
	if e := utils.MarshalToFile(resultPath, metric); e != nil {
		log.ErrorF("save batch async fetch check result to path:%s error:%v", resultPath, e)
	} else {
		log.DebugF("save batch async fetch check result to path:%s", resultPath)
	}

	log.Info("")
	log.Info("------------ Batch Fetch Check Result ------------")
	log.InfoF("%20s%10d", "Total:", metric.TotalCount)
	log.InfoF("%20s%10d", "Success:", metric.SuccessCount)
	log.InfoF("%20s%10d", "Failure:", metric.FailureCount)
	log.InfoF("%20s%10d", "Skipped:", metric.SkippedCount)
	log.InfoF("%20s%10ds", "Duration:", metric.Duration)
	log.InfoF("--------------------------------------------------")
}

func asyncFetchCheckMaxDuration(size uint64) int {
	duration := 10
	if size >= 500*utils.MB {
		duration = 600
	} else if size >= 200*utils.MB {
		duration = 300
	} else if size >= 100*utils.MB {
		duration = 150
	} else if size >= 50*utils.MB {
		duration = 60
	} else if size >= 10*utils.MB {
		duration = 25
	} else if size == 0 {
		duration = 120
	}
	return duration
}

type asyncFetchItem struct {
	fileSize uint64
	info     object.AsyncFetchApiInfo
}

func (f *asyncFetchItem) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", f.info.Url, f.info.Bucket, f.info.Key)
}

type asyncFetchResult struct {
	Bucket   string                      `json:"bucket"`
	Key      string                      `json:"key"`
	Url      string                      `json:"url"`
	FileSize uint64                      `json:"file_size"`
	Info     *object.AsyncFetchApiResult `json:"info"`
}

var _ flow.Work = (*asyncFetchResult)(nil)
var _ flow.Result = (*asyncFetchResult)(nil)

func (f *asyncFetchResult) String() string {
	return fmt.Sprintf("%s => [%s:%s]", f.Url, f.Bucket, f.Key)
}

func (f *asyncFetchResult) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", f.Bucket, f.Key, f.Url)
}

func (f *asyncFetchResult) IsValid() bool {
	return len(f.Bucket) > 0 && len(f.Key) > 0 && len(f.Url) > 0 && f.Info != nil
}
