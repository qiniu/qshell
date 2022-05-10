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
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
	"time"
)

type FetchInfo object.FetchApiInfo

func (info *FetchInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.FromUrl) == 0 {
		return alert.CannotEmptyError("RemoteResourceUrl", "")
	}
	return nil
}

func Fetch(cfg *iqshell.Config, info FetchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	result, err := object.Fetch(object.FetchApiInfo(info))
	if err != nil {
		log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error:%v",
			info.FromUrl, info.Bucket, info.Key, err)
	} else {
		log.InfoF("Fetch Success, '%s' => [%s:%s]", info.FromUrl, info.Bucket, info.Key)
		log.AlertF("Key:%s", result.Key)
		log.AlertF("FileHash:%s", result.Hash)
		log.AlertF("Fsize: %d (%s)", result.Fsize, utils.FormatFileSize(result.Fsize))
		log.AlertF("Mime:%s", result.MimeType)
	}
}

type BatchFetchInfo struct {
	BatchInfo batch.Info
	Bucket    string
}

func (info *BatchFetchInfo) Check() *data.CodeError {
	if err := info.BatchInfo.Check(); err != nil {
		return err
	}

	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

//BatchFetch 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchFetch(cfg *iqshell.Config, info BatchFetchInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	exporter, err := export.NewFileExport(info.BatchInfo.FileExporterConfig)
	if err != nil {
		log.Error(err)
		return
	}

	flow.New(info.BatchInfo.Info).
		WorkProviderWithFile(info.BatchInfo.InputFile,
			info.BatchInfo.EnableStdin,
			flow.NewItemsWorkCreator(info.BatchInfo.ItemSeparate, 1, func(items []string) (work flow.Work, err *data.CodeError) {
				key := ""
				fromUrl := items[0]
				if len(items) > 1 {
					key = items[1]
				} else if k, e := utils.KeyFromUrl(fromUrl); e == nil {
					key = k
				}
				if len(key) == 0 || len(fromUrl) == 0 {
					return nil, alert.Error("key or fromUrl invalid", "")
				}

				return &object.FetchApiInfo{
					Bucket:  info.Bucket,
					Key:     key,
					FromUrl: fromUrl,
				}, nil
			})).
		WorkerProvider(flow.NewWorkerProvider(func() (flow.Worker, *data.CodeError) {
			return flow.NewSimpleWorker(func(workInfo *flow.WorkInfo) (flow.Result, *data.CodeError) {
				in := workInfo.Work.(*object.FetchApiInfo)
				return object.Fetch(*in)
			}), nil
		})).
		OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
			in, _ := workInfo.Work.(*object.FetchApiInfo)
			exporter.Success().ExportF("%s\t%s", in.FromUrl, in.Bucket)
			log.InfoF("Fetch Success, '%s' => [%s:%s]", in.FromUrl, info.Bucket, in.Key)
		}).
		OnWorkFail(func(workInfo *flow.WorkInfo, err *data.CodeError) {
			exporter.Fail().ExportF("%s%s%v", workInfo.Data, flow.ErrorSeparate, err)
			if in, ok := workInfo.Work.(*object.FetchApiInfo); ok {
				log.ErrorF("Fetch Failed, '%s' => [%s:%s], Error: %v", in.FromUrl, in.Bucket, in.Key, err)
			} else {
				log.ErrorF("Fetch Failed, %s, Error: %s", workInfo.Data, err)
			}
		}).Build().Start()
}

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

	fetchResultChan := make(chan fetchResult, 10)

	// fetch
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
			OnWorkSuccess(func(workInfo *flow.WorkInfo, result flow.Result) {
				in := workInfo.Work.(fetchItem)
				res := result.(object.AsyncFetchApiResult)
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
			exporter.Fail().ExportF("%s\t%d\t%s", result.url, result.fileSize, result.key)
			log.ErrorF("Fetch Failed, %s => [%s:%s], ID:%s", result.url, result.bucket, result.key, result.info.Id)
		}
	}
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
	info     object.AsyncFetchApiResult
}
