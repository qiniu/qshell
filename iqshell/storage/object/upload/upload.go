package upload

import (
	"bufio"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/export"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_PUT_THRESHOLD   int64 = 10 * 1024 * 1024 //10MB
	MIN_UPLOAD_THREAD_COUNT int64 = 1
	MAX_UPLOAD_THREAD_COUNT int64 = 2000
)

var (
	currentFileCount  int64
	successFileCount  int64
	notOverwriteCount int64
	failureFileCount  int64
	skippedFileCount  int64

	uploadTasks chan func()
	initUpOnce  sync.Once
)

func doUpload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

// QiniuUpload
func QiniuUpload(threadCount int, uploadConfig *config.UploadConfig, exporter *export.FileExporter) {
	var upSettings = storage.Settings{
		TaskQsize: 8,
		Workers:   4,
		ChunkSize: 4 * 1024 * 1024,
		TryTimes:  3,
	}
	timeStart := time.Now()
	//create job id
	jobId := uploadConfig.JobId()
	QShellRootPath := workspace.GetWorkspace()
	if QShellRootPath == "" {
		log.Error("Empty root path")
		os.Exit(data.STATUS_HALT)
	}
	storePath := filepath.Join(QShellRootPath, "qupload", jobId)

	uploadConfig.PrepareLogger(storePath, jobId)

	//chunk upload threshold
	putThreshold := DEFAULT_PUT_THRESHOLD
	if uploadConfig.PutThreshold > 0 {
		putThreshold = uploadConfig.PutThreshold
	}

	//set resume upload settings
	storage.SetSettings(&upSettings)

	//make SrcDir the full path
	uploadConfig.SrcDir, _ = filepath.Abs(uploadConfig.SrcDir)

	cacheResultName, totalFileCount, cErr := uploadConfig.CacheFileNameAndCount(storePath, jobId)
	if cErr != nil {
		os.Exit(data.STATUS_HALT)
	}

	//leveldb folder
	leveldbFileName := filepath.Join(storePath, jobId+".ldb")
	ldb, err := leveldb.OpenFile(leveldbFileName, nil)
	if err != nil {
		log.ErrorF("Open leveldb `%s` failed due to %s", leveldbFileName, err)
		os.Exit(data.STATUS_HALT)
	}
	defer ldb.Close()

	//open cache list file
	cacheResultFileHandle, err := os.Open(cacheResultName)
	if err != nil {
		log.ErrorF("Open list file `%s` failed due to %s", cacheResultName, err)
		os.Exit(data.STATUS_HALT)
	}
	defer cacheResultFileHandle.Close()
	bScanner := bufio.NewScanner(cacheResultFileHandle)

	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	//init wait group
	upWaitGroup := sync.WaitGroup{}

	initUpOnce.Do(func() {
		uploadTasks = make(chan func(), threadCount)
		for i := 0; i < threadCount; i++ {
			go doUpload(uploadTasks)
		}
	})

	bm, err := bucket.GetBucketManager()
	if err != nil {
		return
	}

	mac, err := workspace.GetMac()
	if err != nil {
		return
	}

	//scan lines and upload
	for bScanner.Scan() {
		line := bScanner.Text()
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			log.ErrorF("Invalid cache line `%s`", line)
			continue
		}

		localFileRelativePath := items[0]
		currentFileCount += 1

		//check skip local file or folder
		if skip, prefix := uploadConfig.HitByPathPrefixes(localFileRelativePath); skip {
			log.InfoF("Skip by path prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, prefix := uploadConfig.HitByFilePrefixes(localFileRelativePath); skip {
			log.InfoF("Skip by file prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, fixedStr := uploadConfig.HitByFixesString(localFileRelativePath); skip {
			log.InfoF("Skip by fixed string `%s` for local file path `%s`", fixedStr, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, suffix := uploadConfig.HitBySuffixes(localFileRelativePath); skip {
			log.InfoF("Skip by suffix `%s` for local file `%s`", suffix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		//pack the upload file key
		localFileLastModified, _ := strconv.ParseInt(items[2], 10, 64)
		uploadFileKey := localFileRelativePath

		//check ignore dir
		if uploadConfig.IgnoreDir {
			uploadFileKey = filepath.Base(uploadFileKey)
		}

		//check prefix
		if uploadConfig.KeyPrefix != "" {
			uploadFileKey = strings.Join([]string{uploadConfig.KeyPrefix, uploadFileKey}, "")
		}
		//convert \ to / under windows
		if runtime.GOOS == "windows" {
			uploadFileKey = strings.Replace(uploadFileKey, "\\", "/", -1)
		}

		localFilePath := filepath.Join(uploadConfig.SrcDir, localFileRelativePath)
		localFileStat, statErr := os.Stat(localFilePath)
		if statErr != nil {
			failureFileCount += 1
			log.ErrorF("Error stat local file `%s` due to `%s`", localFilePath, statErr)
			continue
		}

		localFileSize := localFileStat.Size()
		//check file encoding
		if strings.ToLower(uploadConfig.FileEncoding) == "gbk" {
			uploadFileKey, _ = utils.Gbk2Utf8(uploadFileKey)
		}

		ldbKey := fmt.Sprintf("%s => %s", localFilePath, uploadFileKey)

		if totalFileCount != 0 {
			fmt.Printf("Uploading %s [%d/%d, %.1f%%] ...\n", ldbKey, currentFileCount, totalFileCount,
				float32(currentFileCount)*100/float32(totalFileCount))
		} else {
			fmt.Printf("Uploading %s ...\n", ldbKey)
		}

		//check exists
		needToUpload := checkFileNeedToUpload(bm, uploadConfig, ldb, &ldbWOpt, ldbKey, localFilePath,
			uploadFileKey, localFileLastModified, localFileSize)
		if !needToUpload {
			//no need to upload
			continue
		}

		log.InfoF("Uploading file %s => %s : %s", localFilePath, uploadConfig.Bucket, uploadFileKey)

		//start to upload
		upWaitGroup.Add(1)
		uploadTasks <- func() {
			defer upWaitGroup.Done()

			upToken := uploadConfig.UploadToken(mac, uploadFileKey)
			if localFileSize > putThreshold {
				resumableUploadFile(uploadConfig, ldb, &ldbWOpt, ldbKey, upToken, storePath,
					localFilePath, uploadFileKey, localFileLastModified, exporter)
			} else {
				formUploadFile(uploadConfig, ldb, &ldbWOpt, ldbKey, upToken,
					localFilePath, uploadFileKey, localFileLastModified, exporter)
			}
		}
	}

	upWaitGroup.Wait()
	exporter.Close()

	log.InfoF("-------------Upload Result--------------")
	log.InfoF("%20s%10d", "Total:", totalFileCount)
	log.InfoF("%20s%10d", "Success:", successFileCount)
	log.InfoF("%20s%10d", "Failure:", failureFileCount)
	log.InfoF("%20s%10d", "NotOverwrite:", notOverwriteCount)
	log.InfoF("%20s%10d", "Skipped:", skippedFileCount)
	log.InfoF("%20s%15s", "Duration:", time.Since(timeStart))
	log.InfoF("----------------------------------------")
	log.Alert("See upload log at path:", uploadConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(data.STATUS_ERROR)
	} else {
		os.Exit(data.STATUS_OK)
	}
}

func formUploadFile(uploadConfig *config.UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string,
	upToken string, localFilePath, uploadFileKey string, localFileLastModified int64, exporter *export.FileExporter) {

	cfg := workspace.GetConfig()
	r := (&cfg).GetRegion()
	storage.UcHost = cfg.Hosts.UC[0]
	uploader := storage.NewFormUploader(&storage.Config{
		UseHTTPS:      cfg.IsUseHttps(),
		Zone:          r,
		Region:        r,
		CentralRsHost: cfg.Hosts.GetOneRs(),
	})
	putRet := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := uploader.PutFile(workspace.GetContext(), &putRet, upToken, uploadFileKey, localFilePath, &putExtra)
	if err != nil {
		atomic.AddInt64(&failureFileCount, 1)
		log.ErrorF("Form upload file `%s` => `%s` failed due to nerror `%v`", localFilePath, uploadFileKey, err)
		exporter.Fail().ExportF("%s\t%s\t%v\n", localFilePath, uploadFileKey, err)
	} else {
		atomic.AddInt64(&successFileCount, 1)
		log.InfoF("Upload file `%s` => `%s : %s` success", localFilePath, uploadConfig.Bucket, uploadFileKey)
		putErr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFileLastModified)), ldbWOpt)
		if putErr != nil {
			log.ErrorF("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
		}
		//delete on success
		if uploadConfig.DeleteOnSuccess {
			deleteErr := os.Remove(localFilePath)
			if deleteErr != nil {
				log.ErrorF("Delete `%s` on upload success error due to `%s`", localFilePath, deleteErr)
			} else {
				log.InfoF("Delete `%s` on upload success done", localFilePath)
			}
		}

		exporter.Success().ExportF("%s\t%s\n", localFilePath, uploadFileKey)
		if uploadConfig.Overwrite {
			exporter.Override().ExportF("%s\t%s\n", localFilePath, uploadFileKey)
		}
	}
}

var progressRecorder = NewProgressRecorder("")

func resumableUploadFile(uploadConfig *config.UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string, upToken string, storePath, localFilePath, uploadFileKey string, localFileLastModified int64, exporter *export.FileExporter) {
	uploadConfig.Check()

	var progressFilePath string
	if !uploadConfig.DisableResume {
		//progress file
		progressFileKey := utils.Md5Hex(fmt.Sprintf("%s:%s|%s:%s", uploadConfig.SrcDir,
			uploadConfig.Bucket, localFilePath, uploadFileKey))
		progressFilePath = filepath.Join(storePath, fmt.Sprintf("%s.progress", progressFileKey))
		progressRecorder.FilePath = progressFilePath
	}

	cfg := workspace.GetConfig()
	r := (&cfg).GetRegion()

	var err error
	if uploadConfig.ResumableAPIV2 {
		partSize := uploadConfig.ResumableAPIV2PartSize
		log.DebugF("uploadFileKey: %s, partSize: %d", uploadFileKey, partSize)
		var notifyFunc = func(partNumber int64, ret *storage.UploadPartsRet) {
			log.DebugF("uploadFileKey: %s, partIdx: %d, partSize: %d, %v", uploadFileKey, partNumber, partSize, *ret)
			progressRecorder.Parts = append(progressRecorder.Parts, storage.UploadPartInfo{
				Etag:       ret.Etag,
				PartNumber: partNumber,
			})
			progressRecorder.Offset += partSize
		}

		//params
		putRet := storage.UploadPartsRet{}
		putExtra := storage.RputV2Extra{
			PartSize: partSize,
		}
		putExtra.Notify = notifyFunc

		//resumable upload
		uploader := storage.NewResumeUploaderV2(&storage.Config{
			UseHTTPS:      cfg.IsUseHttps(),
			Zone:          r,
			Region:        r,
			CentralRsHost: cfg.Hosts.GetOneRs(),
		})
		err = uploader.PutFile(workspace.GetContext(), &putRet, upToken, uploadFileKey, localFilePath, &putExtra)
	} else {

		var notifyFunc = func(blkIdx, blkSize int, ret *storage.BlkputRet) {
			log.DebugF("uploadFileKey: %s, blkIdx: %d, blkSize: %d, %v", uploadFileKey, blkIdx, blkSize, *ret)
			progressRecorder.BlkCtxs = append(progressRecorder.BlkCtxs, *ret)
			progressRecorder.Offset += int64(blkSize)
		}

		//params
		putRet := storage.PutRet{}
		putExtra := storage.RputExtra{}
		putExtra.Notify = notifyFunc

		//resumable upload
		uploader := storage.NewResumeUploader(&storage.Config{
			UseHTTPS:      cfg.IsUseHttps(),
			Zone:          r,
			Region:        r,
			CentralRsHost: cfg.Hosts.GetOneRs(),
		})
		err = uploader.PutFile(workspace.GetContext(), &putRet, upToken, uploadFileKey, localFilePath, &putExtra)
	}

	if err != nil {
		atomic.AddInt64(&failureFileCount, 1)
		log.ErrorF("Resumable upload file `%s` => `%s` failed due to nerror `%v`", localFilePath, uploadFileKey, err)
		exporter.Fail().Export("%s\t%s\t%v\n", localFilePath, uploadFileKey, err)
	} else {
		if progressFilePath != "" {
			os.Remove(progressFilePath)
		}

		atomic.AddInt64(&successFileCount, 1)
		log.InfoF("Upload file `%s` => `%s` success", localFilePath, uploadFileKey)
		putErr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFileLastModified)), ldbWOpt)
		if putErr != nil {
			log.ErrorF("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
		}
		//delete on success

		if uploadConfig.DeleteOnSuccess {
			deleteErr := os.Remove(localFilePath)
			if deleteErr != nil {
				log.ErrorF("Delete `%s` on upload success error due to `%s`", localFilePath, deleteErr)
			} else {
				log.InfoF("Delete `%s` on upload success done", localFilePath)
			}
		}
		exporter.Success().ExportF("%s\t%s\n", localFilePath, uploadFileKey)
		if uploadConfig.Overwrite {
			exporter.Override().ExportF("%s\t%s\n", localFilePath, uploadFileKey)
		}
	}
}

func checkFileNeedToUpload(bm *storage.BucketManager, uploadConfig *config.UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions,
	ldbKey, localFilePath, uploadFileKey string, localFileLastModified, localFileSize int64) (needToUpload bool) {
	//default to upload
	needToUpload = true

	//check before to upload
	if uploadConfig.CheckExists {
		rsEntry, sErr := bm.Stat(uploadConfig.Bucket, uploadFileKey)
		if sErr != nil {
			if _, ok := sErr.(*storage.ErrorInfo); !ok {
				//not logic error, should be network error
				log.ErrorF("Get file `%s` stat error, %s", uploadFileKey, sErr)
				atomic.AddInt64(&failureFileCount, 1)
				needToUpload = false
			}
			return
		}
		ldbValue := fmt.Sprintf("%d", localFileLastModified)
		if uploadConfig.CheckHash {
			//compare hash
			localEtag, cErr := utils.GetEtag(localFilePath)
			if cErr != nil {
				log.ErrorF("File `%s` calc local hash failed, %s", uploadFileKey, cErr)
				atomic.AddInt64(&failureFileCount, 1)
				needToUpload = false
				return
			}
			if rsEntry.Hash == localEtag {
				log.InfoF("File `%s` exists in bucket, hash match, ignore this upload", uploadFileKey)
				atomic.AddInt64(&skippedFileCount, 1)
				putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
				if putErr != nil {
					log.ErrorF("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
				}
				needToUpload = false
				return
			}
			if !uploadConfig.Overwrite {
				log.WarningF("Skip upload of unmatch hash file `%s` because `overwrite` is false", localFilePath)
				atomic.AddInt64(&notOverwriteCount, 1)
				needToUpload = false
				return
			}
			log.InfoF("File `%s` exists in bucket, but hash not match, go to upload", uploadFileKey)
		} else {
			if uploadConfig.CheckSize {
				if rsEntry.Fsize == localFileSize {
					log.InfoF("File `%s` exists in bucket, size match, ignore this upload", uploadFileKey)
					atomic.AddInt64(&skippedFileCount, 1)
					putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
					if putErr != nil {
						log.ErrorF("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
					}
					needToUpload = false
					return
				}
				if !uploadConfig.Overwrite {
					log.WarningF("Skip upload of unmatch size file `%s` because `overwrite` is false", localFilePath)
					atomic.AddInt64(&notOverwriteCount, 1)
					needToUpload = false
					return
				}
				log.InfoF("File `%s` exists in bucket, but size not match, go to upload", uploadFileKey)
			} else {
				log.InfoF("File `%s` exists in bucket, no hash or size check, ignore this upload", uploadFileKey)
				atomic.AddInt64(&skippedFileCount, 1)
				putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
				if putErr != nil {
					log.ErrorF("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
				}
				needToUpload = false
				return
			}
		}

	} else {
		//check leveldb
		ldbFlmd, err := ldb.Get([]byte(ldbKey), nil)
		flmd, _ := strconv.ParseInt(string(ldbFlmd), 10, 64)
		//not exist, return ErrNotFound
		//check last modified

		if err != nil {
			return
		}
		if localFileLastModified == flmd {
			log.InfoF("Skip by local leveldb log for file `%s`", localFilePath)
			atomic.AddInt64(&skippedFileCount, 1)
			needToUpload = false
		} else {
			if !uploadConfig.Overwrite {
				//no overwrite set
				log.WarningF("Skip upload of changed file `%s` because `overwrite` is false", localFilePath)
				atomic.AddInt64(&notOverwriteCount, 1)
				needToUpload = false
			}
		}

		return
	}
	return
}
