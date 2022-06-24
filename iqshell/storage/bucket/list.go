package bucket

import (
	"bufio"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

type ListApiInfo struct {
	Bucket       string
	Prefix       string
	Marker       string
	Delimiter    string
	Limit        int       //  最大输出条数，默认：-1, 无限输出
	StartTime    time.Time // list item 的 put time 区间的开始时间 【闭区间】
	EndTime      time.Time // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes     []string  // list item 必须包含后缀
	StorageTypes []int     // list item 存储类型，多个使用逗号隔开， 0:普通存储 1:低频存储 2:归档存储 3:深度归档存储
	MimeTypes    []string  // list item Mimetype类型，多个使用逗号隔开
	MinFileSize  int64     // 文件最小值，单位: B
	MaxFileSize  int64     // 文件最大值，单位: B
	MaxRetry     int       // -1: 无限重试
}

type ListObject storage.ListItem

// List list 某个 bucket 所有的文件
func List(info ListApiInfo,
	objectHandler func(marker string, object ListObject) (shouldContinue bool, err *data.CodeError),
	errorHandler func(marker string, err *data.CodeError)) {
	if objectHandler == nil {
		log.Error(alert.CannotEmpty("list bucket: object handler", ""))
		return
	}

	if errorHandler == nil {
		errorHandler = func(marker string, err *data.CodeError) {
			log.ErrorF("marker: %s", info.Marker)
			log.ErrorF("list bucket Error: %v", err)
		}
		log.Warning("list bucket: not set error handler")
	}

	bucketManager, err := GetBucketManager()
	if err != nil {
		errorHandler("", err)
		return
	}

	log.Debug("will list bucket")
	log.DebugF("Suffixes:%s", info.Suffixes)
	shouldCheckPutTime := !info.StartTime.IsZero() || !info.EndTime.IsZero()
	shouldCheckSuffixes := len(info.Suffixes) > 0
	shouldCheckStorageTypes := len(info.StorageTypes) > 0
	shouldCheckMimeTypes := len(info.MimeTypes) > 0
	shouldCheckFileSize := info.MinFileSize > 0 || info.MaxFileSize > 0
	retryCount := 0
	outputCount := 0
	complete := false
	for !complete && (info.MaxRetry < 0 || retryCount <= info.MaxRetry) {
		entries, lErr := bucketManager.ListBucketContext(workspace.GetContext(), info.Bucket, info.Prefix, info.Delimiter, info.Marker)
		if entries == nil && lErr == nil {
			lErr = errors.New("meet empty body when list not completed")
		}

		if lErr != nil {
			errorHandler(info.Marker, data.ConvertError(lErr))
			// 空间不存在，直接结束
			if strings.Contains(lErr.Error(), "query region error") ||
				strings.Contains(lErr.Error(), "incorrect zone") {
				break
			}

			retryCount++
			time.Sleep(1)
			continue
		}

		for listItem := range entries {
			if listItem.Marker != info.Marker {
				info.Marker = listItem.Marker
			}

			if listItem.Item.IsEmpty() {
				log.Debug("filter: item empty")
				continue
			}

			if shouldCheckPutTime {
				putTime := time.Unix(listItem.Item.PutTime/1e7, 0)
				if !filterByPutTime(putTime, info.StartTime, info.EndTime) {
					log.DebugF("filter %s: putTime not match, %s out of range [start:%s ~ end:%s]", listItem.Item.Key, putTime, info.StartTime, info.EndTime)
					continue
				}
			}

			if shouldCheckSuffixes && !filterBySuffixes(listItem.Item.Key, info.Suffixes) {
				log.DebugF("filter %s: key not match, key:%s suffixes:%s ", listItem.Item.Key, listItem.Item.Key, info.Suffixes)
				continue
			}

			if shouldCheckStorageTypes && !filterByStorageType(listItem.Item.Type, info.StorageTypes) {
				log.DebugF("filter %s: key not match, storageType:%d StorageTypes:%s ", listItem.Item.Key, listItem.Item.Type, info.Suffixes)
				continue
			}

			if shouldCheckMimeTypes && !filterByMimeType(listItem.Item.MimeType, info.MimeTypes) {
				log.DebugF("filter %s: key not match, mimeType:%s mimeTypes:%s ", listItem.Item.Key, listItem.Item.MimeType, info.MimeTypes)
				continue
			}

			if shouldCheckFileSize && !filterByFileSize(listItem.Item.Fsize, info.MinFileSize, info.MaxFileSize) {
				log.DebugF("filter %s: key not match, fileSize:%d minSize:%d maxSize:%d", listItem.Item.Key, listItem.Item.Fsize, info.MinFileSize, info.MaxFileSize)
				continue
			}

			shouldContinue, hErr := objectHandler(listItem.Marker, ListObject(listItem.Item))
			if hErr != nil {
				errorHandler(listItem.Marker, hErr)
			}
			if !shouldContinue {
				complete = true
				break
			}

			outputCount++
			if info.Limit > 0 && outputCount >= info.Limit {
				complete = true
				break
			}
		}

		if len(info.Marker) == 0 {
			// 列举结束
			break
		}

		retryCount = 0
	}

	if len(info.Marker) > 0 {
		log.InfoF("Marker: %s", info.Marker)
	}
	log.Debug("list bucket end")
}

type ListToFileApiInfo struct {
	ListApiInfo
	FilePath   string // file 不存在则输出到 stdout
	AppendMode bool
	Readable   bool
}

func ListToFile(info ListToFileApiInfo, errorHandler func(marker string, err *data.CodeError)) {
	if errorHandler == nil {
		errorHandler = func(marker string, err *data.CodeError) {
			log.ErrorF("marker: %s", marker)
			log.ErrorF("list bucket Error: %v", err)
		}
		log.Warning("list bucket to file: not set error handler")
	}

	var listResultFh io.WriteCloser
	if info.FilePath == "" {
		listResultFh = data.Stdout()
		log.Debug("prepare list bucket to stdout")
	} else {
		var openErr error
		var mode int

		if info.AppendMode {
			mode = os.O_APPEND | os.O_RDWR
		} else {
			mode = os.O_CREATE | os.O_RDWR | os.O_TRUNC
		}
		listResultFh, openErr = os.OpenFile(info.FilePath, mode, 0666)
		if openErr != nil {
			errorHandler("", data.NewEmptyError().AppendDescF("failed to open list result file `%s`, error:%v", info.FilePath, openErr))
			return
		}
		defer listResultFh.Close()
		log.Debug("prepare list bucket to file")
	}

	bWriter := bufio.NewWriter(listResultFh)
	if len(info.FilePath) == 0 {
		_, _ = bWriter.WriteString("Key\tFileSize\tHash\tPutTime\tMimeType\tStorageType\tEndUser\t\n")
		_ = bWriter.Flush()
	}
	List(info.ListApiInfo, func(marker string, object ListObject) (bool, *data.CodeError) {
		var fileSize interface{}
		if info.Readable {
			fileSize = utils.BytesToReadable(object.Fsize)
		} else {
			fileSize = object.Fsize
		}

		lineData := fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d\t%s\r\n",
			object.Key, fileSize, object.Hash,
			object.PutTime, object.MimeType, object.Type, object.EndUser)
		if _, wErr := bWriter.WriteString(lineData); wErr != nil {
			return false, data.NewEmptyError().AppendDesc("write error:" + wErr.Error())
		}

		if fErr := bWriter.Flush(); fErr != nil {
			return false, data.NewEmptyError().AppendDesc("flush error:" + fErr.Error())
		}
		return true, nil
	}, errorHandler)
}

