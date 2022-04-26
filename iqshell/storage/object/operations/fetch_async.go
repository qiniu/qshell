package operations

import (
	"fmt"
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
	"path/filepath"
	"strconv"
	"time"
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
	BatchInfo        batch.Info
	Bucket           string // fetch 的目的 bucket
	Host             string // 从指定URL下载时指定的HOST
	CallbackUrl      string // 抓取成功的回调地址
	CallbackBody     string
	CallbackBodyType string
	CallbackHost     string // 回调时使用的HOST
	FileType         int    // 文件存储类型， 0 标准存储， 1 低频存储
	Overwrite        bool
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

	// overseer， 数组源 不类型不记录中间状态
	var overseer flow.Overseer
	if info.BatchInfo.EnableRecord {
		dbPath := filepath.Join(workspace.GetJobDir(), ".recorder")
		log.DebugF("batch async fetch recorder:%s", dbPath)
		if overseer, err = flow.NewDBRecordOverseer(dbPath, func() *flow.WorkRecord {
			return &flow.WorkRecord{
				WorkInfo: &flow.WorkInfo{
					Data: "",
					Work: nil,
				},
				Result: &object.FetchResult{},
				Err:    nil,
			}
		}); err != nil {
			log.ErrorF("batch async fetch create overseer error:%v", err)
			return
		}
	} else {
		log.Debug("batch async fetch recorder:Not Enable")
	}

	fetchResultChan := make(chan fetchResult, info.BatchInfo.Info.WorkerCount*10)

	// fetch
	metric := &batch.Metric{}
	metric.Start()
	go func() {
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

					return fetchItem{
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
					in := workInfo.Work.(fetchItem)
					return object.AsyncFetch(in.info)
				}), nil
			})).
			SetOverseer(overseer).
			ShouldRedo(func(workInfo *flow.WorkInfo, workRecord *flow.WorkRecord) (shouldRedo bool, cause *data.CodeError) {
				if !info.BatchInfo.RecordRedoWhileError {
					return false, nil
				}

				if workRecord.Err != nil {
					return true, workRecord.Err
				}
				result, _ := workRecord.Result.(*object.FetchResult)
				if result == nil {
					return true, data.NewEmptyError().AppendDesc("no result found")
				}
				if result.Invalid() {
					return true, data.NewEmptyError().AppendDesc("result is invalid")
				}
				return false, nil
			}).
			OnWorkSkip(func(work *flow.WorkInfo, result flow.Result, err *data.CodeError) {
				metric.AddCurrentCount(1)
				metric.PrintProgress("Batching")

				operationResult, _ := result.(*object.AsyncFetchApiResult)
				if err != nil && err.Code == data.ErrorCodeAlreadyDone && operationResult != nil {
					if operationResult.Invalid() {
						metric.AddSuccessCount(1)
						log.DebugF("Skip line:%s because have done and success", work.Data)
					} else {
						metric.AddFailureCount(1)
						log.DebugF("Skip line:%s because have done and failure, %v", work.Data, err)
					}
				} else {
					metric.AddSkippedCount(1)
					log.DebugF("Skip line:%s because:%v", work.Data, err)
				}
			}).
			OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
				in := workInfo.Work.(fetchItem)
				res := result.(*object.AsyncFetchApiResult)
				fetchResultChan <- fetchResult{
					bucket:   in.info.Bucket,
					key:      in.info.Key,
					url:      in.info.Url,
					fileSize: in.fileSize,
					info:     res,
				}
				log.DebugF("Fetch Response, '%s' => [%s:%s] id:%s wait:%d", in.info.Url, in.info.Bucket, in.info.Key, res.Id, res.Wait)
			}).
			OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
				metric.AddFailureCount(1)

				if in, ok := workInfo.Work.(fetchItem); ok {
					exporter.Fail().ExportF("%s%s%v", in.info.Url, flow.ErrorSeparate, err)
					log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error: %v", in.info.Url, in.info.Bucket, in.info.Key, err)
				} else {
					exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
					log.ErrorF("Fetch Failed, %s, Error: %v", workInfo.Data, err)
				}
			}).Build().Start()

		close(fetchResultChan)
	}()

	// checker
	for result := range fetchResultChan {
		counter := 0
		maxDuration := asyncFetchCheckMaxDuration(result.fileSize)
		checkTime := time.Now().Add(maxDuration)
		for counter < 3 {
			current := time.Now()
			if counter == 0 || current.After(checkTime) {
				ret, cErr := object.CheckAsyncFetchStatus(result.bucket, result.info.Id)
				if cErr != nil {
					log.ErrorF("CheckAsyncFetchStatus: %v", cErr)
				} else if ret.Wait == -1 { // 视频抓取过一次，有可能成功了，有可能失败了
					counter += 1
					exist, _ := object.Exist(object.ExistApiInfo{
						Bucket: result.bucket,
						Key:    result.key,
					})
					if exist {
						// 标记成功
						metric.AddSuccessCount(1)

						exporter.Success().ExportF("%s\t%s", result.url, result.key)
						log.AlertF("Fetch Success, %s => [%s:%s]", result.url, result.bucket, result.key)
						break
					} else {
						log.ErrorF("Stat:%s error:%v ID:%s", result.key, err, result.info.Id)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
		if counter >= 3 {
			// 标记失败
			metric.AddFailureCount(1)

			exporter.Fail().ExportF("%s\t%d\t%s", result.url, result.fileSize, result.key)
			log.ErrorF("Fetch Failed, %s => [%s:%s], ID:%s", result.url, result.bucket, result.key, result.info.Id)
		}
	}

	metric.End()
	metric.PrintProgress("Batch async fetch")
}

func asyncFetchCheckMaxDuration(size uint64) time.Duration {
	duration := 3
	if size >= 500*utils.MB {
		duration = 40
	} else if size >= 200*utils.MB {
		duration = 30
	} else if size >= 100*utils.MB {
		duration = 20
	} else if size >= 50*utils.MB {
		duration = 10
	} else if size >= 10*utils.MB {
		duration = 6
	}
	return time.Duration(duration) * time.Second
}

type fetchItem struct {
	fileSize uint64
	info     object.AsyncFetchApiInfo
}

func (f fetchItem) WorkId() string {
	return fmt.Sprintf("%s:%s:%s", f.info.Url, f.info.Bucket, f.info.Key)
}

type fetchResult struct {
	bucket   string
	key      string
	url      string
	fileSize uint64
	info     *object.AsyncFetchApiResult
}
