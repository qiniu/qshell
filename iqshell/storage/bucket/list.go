package bucket

import (
	"bufio"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/file"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/internal/list"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ListApiInfo struct {
	Bucket             string    // 空间名	【必选】
	Prefix             string    // 前缀
	Marker             string    // 标记
	Delimiter          string    //
	StartTime          time.Time // list item 的 put time 区间的开始时间 【闭区间】
	EndTime            time.Time // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes           []string  // list item 必须包含后缀
	FileTypes          []int     // list item 存储类型，多个使用逗号隔开， 0:普通存储 1:低频存储 2:归档存储 3:深度归档存储
	MimeTypes          []string  // list item Mimetype类型，多个使用逗号隔开
	MinFileSize        int64     // 文件最小值，单位: B
	MaxFileSize        int64     // 文件最大值，单位: B
	MaxRetry           int       // -1: 无限重试
	ShowFields         []string  // 需要展示的字段  【必选】
	ApiVersion         string    // list api 版本，v1 / v2【可选】
	V1Limit            int       // 每次请求 size ，list v1 特有
	OutputLimit        int       // 最大输出条数，默认：-1, 无限输出
	OutputFieldsSep    string    // 输出信息，每行的分隔符 【必选】
	OutputFileMaxLines int64     // 输出文件的最大行数，超过则自动创建新的文件，0：不限制输出文件的行数 【可选】
	OutputFileMaxSize  int64     // 输出文件的最大 Size，超过则自动创建新的文件，0：不限制输出文件的大小 【可选】
	EnableRecord       bool      // 是否开启 record 记录，开启后会记录 list 信息，下次 list 会自动指定 Marker 继续 list 【可选】
	CacheDir           string    // 历史数据存储路径 【内部使用】
}

func ListObjectField(field string) string {
	for _, f := range listObjectFields {
		if strings.EqualFold(field, f) {
			return f
		}
	}
	return ""
}

type ListObject = list.Item

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
	shouldCheckFileTypes := len(info.FileTypes) > 0
	shouldCheckMimeTypes := len(info.MimeTypes) > 0
	shouldCheckFileSize := info.MinFileSize > 0 || info.MaxFileSize > 0
	isItemExcepted := func(listItem list.Item) (isExcepted bool) {
		if shouldCheckPutTime {
			putTime := time.Unix(listItem.PutTime/1e7, 0)
			if !filterByPutTime(putTime, info.StartTime, info.EndTime) {
				log.DebugF("filter %s: putTime not match, %s out of range [start:%s ~ end:%s]", listItem.Key, putTime, info.StartTime, info.EndTime)
				return false
			}
		}

		if shouldCheckSuffixes && !filterBySuffixes(listItem.Key, info.Suffixes) {
			log.DebugF("filter %s: key not match, key:%s suffixes:%s ", listItem.Key, listItem.Key, info.Suffixes)
			return false
		}

		if shouldCheckFileTypes && !filterByFileType(listItem.Type, info.FileTypes) {
			log.DebugF("filter %s: key not match, fileType:%d FileTypes:%s ", listItem.Key, listItem.Type, info.Suffixes)
			return false
		}

		if shouldCheckMimeTypes && !filterByMimeType(listItem.MimeType, info.MimeTypes) {
			log.DebugF("filter %s: key not match, mimeType:%s mimeTypes:%s ", listItem.Key, listItem.MimeType, info.MimeTypes)
			return false
		}

		if shouldCheckFileSize && !filterByFileSize(listItem.Fsize, info.MinFileSize, info.MaxFileSize) {
			log.DebugF("filter %s: key not match, fileSize:%d minSize:%d maxSize:%d", listItem.Key, listItem.Fsize, info.MinFileSize, info.MaxFileSize)
			return false
		}

		return true
	}

	listWaiter := sync.WaitGroup{}
	listWaiter.Add(1)
	workspace.AddCancelObserver(func(s os.Signal) {
		listWaiter.Wait()
	})

	cache := &listCache{
		enableRecord: info.EnableRecord,
		cachePath:    filepath.Join(info.CacheDir, "info.json"),
	}
	cacheInfoP, err := cache.loadCache()
	if err != nil {
		log.Debug(err)
	}

	if cacheInfoP != nil && len(cacheInfoP.Marker) > 0 {
		if len(info.Marker) == 0 {
			info.Marker = cacheInfoP.Marker
		}
	}

	if cacheInfoP == nil {
		cacheInfoP = &cacheInfo{}
	} else if len(cacheInfoP.Marker) > 0 {
		log.InfoF("use marker:%s", cacheInfoP.Marker)
	}

	retryCount := 0
	outputCount := 0
	complete := false
	for !complete && (info.MaxRetry < 0 || retryCount <= info.MaxRetry) {
		var hasMore = false
		var lErr *data.CodeError = nil

		if !workspace.IsCmdInterrupt() {
			hasMore, lErr = list.ListBucket(workspace.GetContext(), list.ApiInfo{
				Manager:    bucketManager,
				ApiVersion: list.ApiVersion(info.ApiVersion),
				Bucket:     info.Bucket,
				Prefix:     info.Prefix,
				Delimiter:  info.Delimiter,
				Marker:     info.Marker,
				V1Limit:    info.V1Limit,
			}, func(marker string, dir string, listItem list.Item) (stop bool) {
				if marker != info.Marker {
					info.Marker = marker
				}

				if listItem.IsNull() {
					return false
				}

				if !isItemExcepted(listItem) {
					return false
				}

				shouldContinue, hErr := objectHandler(marker, listItem)
				if hErr != nil {
					errorHandler(marker, hErr)
				}
				if !shouldContinue {
					complete = true
					return true
				}

				outputCount++
				if info.OutputLimit > 0 && outputCount >= info.OutputLimit {
					complete = true
					return true
				}

				return false
			})
		}

		// 保存信息
		cacheInfoP.Bucket = info.Bucket
		cacheInfoP.Prefix = info.Prefix
		cacheInfoP.Marker = info.Marker
		_ = cache.saveCache(cacheInfoP)

		if workspace.IsCmdInterrupt() && lErr == nil {
			lErr = data.NewError(0, "list is interrupted")
		}

		if lErr != nil || workspace.IsCmdInterrupt() {
			errorHandler(info.Marker, lErr)

			if workspace.IsCmdInterrupt() || // 取消
				strings.Contains(lErr.Error(), "no such bucket") || // 空间不存在，直接结束
				strings.Contains(lErr.Error(), "incorrect zone") || // 空间不正确
				strings.Contains(lErr.Error(), "query region error") || // 查询空间错误
				strings.Contains(lErr.Error(), "app/accesskey is not found") || // AK/SK 错误
				strings.Contains(lErr.Error(), "invalid list limit") || //  api v1 list limit
				strings.Contains(lErr.Error(), "context canceled") { // 取消
				break
			}

			retryCount++
			time.Sleep(time.Millisecond * 500)
			continue
		}

		if !hasMore || complete || workspace.IsCmdInterrupt() {
			break
		}

		retryCount = 0
	}

	if len(info.Marker) == 0 {
		if rErr := cache.removeCache(); rErr != nil {
			log.ErrorF("list remove cache status error: %v", rErr)
		} else {
			log.InfoF("list success, remove cache status: %s", cache.cachePath)
		}
	} else {
		log.InfoF("Marker: %s", info.Marker)
	}

	log.Debug("list bucket end")

	listWaiter.Done()
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

	if len(info.ShowFields) == 0 {
		info.ShowFields = listObjectFields
	}

	if len(info.OutputFieldsSep) == 0 {
		info.OutputFieldsSep = data.DefaultLineSeparate
	}

	// 文件头
	title := strings.Join(info.ShowFields, info.OutputFieldsSep)

	var output io.WriteCloser
	if info.FilePath == "" {
		output = data.Stdout()
		_, _ = output.Write([]byte(title + "\n"))
		log.Debug("prepare list bucket to stdout")
	} else {
		var nErr *data.CodeError
		output, nErr = file.NewRotateFile(info.FilePath,
			file.RotateOptionMaxSize(info.OutputFileMaxSize),
			file.RotateOptionMaxLine(info.OutputFileMaxLines),
			file.RotateOptionAppendMode(info.AppendMode),
			file.RotateOptionFileHeader(title),
			file.RotateOptionOnOpenFile(func(filename string) {
				log.InfoF("open output file and prepare to write:%v", filename)
			}))

		if nErr != nil {
			errorHandler("", data.NewEmptyError().AppendDescF("failed to create rotate file:`%s`, error:%v", info.FilePath, nErr))
			return
		}
		defer output.Close()
		log.Debug("prepare list bucket to file")
	}

	bWriter := bufio.NewWriter(output)
	lineCreator := &ListLineCreator{
		Fields:   info.ShowFields,
		Sep:      info.OutputFieldsSep,
		Readable: info.Readable,
	}
	List(info.ListApiInfo, func(marker string, object ListObject) (bool, *data.CodeError) {
		lineData := lineCreator.Create(&object)
		if _, wErr := bWriter.WriteString(lineData + "\n"); wErr != nil {
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

func filterByFileType(fileType int, fileTypes []int) bool {
	hasFileType := false
	if len(fileTypes) == 0 {
		hasFileType = true
	}
	for _, s := range fileTypes {
		if fileType == s {
			hasFileType = true
			break
		}
	}
	return hasFileType
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
