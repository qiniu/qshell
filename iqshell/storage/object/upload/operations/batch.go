package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/group"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/synchronized"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/work"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type BatchUploadInfo struct {
	// 输入文件通过 upload config 输入
	GroupInfo group.Info
}

func (info *BatchUploadInfo) Check() error {
	if info.GroupInfo.WorkCount < 1 || info.GroupInfo.WorkCount > 2000 {
		info.GroupInfo.WorkCount = 5
		log.WarningF("Tip: you can set <ThreadCount> value between 1 and 200 to improve speed, and now ThreadCount change to: %d",
			info.GroupInfo.Info.WorkCount)
	}
	if err := info.GroupInfo.Check(); err != nil {
		return err
	}
	return nil
}

// BatchUpload 该命令会读取配置文件， 上传本地文件系统的文件到七牛存储中;
// 可以设置多线程上传，默认的线程区间在[iqshell.min_upload_thread_count, iqshell.max_upload_thread_count]
func BatchUpload(cfg *iqshell.Config, info BatchUploadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if len(cfg.UploadConfigFile) == 0 {
		log.Error("LocalDownloadConfig can't empty")
		return
	}

	batchUpload(info)
}

func BatchUpload2(cfg *iqshell.Config, info BatchUploadInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	batchUpload(info)
}

func batchUpload(info BatchUploadInfo) {

	uploadConfig := workspace.GetConfig().Up
	if err := uploadConfig.Check(); err != nil {
		log.ErrorF("batch upload:%v", err)
		return
	}

	log.AlertF("Writing upload log to file:%s \n\n", uploadConfig.LogFile.Value())

	jobId := uploadConfig.JobId()
	cachePath := workspace.UploadCachePath()
	dbPath := filepath.Join(cachePath, jobId+".ldb")
	log.InfoF("upload status db file path:%s", dbPath)

	// 扫描本地文件
	needScanLocal := false

	if data.Empty(uploadConfig.FileList) {
		needScanLocal = true
	} else {
		if _, err := os.Stat(uploadConfig.FileList.Value()); err == nil {
			// 存在 file list 无需再重新扫描
			needScanLocal = false
			info.GroupInfo.InputFile = uploadConfig.FileList.Value()
		} else {
			info.GroupInfo.InputFile = filepath.Join(cachePath, jobId+".cache")
			if _, statErr := os.Stat(info.GroupInfo.InputFile); statErr == nil {
				//file exists
				needScanLocal = uploadConfig.IsRescanLocal()
			} else {
				needScanLocal = true
			}
		}
	}

	if needScanLocal {
		if uploadConfig.SrcDir != nil {
			log.ErrorF("scan error: src dir was empty")
			return
		}

		_, err := utils.DirCache(uploadConfig.SrcDir.Value(), info.GroupInfo.InputFile)
		if err != nil {
			log.ErrorF("create dir files cache error:%v", err)
			return
		}
	}

	batchUploadFlow(info, uploadConfig, dbPath)
}

