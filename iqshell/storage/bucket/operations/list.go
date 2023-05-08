package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/internal/list"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ListInfo struct {
	Bucket             string // 指定空间【必选】
	Prefix             string // 指定前缀，只有资源名匹配该前缀的资源会被列出 【可选】
	Marker             string // 上一次列举返回的位置标记，作为本次列举的起点信息 【可选】
	Delimiter          string // 指定目录分隔符，列出所有公共前缀（模拟列出目录效果），默认值为空字符串【可选】
	StartDate          string // list item 的 put time 区间的开始时间 【闭区间】 【可选】
	EndDate            string // list item 的 put time 区间的终止时间 【闭区间】 【可选】
	Suffixes           string // list item 必须包含后缀 【可选】
	FileTypes          string // list item 存储类型，多个使用逗号隔开， 0:普通存储 1:低频存储 2:归档存储 3:深度归档存储 【可选】
	MimeTypes          string // list item Mimetype类型，多个使用逗号隔开 【可选】
	MinFileSize        string // 文件最小值，单位: B 【可选】
	MaxFileSize        string // 文件最大值，单位: B 【可选】
	MaxRetry           int    // -1: 无限重试 【可选】
	SaveToFile         string // 【可选】
	AppendMode         bool   // 【可选】
	Readable           bool   // 【可选】
	ShowFields         string // 需要展示的字段
	ApiVersion         string // list api 版本，v1 / v2【可选】
	ApiLimit           int    // 每次请求 size ，当前仅支持 list v1 【可选】
	OutputLimit        int    // 最大输出条数，默认：-1, 无限输出 【可选】
	OutputFieldsSep    string // 输出信息，每行的分隔符 【可选】
	OutputFileMaxLines int64  // 输出文件的最大行数，超过则自动创建新的文件，0：不限制输出文件的行数 【可选】
	OutputFileMaxSize  int64  // 输出文件的最大 Size，超过则自动创建新的文件，0：不限制输出文件的大小 【可选】
	EnableRecord       bool   // 是否开启 record 记录，开启后会记录 list 信息，下次 list 会自动指定 Marker 继续 list 【可选】
}

func (info *ListInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if len(info.ApiVersion) == 0 {
		info.ApiVersion = list.ApiVersionV1
	} else if info.ApiVersion == list.ApiVersionV2 {
		log.Warning("list bucket: api v2 is deprecated, recommend to use v1 instead")
	} else if info.ApiVersion != list.ApiVersionV1 && info.ApiVersion != list.ApiVersionV2 {
		return alert.Error("list bucket: api version is error, should set one of v1 and v2", "")
	}

	if len(info.ShowFields) > 0 {
		var fieldsNew []string
		fields := info.getShowFields()
		if len(fields) > 0 {
			for _, field := range fields {
				f := bucket.ListObjectField(field)
				if len(f) == 0 {
					return data.NewEmptyError().AppendDescF("show-fields value error:%s not support", field)
				}
				fieldsNew = append(fieldsNew, f)
			}
		}
		info.ShowFields = strings.Join(fieldsNew, ",")
	}

	if info.EnableRecord {
		// 记录模式开启 append
		info.AppendMode = true
	}

	if len(info.SaveToFile) > 0 {
		if absFilePath, aErr := filepath.Abs(info.SaveToFile); aErr != nil {
			log.ErrorF("=============:%v %s", aErr, info.SaveToFile)
			return data.ConvertError(aErr)
		} else {
			info.SaveToFile = absFilePath
		}
	}

	return nil
}

func (info *ListInfo) JobId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s", info.Prefix, info.Bucket, info.SaveToFile))
}

func (info *ListInfo) getShowFields() []string {
	if len(info.ShowFields) == 0 {
		return nil
	}
	return strings.Split(info.ShowFields, ",")
}

func (info *ListInfo) getOutputFieldsSep() string {
	if len(info.OutputFieldsSep) == 0 {
		return "\t"
	}
	return info.OutputFieldsSep
}

func List(cfg *iqshell.Config, info ListInfo) {
	cfg.JobPathBuilder = func(cmdPath string) string {
		return filepath.Join(cmdPath, info.JobId())
	}

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

	cacheDir := workspace.GetJobDir()
	log.InfoF("list status cache dir:%s ,you can delete it if you don't needed.", cacheDir)

	bucket.ListToFile(bucket.ListToFileApiInfo{
		ListApiInfo: bucket.ListApiInfo{
			Bucket:             info.Bucket,
			Prefix:             info.Prefix,
			Marker:             info.Marker,
			Delimiter:          info.Delimiter,
			StartTime:          startTime,
			EndTime:            endTime,
			Suffixes:           info.getSuffixes(),
			FileTypes:          info.getFileTypes(),
			MimeTypes:          info.getMimeTypes(),
			MinFileSize:        info.getMinFileSize(),
			MaxFileSize:        info.getMaxFileSize(),
			MaxRetry:           info.MaxRetry,
			ShowFields:         info.getShowFields(),
			ApiVersion:         info.ApiVersion,
			V1Limit:            info.ApiLimit,
			OutputLimit:        info.OutputLimit,
			OutputFieldsSep:    info.getOutputFieldsSep(),
			OutputFileMaxLines: info.OutputFileMaxLines,
			OutputFileMaxSize:  info.OutputFileMaxSize,
			CacheDir:           cacheDir,
			EnableRecord:       info.EnableRecord,
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

func (info *ListInfo) getStartDate() (time.Time, *data.CodeError) {
	return parseDate(info.StartDate)
}

func (info *ListInfo) getEndDate() (time.Time, *data.CodeError) {
	return parseDate(info.EndDate)
}

func (info *ListInfo) getSuffixes() []string {
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

func (info *ListInfo) getFileTypes() []int {
	ret := make([]int, 0)
	fileTypes := strings.Split(info.FileTypes, ",")
	for _, s := range fileTypes {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			if fileType, e := strconv.Atoi(s); e == nil {
				ret = append(ret, fileType)
			} else {
				log.WarningF("fileType(%s) config error:%v", s, e)
			}
		}
	}
	return ret
}

func (info *ListInfo) getMimeTypes() []string {
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

func (info *ListInfo) getMinFileSize() int64 {
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

func (info *ListInfo) getMaxFileSize() int64 {
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
