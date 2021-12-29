package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
	"strconv"
	"sync"
	"time"
)

type FetchInfo rs.FetchApiInfo

func Fetch(info FetchInfo) {
	result, err := rs.Fetch(rs.FetchApiInfo(info))
	if err != nil {
		log.ErrorF("Fetch error: %v", err)
		os.Exit(data.STATUS_ERROR)
	} else {
		log.AlertF("Key:%s", result.Key)
		log.AlertF("Hash:%s", result.Hash)
		log.AlertF("Fsize: %d (%s)", result.Fsize, utils.FormatFileSize(result.Fsize))
		log.AlertF("Mime:%s", result.MimeType)
	}
}

type BatchFetchInfo struct {
	BatchInfo BatchInfo
	Bucket    string
}

//BatchFetch 批量删除，由于和批量删除的输入读取逻辑不同，所以分开
func BatchFetch(info BatchFetchInfo) {
	if !prepareToBatch(info.BatchInfo) {
		return
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:     info.BatchInfo.FailExportFilePath,
		OverrideExportFilePath: info.BatchInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	scanner, err := newBatchScanner(info.BatchInfo)
	if err != nil {
		log.ErrorF("get scanner error:%v", err)
		return
	}

	for {
		line, success := scanner.scanLine()
		if !success {
			break
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
		if key != "" && fromUrl != "" {
			_, err := rs.Fetch(rs.FetchApiInfo{
				Bucket:  info.Bucket,
				Key:     key,
				FromUrl: fromUrl,
			})

			if err != nil {
				resultExport.Fail().ExportF("%s\t%s\t%v", fromUrl, key, err)
				log.ErrorF("Fetch '%s' => %s:%s Failed, Error: %v", fromUrl, info.Bucket, key, err)
			} else {
				resultExport.Success().ExportF("%s\t%s", fromUrl, info.Bucket)
				log.AlertF("Fetch '%s' => %s:%s Success", fromUrl, info.Bucket, key)
			}
		}
	}
}

type CheckAsyncFetchStatusInfo struct {
	Bucket string
	Id     string
}

func CheckAsyncFetchStatus(info CheckAsyncFetchStatusInfo) {
	ret, err := rs.CheckAsyncFetchStatus(info.Bucket, info.Id)
	if err != nil {
		log.ErrorF("CheckAsyncFetchStatus: %v", err)
	} else {
		log.Alert(ret)
	}
}

type BatchAsyncFetchInfo struct {
	BatchInfo        BatchInfo
	Bucket           string // fetch 的目的 bucket
	Host             string // 从指定URL下载时指定的HOST
	Md5              string // 设置了该值，抓取的过程使用文件md5值进行校验, 校验失败不存在七牛空间
	Etag             string // 设置了该值， 抓取的过程中使用etag进行校验，失败不保存在存储空间中
	CallbackUrl      string // 抓取成功的回调地址
	CallbackBody     string
	CallbackBodyType string
	CallbackHost     string // 回调时使用的HOST
	FileType         int    // 文件存储类型， 0 标准存储， 1 低频存储
	InputFile        string // 输入访问地址列表
}

func BatchAsyncFetch(info BatchAsyncFetchInfo) {
	if !prepareToBatch(info.BatchInfo) {
		return
	}

	resultExport, err := export.NewFileExport(export.FileExporterConfig{
		SuccessExportFilePath:  info.BatchInfo.SuccessExportFilePath,
		FailExportFilePath:     info.BatchInfo.FailExportFilePath,
		OverrideExportFilePath: info.BatchInfo.OverrideExportFilePath,
	})
	if err != nil {
		log.ErrorF("get export error:%v", err)
		return
	}

	scanner, err := newBatchScanner(info.BatchInfo)
	if err != nil {
		log.ErrorF("get scanner error:%v", err)
		return
	}

	type fetchItem struct {
		fileSize uint64
		info     rs.AsyncFetchApiInfo
	}
	type fetchResult struct {
		bucket   string
		key      string
		url      string
		fileSize uint64
		info     rs.AsyncFetchApiResult
	}
	fetchInfoChan := make(chan fetchItem, info.BatchInfo.Worker)
	fetchResultChan := make(chan fetchResult)

	// 生产者
	go func() {
		for {
			line, success := scanner.scanLine()
			if !success {
				break
			}

			items := utils.SplitString(line, info.BatchInfo.ItemSeparate)
			if len(items) <= 0 {
				continue
			}

			var size uint64 = 0
			fromUrl := items[0]
			if len(items) >= 2 {
				s, pErr := strconv.ParseUint(items[1], 10, 64)
				if pErr != nil {
					resultExport.Fail().ExportF("%s: %v", line, pErr)
					continue
				}
				size = s
			}

			saveKey, pErr := utils.KeyFromUrl(fromUrl)
			if pErr != nil {
				resultExport.Fail().ExportF("%s: %v", line, pErr)
				continue
			}

			fetchInfoChan <- fetchItem{
				fileSize: size,
				info: rs.AsyncFetchApiInfo{
					Url:              fromUrl,
					Host:             info.Host,
					Bucket:           info.Bucket,
					Key:              saveKey,
					Md5:              info.Md5,
					Etag:             info.Etag,
					CallbackURL:      info.CallbackUrl,
					CallbackBody:     info.CallbackBody,
					CallbackBodyType: info.CallbackBodyType,
					FileType:         info.FileType,
				},
			}
		}
		close(fetchInfoChan)
	}()

	// 消费者
	fetchResultChanCloserOnce := &sync.Once{}
	fetchResultChanCloser := func() {
		fetchResultChanCloserOnce.Do(func() {
			close(fetchResultChan)
		})
	}
	for i := 0; i < info.BatchInfo.Worker; i++ {
		go func() {
			for info := range fetchInfoChan {
				ret, err := rs.AsyncFetch(info.info)
				if err != nil {
					resultExport.Fail().ExportF("%s: %v\n", info.info.Url, err)
					break
				}
				fetchResultChan <- fetchResult{
					bucket:   info.info.Bucket,
					key:      info.info.Key,
					url:      info.info.Url,
					fileSize: info.fileSize,
					info:     ret,
				}
			}
			fetchResultChanCloser()
		}()
	}

	// checker
	for result := range fetchResultChan {
		counter := 0
		maxDuration := asyncFetchCheckMaxDuration(result.fileSize)
		deadline := time.Now().Add(maxDuration)
		for counter < 3 {
			current := time.Now()
			if current.Before(deadline) {
				ret, cErr := rs.CheckAsyncFetchStatus(result.bucket, result.info.Id)
				if cErr != nil {
					log.ErrorF("CheckAsyncFetchStatus: %v", cErr)
				} else if ret.Wait == -1 { // 视频抓取过一次，有可能成功了，有可能失败了
					counter += 1
					exist, _ := bucket.CheckExists(result.bucket, result.key)
					if exist {
						resultExport.Success().ExportF("%s\t%s", result.url, result.key)
						log.Alert("fetch %s => %s:%s success", result.url, result.bucket, result.key)
						break
					} else {
						log.ErrorF("Stat: %s: %v", result.key, err)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
		if counter >= 3 {
			resultExport.Fail().ExportF("%s\t%d\t%s", result.url, result.fileSize, result.key)
			log.ErrorF("fetch %s => %s:%s failed", result.url, result.bucket, result.key)
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
