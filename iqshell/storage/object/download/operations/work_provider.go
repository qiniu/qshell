package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"time"
)

func NewWorkProvider(bucket, keyPrefix, inputFile, itemSeparate string, infoResetHandler apiInfoResetHandler) flow.WorkProvider {
	provider := &workProvider{
		totalCount:       0,
		bucket:           bucket,
		keyPrefix:        keyPrefix,
		inputFile:        inputFile,
		itemSeparate:     itemSeparate,
		infoResetHandler: infoResetHandler,
		downloadItemChan: make(chan *downloadItem),
	}
	if len(inputFile) > 0 {
		provider.getWorkInfoFromFile()
	} else {
		provider.getWorkInfoFromBucket()
	}
	return provider
}

type apiInfoResetHandler func(apiInfo *download.DownloadActionInfo) *data.CodeError

type workProvider struct {
	totalCount       int64
	itemSeparate     string
	inputFile        string
	bucket           string
	keyPrefix        string
	infoResetHandler apiInfoResetHandler
	downloadItemChan chan *downloadItem
}

func (w *workProvider) WorkTotalCount() int64 {
	return w.totalCount
}

func (w *workProvider) Provide() (hasMore bool, workInfo *flow.WorkInfo, err *data.CodeError) {
	for item := range w.downloadItemChan {
		hasMore = true
		workInfo = item.workInfo
		err = item.err
		break
	}
	return
}

func (w *workProvider) getWorkInfoFromFile() {
	if len(w.inputFile) == 0 {
		return
	}

	w.totalCount = utils.GetFileLineCount(w.inputFile)

	go func() {
		lineParser := bucket.NewListLineParser()
		workPro, err := flow.NewWorkProviderOfFile(w.inputFile, false, flow.NewItemsWorkCreator(w.itemSeparate,
			1,
			func(items []string) (work flow.Work, err *data.CodeError) {
				listObject, e := lineParser.Parse(items)
				if e != nil {
					return nil, e
				}

				if len(listObject.Key) == 0 {
					return nil, alert.Error("key invalid", "")
				}

				info := &download.DownloadActionInfo{
					Key:               listObject.Key,
					ServerFileSize:    listObject.Fsize,
					ServerFileHash:    listObject.Hash,
					ServerFilePutTime: listObject.PutTime,
				}
				if w.infoResetHandler != nil {
					if e = w.infoResetHandler(info); e != nil {
						return nil, e
					}
				}
				return info, nil
			}))

		if err != nil {
			log.ErrorF("download create work provider error:%v", err)
			close(w.downloadItemChan)
			return
		}

		var keys []string
		for {
			if len(keys) == 10 {
				w.getWorkInfoOfKeys(keys)
				keys = nil
			}

			if keys == nil {
				keys = make([]string, 0, 10)
			}

			hasMore, workInfo, pErr := workPro.Provide()
			if pErr != nil {
				w.downloadItemChan <- &downloadItem{
					workInfo: workInfo,
					err:      pErr,
				}
			} else if workInfo != nil && workInfo.Work != nil {
				downloadApiInfo, _ := (workInfo.Work).(*download.DownloadActionInfo)
				if downloadApiInfo.ServerFilePutTime < 1 {
					keys = append(keys, downloadApiInfo.Key)
				} else {
					w.downloadItemChan <- &downloadItem{
						workInfo: workInfo,
						err:      pErr,
					}
				}
			}

			if !hasMore {
				w.getWorkInfoOfKeys(keys)
				keys = nil
				break
			}
		}

		close(w.downloadItemChan)
	}()
}

func (w *workProvider) getWorkInfoOfKeys(keys []string) {
	if len(keys) == 0 {
		return
	}

	operations := make([]batch.Operation, 0, len(keys))
	for _, key := range keys {
		if len(key) > 0 {
			operations = append(operations, object.StatusApiInfo{
				Bucket: w.bucket,
				Key:    key,
			})
		}
	}

	results, err := batch.Some(operations)
	if len(results) == len(operations) {
		for i, result := range results {
			item := operations[i].(object.StatusApiInfo)
			downItem := &downloadItem{}
			if result.Code != 200 || result.Error != "" {
				downItem.workInfo = &flow.WorkInfo{
					Data: item.Key,
				}
				downItem.err = data.NewError(result.Code, result.Error)
			} else {
				info := &download.DownloadActionInfo{
					Bucket:            w.bucket,
					Key:               item.Key,
					ServerFileHash:    result.Hash,
					ServerFileSize:    result.FSize,
					ServerFilePutTime: result.PutTime,
				}
				if w.infoResetHandler != nil {
					if e := w.infoResetHandler(info); e != nil {
						log.ErrorF("reset download api error:%v", e)
						continue
					}
				}
				downItem.workInfo = &flow.WorkInfo{
					Data: fmt.Sprintf("%s%s%d%s%s%s%d",
						item.Key, w.itemSeparate,
						result.FSize, w.itemSeparate,
						result.Hash, w.itemSeparate,
						result.PutTime),
					Work: info,
				}
			}
			w.downloadItemChan <- downItem
		}
	} else if err != nil {
		for _, operation := range operations {
			item := operation.(object.StatusApiInfo)
			w.downloadItemChan <- &downloadItem{
				workInfo: &flow.WorkInfo{
					Data: item.Key,
					Work: item,
				},
				err: err,
			}
		}
	}
}

func (w *workProvider) getWorkInfoFromBucket() {
	go func() {
		bucket.List(bucket.ListApiInfo{
			Bucket:    w.bucket,
			Prefix:    w.keyPrefix,
			Marker:    "",
			Delimiter: "",
			StartTime: time.Time{},
			EndTime:   time.Time{},
			Suffixes:  nil,
			MaxRetry:  20,
		}, func(marker string, object bucket.ListObject) (bool, *data.CodeError) {
			info := &download.DownloadActionInfo{
				Bucket:            w.bucket,
				Key:               object.Key,
				ServerFileHash:    object.Hash,
				ServerFileSize:    object.Fsize,
				ServerFilePutTime: object.PutTime,
			}
			if w.infoResetHandler != nil {
				if err := w.infoResetHandler(info); err != nil {
					return false, err
				}
			}

			w.downloadItemChan <- &downloadItem{
				workInfo: &flow.WorkInfo{
					Data: fmt.Sprintf("%s%s%d%s%s%s%d",
						object.Key, w.itemSeparate,
						object.Fsize, w.itemSeparate,
						object.Hash, w.itemSeparate,
						object.PutTime),
					Work: info,
				},
				err: nil,
			}
			return true, nil
		}, func(marker string, err *data.CodeError) {
			if err != nil {
				log.ErrorF("download list bucket error:%v", err)
			}
		})
		close(w.downloadItemChan)
	}()
}

type downloadItem struct {
	workInfo *flow.WorkInfo
	err      *data.CodeError
}
