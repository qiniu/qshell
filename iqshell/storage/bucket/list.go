package bucket

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"strings"
	"time"
)

type ListApiInfo struct {
	Bucket            string
	Prefix            string
	Marker            string
	Delimiter         string
	StartTime         time.Time // list item 的 put time 区间的开始时间 【闭区间】
	EndTime           time.Time // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes          []string  // list item 必须包含前缀
	MaxRetry          int       // -1: 无限重试
	StopWhenListError bool      // 当 list 过程中出现错误是否停止 list
}

type ListObject storage.ListItem

// List list 某个 bucket 所有的文件
func List(info ListApiInfo, objectHandler func(marker string, object ListObject) error, errorHandler func(marker string, err error)) {
	if objectHandler == nil {
		logs.Error(alert.CannotEmpty("list bucket: object handler", ""))
		return
	}

	if errorHandler == nil {
		errorHandler = func(marker string, err error) {
			log.ErrorF("marker: %s", info.Marker)
			log.ErrorF("list bucket Error: %v", err)
		}
		logs.Warning("list bucket: not set error handler")
	}

	bucketManager, err := GetBucketManager()
	if err != nil {
		errorHandler("", err)
		return
	}

	shouldCheckPutTime := !info.StartTime.IsZero() || !info.StartTime.IsZero()
	shouldCheckSuffixes := len(info.Suffixes) > 0
	retryCount := 0
	complete := false
	for ;!complete && (info.MaxRetry < 0 || retryCount <= info.MaxRetry); {
		entries, lErr := bucketManager.ListBucketContext(workspace.GetContext(), info.Bucket, info.Prefix, info.Delimiter, info.Marker)
		if entries == nil && lErr == nil {
			// no data
			if info.Marker == "" {
				complete = true
				break
			} else {
				lErr = errors.New("meet empty body when list not completed")
			}
		}

		if lErr != nil {
			errorHandler(info.Marker, lErr)
			retryCount++
			time.Sleep(1)
			continue
		}

		for listItem := range entries {
			if listItem.Marker != info.Marker {
				info.Marker = listItem.Marker
			}

			if listItem.Item.IsEmpty() {
				continue
			}

			if shouldCheckPutTime {
				putTime := time.Unix(listItem.Item.PutTime/1e7, 0)
				if !filterByPutTime(putTime, info.StartTime, info.EndTime) {
					continue
				}
			}

			if shouldCheckSuffixes && !filterBySuffixes(listItem.Item.Key, info.Suffixes) {
				continue
			}

			hErr := objectHandler(listItem.Marker, ListObject(listItem.Item))
			if hErr != nil {
				errorHandler(listItem.Marker, hErr)
			}
		}

		retryCount = 0
	}

	if len(info.Marker) > 0 {
		log.InfoF("Marker: %s", info.Marker)
	}
}

type ListToFileApiInfo struct {
	ListApiInfo
	FilePath   string // file 不存在则输出到 stdout
	AppendMode bool
	Readable   bool
}

func ListToFile(info ListToFileApiInfo, errorHandler func(marker string, err error)) {
	if errorHandler == nil {
		errorHandler = func(marker string, err error) {
			log.ErrorF("marker: %s", marker)
			log.ErrorF("list bucket Error: %v", err)
		}
		logs.Warning("list bucket to file: not set error handler")
	}

	var listResultFh *os.File
	if info.FilePath == "" {
		listResultFh = os.Stdout
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
			errorHandler("", fmt.Errorf("failed to open list result file `%s`, error:%v", info.FilePath, openErr))
			return
		}
		defer listResultFh.Close()
	}

	bWriter := bufio.NewWriter(listResultFh)
	List(info.ListApiInfo, func(marker string, object ListObject) error {
		var fileSize interface{}
		if info.Readable {
			fileSize = utils.BytesToReadable(object.Fsize)
		} else {
			fileSize = object.Fsize
		}

		lineData := fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d\t%s\r\n",
			object.Key, fileSize, object.Hash,
			object.PutTime, object.MimeType, object.Type, object.EndUser)
		_, wErr := bWriter.WriteString(lineData)
		if wErr == nil {
			return nil
		}
		return errors.New("flush error:" + wErr.Error())
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
	if hasSuffix {
		return true
	} else {
		return false
	}
}
