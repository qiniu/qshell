package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"strconv"
	"strings"
	"time"
)

type ListInfo struct {
	Bucket     string
	Prefix     string
	Marker     string
	Delimiter  string
	StartDate  string // list item 的 put time 区间的开始时间 【闭区间】
	EndDate    string // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes   string // list item 必须包含前缀
	MaxRetry   int    // -1: 无限重试
	SaveToFile string
	AppendMode bool
	Readable   bool
}

func List(info ListInfo) {
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
			Bucket:            info.Bucket,
			Prefix:            info.Prefix,
			Marker:            info.Marker,
			Delimiter:         info.Delimiter,
			StartTime:         startTime,
			EndTime:           endTime,
			Suffixes:          info.getSuffixes(),
			MaxRetry:          info.MaxRetry,
			StopWhenListError: false,
		},
		FilePath:    info.SaveToFile,
		AppendMode:  info.AppendMode,
		Readable:    info.Readable,
	}, func(marker string, err error) {
		log.ErrorF("marker: %s", marker)
		log.ErrorF("list bucket Error: %v", err)
	})
}

func parseDate(dateString string) (time.Time, error) {
	if len(dateString) == 0 {
		return time.Time{}, nil
	}

	fields := strings.Split(dateString, "-")
	if len(fields) > 6 {
		return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second")
	}

	var dateItems [6]int
	for ind, field := range fields {
		field, err := strconv.Atoi(field)
		if err != nil {
			return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second, each field must be integer")
		}
		dateItems[ind] = field
	}
	return time.Date(dateItems[0], time.Month(dateItems[1]), dateItems[2], dateItems[3], dateItems[4], dateItems[5], 0, time.Local), nil
}

func (info ListInfo) getStartDate() (time.Time, error) {
	return parseDate(info.StartDate)
}

func (info ListInfo) getEndDate() (time.Time, error) {
	return parseDate(info.EndDate)
}

func (info ListInfo) getSuffixes() []string {
	sf := make([]string, 0)
	for _, s := range strings.Split(info.Suffixes, ",") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			sf = append(sf, strings.TrimSpace(s))
		}
	}
	return nil
}