func filterByPutTime(putTime, startDate, endDate time.Time) bool {
	switch {
	case startDate.IsZero() && endDate.IsZero():
		return true
	case !startDate.IsZero() && endDate.IsZero() && putTime.After(startDate):
		return true
	case !endDate.IsZero() && startDate.IsZero() && putTime.Before(endDate):
		return true
	case putTime.After(startDate) && putTime.Before(endDate):
		return true
	default:
		return false
	}
}

func filterBySuffixes(key string, suffixes []string) bool {
	hasSuffix := false
	if len(suffixes) == 0 {
		hasSuffix = true
	}
	for _, s := range suffixes {
		if strings.HasSuffix(key, s) {
			hasSuffix = true
			break
		}
	}
	return hasSuffix
}

func filterByStorageType(storageType int, storageTypes []int) bool {
	hasStorageType := false
	if len(storageTypes) == 0 {
		hasStorageType = true
	}
	for _, s := range storageTypes {
		if storageType == s {
			hasStorageType = true
			break
		}
	}
	return hasStorageType
}

func filterByMimeType(mimeType string, mimeTypes []string) bool {
	hasMimeType := false
	if len(mimeTypes) == 0 {
		hasMimeType = true
	}
	for _, s := range mimeTypes {
		if strings.Contains(s, "*") {
			sp := strings.ReplaceAll(s, "*", "")
			if strings.Contains(mimeType, sp) {
				hasMimeType = true
				break
			}
		} else if mimeType == s {
			hasMimeType = true
			break
		}
	}
	return hasMimeType
}

func filterByFileSize(fileSize, minSize, maxSize int64) bool {
	if maxSize < 0 {
		maxSize = math.MaxInt64
	}
	return fileSize >= minSize && fileSize <= maxSize
}