func batchUploadFlow(info BatchUploadInfo, uploadConfig *config.Up, dbPath string) {
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
	syncLocker := synchronized.NewSynchronized(nil)
	var totalFileCount = handler.Scanner().LineCount()
	var currentFileCount int64
	var successFileCount int64
	var failureFileCount int64
	var notOverwriteCount int64
	var skippedFileCount int64
	work.NewFlowHandler(info.GroupInfo.Info).ReadWork(func() (work work.Work, hasMore bool) {
		line, hasMore := handler.Scanner().ScanLine()
		if !hasMore {
			return
		}

		if len(line) == 0 {
			syncLocker.Do(func() { skippedFileCount += 1 })
			log.InfoF("Skip by invalid line, items should more than 2:%s", line)
			return
		}

		items := strings.Split(line, info.GroupInfo.ItemSeparate)
		if len(items) < 3 {
			syncLocker.Do(func() { skippedFileCount += 1 })
			log.InfoF("Skip by invalid line, items should more than 2:%s", line)
			return nil, true
		}
		fileRelativePath := items[0]

		//check skip local file or folder
		if skip, prefix := uploadConfig.HitByPathPrefixes(fileRelativePath); skip {
			log.InfoF("Skip by path prefix `%s` for local file path `%s`", prefix, fileRelativePath)
			syncLocker.Do(func() { skippedFileCount += 1 })
			return nil, true
		}

		if skip, prefix := uploadConfig.HitByFilePrefixes(fileRelativePath); skip {
			log.InfoF("Skip by file prefix `%s` for local file path `%s`", prefix, fileRelativePath)
			syncLocker.Do(func() { skippedFileCount += 1 })
			return nil, true
		}

		if skip, fixedStr := uploadConfig.HitByFixesString(fileRelativePath); skip {
			log.InfoF("Skip by fixed string `%s` for local file path `%s`", fixedStr, fileRelativePath)
			syncLocker.Do(func() { skippedFileCount += 1 })
			return nil, true
		}

		if skip, suffix := uploadConfig.HitBySuffixes(fileRelativePath); skip {
			log.InfoF("Skip by suffix `%s` for local file `%s`", suffix, fileRelativePath)
			syncLocker.Do(func() { skippedFileCount += 1 })
			return nil, true
		}

		//pack the upload file key
		fileSize, _ := strconv.ParseInt(items[1], 10, 64)
		modifyTime, _ := strconv.ParseInt(items[2], 10, 64)
		key := fileRelativePath
		//check ignore dir
		if uploadConfig.IsIgnoreDir() {
			key = filepath.Base(key)
		}
		//check prefix
		if data.NotEmpty(uploadConfig.KeyPrefix) {
			key = strings.Join([]string{uploadConfig.KeyPrefix.Value(), key}, "")
		}
		//convert \ to / under windows
		if utils.IsWindowsOS() {
			key = strings.Replace(key, "\\", "/", -1)
		}
		//check file encoding
		if data.NotEmpty(uploadConfig.FileEncoding) && utils.IsGBKEncoding(uploadConfig.FileEncoding.Value()) {
			key, _ = utils.Gbk2Utf8(key)
		}

		localFilePath := filepath.Join(uploadConfig.SrcDir.Value(), fileRelativePath)
		apiInfo := &UploadInfo{
			FilePath:         localFilePath,
			Bucket:           uploadConfig.Bucket.Value(),
			Key:              key,
			MimeType:         "",
			FileStatusDBPath: dbPath,
			FileSize:         fileSize,
			FileModifyTime:   modifyTime,
			TokenProvider:    nil,
		}
		apiInfo.TokenProvider = createTokenProviderWithMac(mac, uploadConfig, apiInfo)
		return apiInfo, hasMore
	}).DoWork(func(work work.Work) (work.Result, error) {
		syncLocker.Do(func() {
			currentFileCount += 1
		})
		apiInfo := work.(*UploadInfo)

		log.AlertF("Uploading %s [%d/%d, %.1f%%] ...", apiInfo.FilePath, currentFileCount, totalFileCount,
			float32(currentFileCount)*100/float32(totalFileCount))

		res, err := uploadFile(*apiInfo)

		if err != nil {
			return nil, err
		}
		return res, nil
	}).OnWorkResult(func(work work.Work, result work.Result) {
		apiInfo := work.(*UploadInfo)
		res := result.(upload.ApiResult)
		if res.IsOverWrite {
			handler.Export().Override().ExportF("upload overwrite, %s => [%s:%s]", apiInfo.FilePath, apiInfo.Bucket, apiInfo.Key)
		} else if res.IsSkip {

		} else {
			handler.Export().Success().ExportF("upload success, %s => [%s:%s]", apiInfo.FilePath, apiInfo.Bucket, apiInfo.Key)
		}

		syncLocker.Do(func() {
			if res.IsNotOverWrite {
				notOverwriteCount += 1
			} else if res.IsSkip {
				skippedFileCount += 1
			} else {
				successFileCount += 1
			}
		})
	}).OnWorkError(func(work work.Work, err error) {
		syncLocker.Do(func() {
			failureFileCount += 1
		})

		apiInfo := work.(*UploadInfo)
		handler.Export().Fail().ExportF("%s%s%ld%s%s%s%ld%s error:%s", /* path fileSize fileModifyTime */
			apiInfo.FilePath, info.GroupInfo.ItemSeparate,
			apiInfo.FileSize, info.GroupInfo.ItemSeparate,
			apiInfo.FileModifyTime, info.GroupInfo.ItemSeparate,
			err)
	}).Start()

	log.Alert("--------------- Upload Result ---------------")
	log.AlertF("%20s%10d", "Total:", totalFileCount)
	log.AlertF("%20s%10d", "Success:", successFileCount)
	log.AlertF("%20s%10d", "Failure:", failureFileCount)
	log.AlertF("%20s%10d", "NotOverwrite:", notOverwriteCount)
	log.AlertF("%20s%10d", "Skipped:", skippedFileCount)
	log.AlertF("%20s%15s", "Duration:", time.Since(timeStart))
	log.AlertF("---------------------------------------------")
	log.AlertF("See upload log at path:%s \n\n", uploadConfig.LogFile.Value())
}

type BatchUploadConfigMouldInfo struct {
}

func BatchUploadConfigMould(cfg *iqshell.Config, info BatchUploadConfigMouldInfo) {
	log.Alert(uploadConfigMouldJsonString)
}
