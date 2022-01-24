package operations

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	MIN_UPLOAD_THREAD_COUNT = 20
	MAX_UPLOAD_THREAD_COUNT = 2000
)

type BatchUploadInfo struct {
	GroupInfo group.Info
}

// [qupload]命令， 上传本地文件到七牛存储中

// BatchUpload 该命令会读取配置文件， 上传本地文件系统的文件到七牛存储中;
// 可以设置多线程上传，默认的线程区间在[iqshell.MIN_UPLOAD_THREAD_COUNT, iqshell.MAX_UPLOAD_THREAD_COUNT]
func BatchUpload(info BatchUploadInfo) {
	uploadConfig := workspace.GetConfig().Up
	if err := uploadConfig.Check(); err != nil {
		log.ErrorF("batch upload:%v", err)
		return
	}

	//upload
	if info.GroupInfo.Info.WorkCount < MIN_UPLOAD_THREAD_COUNT {
		info.GroupInfo.Info.WorkCount = MIN_UPLOAD_THREAD_COUNT
		log.WarningF("Tip: you can set <ThreadCount> value between %d and %d to improve speed, and ThreadCount change to:%d",
			MIN_UPLOAD_THREAD_COUNT, MAX_UPLOAD_THREAD_COUNT, info.GroupInfo.Info.WorkCount)
	}

	if info.GroupInfo.Info.WorkCount > MAX_UPLOAD_THREAD_COUNT {
		info.GroupInfo.Info.WorkCount = MAX_UPLOAD_THREAD_COUNT
		log.WarningF("Tip: you can set <ThreadCount> value between %d and %d to improve speed, and ThreadCount change to:%d",
			MIN_UPLOAD_THREAD_COUNT, MAX_UPLOAD_THREAD_COUNT, info.GroupInfo.Info.WorkCount)
	}

	jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket, uploadConfig.FileList))
	dbPath := ""
	//if len(uploadConfig.RecordRoot) == 0 {
	dbPath = filepath.Join(workspace.GetWorkspace(), "upload", jobId, ".list")
	//} else {
	//	dbPath = filepath.Join(downloadCfg.RecordRoot, "download", jobId, ".list")
	//}
	log.InfoF("download db dir:%s", dbPath)

}

func batchUpload(info BatchUploadInfo, uploadConfig *config.Up) {
	handler, err := group.NewHandler(info.GroupInfo)
	if err != nil {
		log.Error(err)
		return
	}

	timeStart := time.Now()
	var locker sync.Mutex
	var totalFileCount int64 = 0 // ds.getFileLineCount()
	var currentFileCount int64
	var existsFileCount int64
	var updateFileCount int64
	var successFileCount int64
	var failureFileCount int64
	var skippedFileCount int64

	//hasPrefixes := len(uploadConfig.SkipPathPrefixes) > 0
	//prefixes := strings.Split(downloadCfg.Prefix, ",")
	//filterPrefix := func(name string) bool {
	//	if !hasPrefixes {
	//		return false
	//	}
	//
	//	for _, prefix := range prefixes {
	//		if strings.HasPrefix(name, prefix) {
	//			return false
	//		}
	//	}
	//	return true
	//}
	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, hasMore := handler.Scanner().ScanLine()
		if len(line) == 0 {
			return
		}

		items := strings.Split(line, info.GroupInfo.ItemSeparate)
		if len(items) < 3 {
			return nil, true
		}
		fileRelativePath := items[0]
		currentFileCount += 1

		//check skip local file or folder
		if skip, prefix := uploadConfig.HitByPathPrefixes(fileRelativePath); skip {
			log.InfoF("Skip by path prefix `%s` for local file path `%s`", prefix, fileRelativePath)
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
			return nil, true
		}

		if skip, prefix := uploadConfig.HitByFilePrefixes(fileRelativePath); skip {
			log.InfoF("Skip by file prefix `%s` for local file path `%s`", prefix, fileRelativePath)
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
			return nil, true
		}

		if skip, fixedStr := uploadConfig.HitByFixesString(fileRelativePath); skip {
			log.InfoF("Skip by fixed string `%s` for local file path `%s`", fixedStr, fileRelativePath)
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
			return nil, true
		}

		if skip, suffix := uploadConfig.HitBySuffixes(fileRelativePath); skip {
			log.InfoF("Skip by suffix `%s` for local file `%s`", suffix, fileRelativePath)
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
			return nil, true
		}

		apiInfo := &upload.ApiInfo{
			FilePath:         fileRelativePath,
			CheckExist:       false,
			CheckHash:        false,
			CheckSize:        false,
			OverWrite:        false,
			FileStatusDBPath: "",
			ToBucket:         "",
			SaveKey:          "",
			TokenProvider:    nil,
			TryTimes:         0,
			FileSize:         0,
			FileModifyTime:   0,
		}
		return apiInfo, hasMore
	}).DoWork(func(work work.Work) (work.Result, error) {
		locker.Lock()
		currentFileCount += 1
		locker.Unlock()

		apiInfo := work.(*upload.ApiInfo)
		res, err := uploadFile(*apiInfo)
		log.AlertF("Uploading %s [%d/%d, %.1f%%] ...", apiInfo.FilePath, currentFileCount, totalFileCount,
			float32(currentFileCount)*100/float32(totalFileCount))

		if res.IsSkip {
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
		}

		if err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(*download.ApiInfo)
		res := result.(download.ApiResult)
		handler.Export().Success().ExportF("download success, [%s:%s] => %s", apiInfo.Bucket, apiInfo.Key, res.FileAbsPath)

		locker.Lock()
		if res.IsExist {
			existsFileCount += 1
		} else if res.IsUpdate {
			updateFileCount += 1
		} else {
			successFileCount += 1
		}
		locker.Unlock()
	}).OnWorkError(func(work work.Work, err error) {
		locker.Lock()
		failureFileCount += 1
		locker.Unlock()

		apiInfo := work.(*upload.ApiInfo)
		handler.Export().Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* path fileSize fileModifyTime */
			apiInfo.FilePath, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()

	//skippedFileCount = totalFileCount - existsFileCount - successFileCount - updateFileCount - failureFileCount

	log.Alert("-------Upload Result-------")
	log.AlertF("%10s%10d", "Total:", totalFileCount)
	log.AlertF("%10s%10d", "Skipped:", skippedFileCount)
	log.AlertF("%10s%10d", "Exists:", existsFileCount)
	log.AlertF("%10s%10d", "Success:", successFileCount)
	log.AlertF("%10s%10d", "Update:", updateFileCount)
	log.AlertF("%10s%10d", "Failure:", failureFileCount)
	log.AlertF("%10s%15s", "Duration:", time.Since(timeStart))
	log.AlertF("-----------------------------")
	log.AlertF("See upload log at path:%s", uploadConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(data.STATUS_ERROR)
	}
}
