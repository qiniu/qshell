package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"strconv"
	"strings"
	"time"
)

type ListInfo struct {
	Bucket       string
	Prefix       string
	Marker       string
	Delimiter    string
	Limit        int    //  最大输出条数，默认：-1, 无限输出
	StartDate    string // list item 的 put time 区间的开始时间 【闭区间】
	EndDate      string // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes     string // list item 必须包含后缀
	StorageTypes string // list item 存储类型，多个使用逗号隔开， 0:普通存储 1:低频存储 2:归档存储 3:深度归档存储
	MimeTypes    string // list item Mimetype类型，多个使用逗号隔开
	MinFileSize  string // 文件最小值，单位: B
	MaxFileSize  string // 文件最大值，单位: B
	MaxRetry     int    // -1: 无限重试
	SaveToFile   string
	AppendMode   bool
	Readable     bool
}

func (info *ListInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	return nil
}

func List(cfg *iqshell.Config, info ListInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	startTime, err := info.getStartDate()
	if err != nil {
		log.Error(err)
		return
	}
	endTime, err := info.getEndDate()
	if err != nil {
		log.Error(err)
		return
	}

	bucket.ListToFile(bucket.ListToFileApiInfo{
		ListApiInfo: bucket.ListApiInfo{
			Bucket:       info.Bucket,
			Prefix:       info.Prefix,
			Marker:       info.Marker,
			Delimiter:    info.Delimiter,
			Limit:        info.Limit,
			StartTime:    startTime,
			EndTime:      endTime,
			Suffixes:     info.getSuffixes(),
			StorageTypes: info.getStorageTypes(),
			MimeTypes:    info.getMimeTypes(),
			MinFileSize:  info.getMinFileSize(),
			MaxFileSize:  info.getMaxFileSize(),
			MaxRetry:     info.MaxRetry,
		},
		FilePath:   info.SaveToFile,
		AppendMode: info.AppendMode,
		Readable:   info.Readable,
	}, func(marker string, err *data.CodeError) {
		log.ErrorF("marker: %s", marker)
		log.ErrorF("list bucket Error: %v", err)
	})
}

func parseDate(dateString string) (time.Time, *data.CodeError) {
	if len(dateString) == 0 {
		return time.Time{}, nil
	}

	fields := strings.Split(dateString, "-")
	if len(fields) > 6 {
		return time.Time{}, data.NewEmptyError().AppendDescF("date format must be year-month-day-hour-minute-second")
	}

	var dateItems [6]int
	for ind, field := range fields {
		field, err := strconv.Atoi(field)
		if err != nil {
			return time.Time{}, data.NewEmptyError().AppendDescF("date format must be year-month-day-hour-minute-second, each field must be integer")
		}
		dateItems[ind] = field
	}
	return time.Date(dateItems[0], time.Month(dateItems[1]), dateItems[2], dateItems[3], dateItems[4], dateItems[5], 0, time.Local), nil
}

func (info ListInfo) getStartDate() (time.Time, *data.CodeError) {
	return parseDate(info.StartDate)
}

func (info ListInfo) getEndDate() (time.Time, *data.CodeError) {
	return parseDate(info.EndDate)
}

func (info ListInfo) getSuffixes() []string {
	ret := make([]string, 0)
	suffixes := strings.Split(info.Suffixes, ",")
	for _, s := range suffixes {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			ret = append(ret, strings.TrimSpace(s))
		}
	}
	return ret
}

func (info ListInfo) getStorageTypes() []int {
	ret := make([]int, 0)
	storageTypes := strings.Split(info.StorageTypes, ",")
	for _, s := range storageTypes {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			if storageType, e := strconv.Atoi(s); e == nil {
				ret = append(ret, storageType)
			} else {
				log.WarningF("storageType(%s) config error:%v", s, e)
			}
		}
	}
	return ret
}

func (info ListInfo) getMimeTypes() []string {
	ret := make([]string, 0)
	mimeTypes := strings.Split(info.MimeTypes, ",")
	for _, m := range mimeTypes {
		m = strings.TrimSpace(m)
		if len(m) > 0 {
			ret = append(ret, strings.TrimSpace(m))
		}
	}
	return ret
}

func (info ListInfo) getMinFileSize() int64 {
	if len(info.MinFileSize) == 0 {
		return -1
	}
	if size, e := strconv.ParseInt(info.MinFileSize, 10, 64); e == nil {
		return size
	} else {
		log.WarningF("minFileSize(%s) config error:%v", info.MinFileSize, e)
		return -1
	}
}

func (info ListInfo) getMaxFileSize() int64 {
	if len(info.MaxFileSize) == 0 {
		return -1
	}
	if size, e := strconv.ParseInt(info.MaxFileSize, 10, 64); e == nil {
		return size
	} else {
		log.WarningF("maxFileSize(%s) config error:%v", info.MaxFileSize, e)
		return -1
	}
}
