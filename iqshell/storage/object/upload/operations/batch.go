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
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MIN_UPLOAD_THREAD_COUNT = 20
	MAX_UPLOAD_THREAD_COUNT = 2000
)

type BatchUploadInfo struct {
	// 输入文件通过 upload config 输入
	GroupInfo group.Info
}

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

	cachePath := uploadConfig.RecordRoot
	if len(cachePath) == 0 {
		cachePath = workspace.GetWorkspace()
	}
	jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket, uploadConfig.FileList))
	cachePath = filepath.Join(cachePath, jobId, "qupload")

	dbPath := filepath.Join(cachePath, jobId, ".ldb")
	log.InfoF("upload status db file path:%s", dbPath)

	needRescanLocal := uploadConfig.RescanLocal
	_, localFileStatErr := os.Stat(uploadConfig.FileList)
	if uploadConfig.FileList != "" && localFileStatErr == nil {
		info.GroupInfo.InputFile = uploadConfig.FileList
	} else {
		info.GroupInfo.InputFile = filepath.Join(cachePath, jobId, ".cache")
		needRescanLocal = true
	}
	if needRescanLocal {
		_, err := utils.DirCache(uploadConfig.SrcDir, info.GroupInfo.InputFile)
		if err != nil {
			log.ErrorF("create dir files cache error:%v", err)
			return
		}
	}

	batchUpload(info, uploadConfig, dbPath)
}

func batchUpload(info BatchUploadInfo, uploadConfig *config.Up, dbPath string) {
	handler, err := group.NewHandler(info.GroupInfo)
	if err != nil {
		log.Error(err)
		return
	}

	mac, err := workspace.GetMac()
	if err != nil {
		log.Error("get mac error:" + err.Error())
		return
	}

	timeStart := time.Now()
	var locker sync.Mutex
	var totalFileCount = handler.Scanner().LineCount()
	var currentFileCount int64
	var successFileCount int64
	var failureFileCount int64
	var notOverwriteCount int64
	var skippedFileCount int64
	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, hasMore := handler.Scanner().ScanLine()
		if len(line) == 0 {
			return
		}

		items := strings.Split(line, info.GroupInfo.ItemSeparate)
		if len(items) < 3 {
			locker.Lock()
			skippedFileCount += 1
			locker.Unlock()
			log.InfoF("Skip by invalid line, items should more than 2:%s", line)
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

		//pack the upload file key
		fileSize, _ := strconv.ParseInt(items[1], 10, 64)
		modifyTime, _ := strconv.ParseInt(items[2], 10, 64)
		key := fileRelativePath
		//check ignore dir
		if uploadConfig.IgnoreDir {
			key = filepath.Base(key)
		}
		//check prefix
		if uploadConfig.KeyPrefix != "" {
			key = strings.Join([]string{uploadConfig.KeyPrefix, key}, "")
		}
		//convert \ to / under windows
		if utils.IsWindowsOS() {
			key = strings.Replace(key, "\\", "/", -1)
		}
		//check file encoding
		if utils.IsGBKEncoding(uploadConfig.FileEncoding) {
			key, _ = utils.Gbk2Utf8(key)
		}

		localFilePath := filepath.Join(uploadConfig.SrcDir, fileRelativePath)
		apiInfo := &upload.ApiInfo{
			FilePath:         localFilePath,
			CheckExist:       uploadConfig.CheckExists,
			CheckHash:        uploadConfig.CheckHash,
			CheckSize:        uploadConfig.CheckSize,
			Overwrite:        uploadConfig.Overwrite,
			FileStatusDBPath: dbPath,
			ToBucket:         uploadConfig.Bucket,
			SaveKey:          key,
			TokenProvider:    createTokenProviderWithConfig(mac, uploadConfig, key),
			TryTimes:         0,
			FileSize:         fileSize,
			FileModifyTime:   modifyTime,
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

		if err != nil {
			return nil, err
		}
		return res, nil
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(*upload.ApiInfo)
		res := result.(upload.ApiResult)
		handler.Export().Success().ExportF("upload success, %s => [%s:%s]", apiInfo.FilePath, apiInfo.ToBucket, apiInfo.SaveKey)

		locker.Lock()
		if res.IsNotOverWrite {
			notOverwriteCount += 1
		} else if res.IsSkip {
			skippedFileCount += 1
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

	log.Alert("-------Upload ApiResult-------")
	log.AlertF("%20s%10d", "Total:", totalFileCount)
	log.AlertF("%20s%10d", "Success:", successFileCount)
	log.AlertF("%20s%10d", "Failure:", failureFileCount)
	log.AlertF("%20s%10d", "NotOverwrite:", notOverwriteCount)
	log.AlertF("%20s%10d", "Skipped:", skippedFileCount)
	log.AlertF("%20s%15s", "Duration:", time.Since(timeStart))
	log.AlertF("-----------------------------")
	log.AlertF("See upload log at path:%s", uploadConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(data.STATUS_ERROR)
	}
}
