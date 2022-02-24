package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"strconv"
	"time"
)

type FetchInfo object.FetchApiInfo

func (info *FetchInfo) Check() error {
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
		log.ErrorF("Fetch error: %v", err)
	} else {
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

func (info *BatchFetchInfo) Check() error {
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

	handler, err := group.NewHandler(info.BatchInfo.Info)
	if err != nil {
		log.Error(err)
		return
	}

	work.NewFlowHandler(info.BatchInfo.Info.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, success := handler.Scanner().ScanLine()
		if !success {
			return nil, false
		}
		items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
		key, fromUrl := "", ""
		if len(items) > 0 {
			fromUrl = items[0]
			if len(items) > 1 {
				key = items[1]
			} else if k, err := utils.KeyFromUrl(fromUrl); err == nil {
				key = k
			}
		}
		if len(key) == 0 || len(fromUrl) == 0 {
			return nil, true
		}
		return object.FetchApiInfo{
			Bucket:  info.Bucket,
			Key:     key,
			FromUrl: fromUrl,
		}, true
	}).DoWork(func(work work.Work) (work.Result, error) {
		in := work.(object.FetchApiInfo)
		return object.Fetch(in)
	}).OnWorkResult(func(work work.Work, result work.Result) {
		in := work.(object.FetchApiInfo)
		handler.Export().Success().ExportF("%s\t%s", in.FromUrl, in.Bucket)
		log.InfoF("Fetch '%s' => %s:%s Success", in.FromUrl, info.Bucket, in.Key)
	}).OnWorkError(func(work work.Work, err error) {
		in := work.(object.FetchApiInfo)
		handler.Export().Fail().ExportF("%s\t%s\t%v", in.FromUrl, in.Key, err)
		log.ErrorF("Fetch '%s' => %s:%s Failed, Error: %v", in.FromUrl, in.Bucket, in.Key, err)
	}).Start()
}

type CheckAsyncFetchStatusInfo struct {
	Bucket string
	Id     string
}

func (info *CheckAsyncFetchStatusInfo) Check() error {
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
	GroupInfo        group.Info
	Bucket           string // fetch 的目的 bucket
	Host             string // 从指定URL下载时指定的HOST
	CallbackUrl      string // 抓取成功的回调地址
	CallbackBody     string
	CallbackBodyType string
	CallbackHost     string // 回调时使用的HOST
	FileType         int    // 文件存储类型， 0 标准存储， 1 低频存储
}

func (info *BatchAsyncFetchInfo) Check() error {
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

	info.GroupInfo.Force = true
	handler, err := group.NewHandler(info.GroupInfo)
	if err != nil {
		log.Error(err)
		return
	}

	type fetchItem struct {
		fileSize uint64
		info     object.AsyncFetchApiInfo
	}
	type fetchResult struct {
		bucket   string
		key      string
		url      string
		fileSize uint64
		info     object.AsyncFetchApiResult
	}
	fetchResultChan := make(chan fetchResult, 10)

	// fetch
	go func() {
		work.NewFlowHandler(info.GroupInfo.Info).
			ReadWork(func() (work work.Work, hasMore bool) {
				line, success := handler.Scanner().ScanLine()
				if !success {
					return nil, false
				}

				items := utils.SplitString(line, info.GroupInfo.ItemSeparate)
				if len(items) <= 0 {
					return nil, true
				}

				var size uint64 = 0
				fromUrl := items[0]
				if len(items) > 1 {
					s, pErr := strconv.ParseUint(items[1], 10, 64)
					if pErr != nil {
						handler.Export().Fail().ExportF("%s: %v", line, pErr)
						return nil, true
					}
					size = s
				}

				saveKey := ""
				if len(items) > 2 && len(items[2]) > 0 {
					saveKey = items[2]
				} else {
					key, pErr := utils.KeyFromUrl(fromUrl)
					if pErr != nil {
						handler.Export().Fail().ExportF("%s: %v", line, pErr)
						return nil, true
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
					},
				}, true
			}).
			DoWork(func(work work.Work) (work.Result, error) {
				in := work.(fetchItem)
				return object.AsyncFetch(in.info)
			}).
			OnWorkResult(func(work work.Work, result work.Result) {
				in := work.(fetchItem)
				res := result.(object.AsyncFetchApiResult)
				fetchResultChan <- fetchResult{
					bucket:   in.info.Bucket,
					key:      in.info.Key,
					url:      in.info.Url,
					fileSize: in.fileSize,
					info:     res,
				}
				log.DebugF("Fetch Response '%s' => %s:%s id:%s wait:%d", in.info.Url, in.info.Bucket, in.info.Key, res.Id, res.Wait)
			}).
			OnWorkError(func(work work.Work, err error) {
				in := work.(fetchItem)
				handler.Export().Fail().ExportF("%s: %v", in.info.Url, err)
				log.ErrorF("Fetch '%s' => %s:%s Failed, Error: %v", in.info.Url, in.info.Bucket, in.info.Key, err)
			}).
			OnWorksComplete(func() {
				close(fetchResultChan)
			}).
			Start()
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
						handler.Export().Success().ExportF("%s\t%s", result.url, result.key)
						log.AlertF("fetch %s => %s:%s success", result.url, result.bucket, result.key)
						break
					} else {
						log.ErrorF("Stat:%s error:%v ID:%s", result.key, err, result.info.Id)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
		if counter >= 3 {
			handler.Export().Fail().ExportF("%s\t%d\t%s", result.url, result.fileSize, result.key)
			log.ErrorF("fetch %s => %s:%s, ID:%s failed", result.url, result.bucket, result.key, result.info.Id)
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
