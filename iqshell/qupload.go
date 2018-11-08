package qshell

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/storage"
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

/*
Config file like:

{
	"src_dir"		:	"/Users/jemy/Photos",
	"file_list"             :       "",
	"bucket"		:	"test-bucket",
	"put_threshold"		:	10000000,
	"key_prefix"		:	"2014/12/01/",
	"ignore_dir"		:	false,
	"overwrite"		:	false,
	"check_exists"		:	true,
	"skip_file_prefixes"	:	"IMG_",
	"skip_path_prefixes"	:	"tmp/,bin/,obj/",
	"skip_suffixes"		:	".exe,.obj,.class",
	"skip_fixed_strings"    :       ".svn,.git",
	"up_host"		:	"http://upload.qiniu.com",
	"bind_up_ip"		:	"",
	"bind_rs_ip"		:	"",
	"bind_nic_ip"		:	"",
	"rescan_local"		:	false,
	"delete_on_success" :   false
 }

or the simplest one

{
	"src_dir" 	:	"/Users/jemy/Photos",
	"bucket"	:	"test-bucket",
}

*/

const (
	DEFAULT_PUT_THRESHOLD   int64 = 10 * 1024 * 1024 //10MB
	MIN_UPLOAD_THREAD_COUNT int64 = 1
	MAX_UPLOAD_THREAD_COUNT int64 = 2000
)

type UploadInfo struct {
	TotalFileCount int64 `json:"total_file_count"`
}

type UploadConfig struct {
	//basic config
	SrcDir string `json:"src_dir"`
	Bucket string `json:"bucket"`

	//optional config
	FileList         string `json:"file_list,omitempty"`
	PutThreshold     int64  `json:"put_threshold,omitempty"`
	KeyPrefix        string `json:"key_prefix,omitempty"`
	IgnoreDir        bool   `json:"ignore_dir,omitempty"`
	Overwrite        bool   `json:"overwrite,omitempty"`
	CheckExists      bool   `json:"check_exists,omitempty"`
	CheckHash        bool   `json:"check_hash,omitempty"`
	CheckSize        bool   `json:"check_size,omitempty"`
	SkipFilePrefixes string `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes string `json:"skip_path_prefixes,omitempty"`
	SkipFixedStrings string `json:"skip_fixed_strings,omitempty"`
	SkipSuffixes     string `json:"skip_suffixes,omitempty"`
	RescanLocal      bool   `json:"rescan_local,omitempty"`
	FileType         int    `json:"file_type,omitempty"`

	//advanced config
	UpHost string `json:"up_host,omitempty"`

	BindUpIp string `json:"bind_up_ip,omitempty"`
	BindRsIp string `json:"bind_rs_ip,omitempty"`
	//local network interface card config
	BindNicIp string `json:"bind_nic_ip,omitempty"`

	//log settings
	LogLevel  string `json:"log_level,omitempty"`
	LogFile   string `json:"log_file,omitempty"`
	LogRotate int    `json:"log_rotate,omitempty"`
	LogStdout bool   `json:"log_stdout,omitempty"`

	//more settings
	DeleteOnSuccess bool `json:"delete_on_success,omitempty"`
	DisableResume   bool `json:"disable_resume,omitempty"`
}

func (cfg *UploadConfig) JobId() string {

	return Md5Hex(fmt.Sprintf("%s:%s", cfg.SrcDir, cfg.Bucket))
}

func (cfg *UploadConfig) GetLogLevel() int {

	//init log level
	logLevel := logs.LevelInformational
	switch cfg.LogLevel {
	case "debug":
		logLevel = logs.LevelDebug
	case "info":
		logLevel = logs.LevelInformational
	case "warn":
		logLevel = logs.LevelWarning
	case "error":
		logLevel = logs.LevelError
	default:
		logLevel = logs.LevelInformational
	}
	return logLevel
}

func (cfg *UploadConfig) GetLogRotate() int {
	logRotate := 1
	if cfg.LogRotate > 0 {
		logRotate = cfg.LogRotate
	}
	return logRotate
}

func (cfg *UploadConfig) HitByPathPrefixes(localFileRelativePath string) (hit bool, pathPrefix string) {

	if cfg.SkipPathPrefixes != "" {
		//unpack skip prefix
		pathPrefixes := strings.Split(cfg.SkipPathPrefixes, ",")
		for _, prefix := range pathPrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			if strings.HasPrefix(localFileRelativePath, prefix) {
				pathPrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) HitByFilePrefixes(localFileRelativePath string) (hit bool, filePrefix string) {
	if cfg.SkipFilePrefixes != "" {
		//unpack skip prefix
		filePrefixes := strings.Split(cfg.SkipFilePrefixes, ",")
		for _, prefix := range filePrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			localFileName := filepath.Base(localFileRelativePath)
			if strings.HasPrefix(localFileName, prefix) {
				filePrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) HitByFixesString(localFileRelativePath string) (hit bool, hitFixedStr string) {
	if cfg.SkipFixedStrings != "" {
		//unpack fixed strings
		fixedStrings := strings.Split(cfg.SkipFixedStrings, ",")
		for _, fixedStr := range fixedStrings {
			if strings.TrimSpace(fixedStr) == "" {
				continue
			}

			if strings.Contains(localFileRelativePath, fixedStr) {
				hitFixedStr = fixedStr
				hit = true
				break
			}
		}
	}
	return

}

func (cfg *UploadConfig) HitBySuffixes(localFileRelativePath string) (hit bool, hitSuffix string) {
	if cfg.SkipSuffixes != "" {
		suffixes := strings.Split(cfg.SkipSuffixes, ",")
		for _, suffix := range suffixes {
			if strings.TrimSpace(suffix) == "" {
				continue
			}

			if strings.HasSuffix(localFileRelativePath, suffix) {
				hitSuffix = suffix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) CacheFileNameAndCount(storePath, jobId string) (cacheResultName string, totalFileCount int64, cacheErr error) {
	//find the local file list, by specified or by config

	_, localFileStatErr := os.Stat(cfg.FileList)
	if cfg.FileList != "" && localFileStatErr == nil {
		//use specified file list
		cacheResultName = cfg.FileList
		totalFileCount = GetFileLineCount(cacheResultName)
	} else {
		cacheResultName = filepath.Join(storePath, fmt.Sprintf("%s.cache", jobId))
		cacheCountName := filepath.Join(storePath, fmt.Sprintf("%s.count", jobId))
		totalFileCount, cacheErr = prepareCacheFileList(cacheResultName, cacheCountName, cfg.SrcDir, cfg.RescanLocal)
		if cacheErr != nil {
			return
		}
	}
	return

}

func (cfg *UploadConfig) DefaultLogFile(storePath, jobId string) (defaultLogFile string, err error) {
	//local storage path
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		err = mkdirErr
		return
	}
	defaultLogFile = filepath.Join(storePath, fmt.Sprintf("%s.log", jobId))
	return
}

func (cfg *UploadConfig) PrepareLogger(storePath, jobId string) {

	defaultLogFile, err := cfg.DefaultLogFile(storePath, jobId)
	if err != nil {
		os.Exit(STATUS_HALT)
	}
	logLevel := cfg.GetLogLevel()
	logRotate := cfg.GetLogRotate()

	//init log writer
	if cfg.LogFile == "" {
		//set default log file
		cfg.LogFile = defaultLogFile
	}

	if !cfg.LogStdout {
		logs.GetBeeLogger().DelLogger(logs.AdapterConsole)
	}
	//open log file
	fmt.Println("Writing upload log to file", cfg.LogFile)

	//daily rotate
	logCfg := BeeLogConfig{
		Filename: cfg.LogFile,
		Level:    logLevel,
		Daily:    true,
		MaxDays:  logRotate,
	}
	logs.SetLogger(logs.AdapterFile, logCfg.ToJson())
	fmt.Println()
}

var uploadTasks chan func()
var initUpOnce sync.Once

func doUpload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

// FileExporter
type FileExporter struct {
	SuccessFhandle   *os.File
	SuccessLock      sync.RWMutex
	SuccessWriter    *bufio.Writer
	FailureFhandle   *os.File
	FailureLock      sync.RWMutex
	FailureWriter    *bufio.Writer
	OverwriteFhandle *os.File
	OverwriteLock    sync.RWMutex
	OverwriteWriter  *bufio.Writer
}

func NewFileExporter(successFname, failureFname, overwriteFname string) (ex *FileExporter, err error) {
	ex = new(FileExporter)
	//init file list writer
	var (
		successListFp   *os.File
		failureListFp   *os.File
		overwriteListFp *os.File
		openErr         error
	)

	if successFname != "" {
		successListFp, openErr = os.Create(successFname)
		if openErr != nil {
			err = fmt.Errorf("open file: %s: %v\n", successFname, openErr)
			return
		}
		ex.SuccessFhandle = successListFp
		ex.SuccessWriter = bufio.NewWriter(successListFp)
	}

	if failureFname != "" {
		failureListFp, openErr = os.Create(failureFname)
		if openErr != nil {
			err = fmt.Errorf("open file: %s: %v\n", failureFname, openErr)
			return
		}
		ex.FailureFhandle = failureListFp
		ex.FailureWriter = bufio.NewWriter(failureListFp)
	}

	if overwriteFname != "" {
		overwriteListFp, openErr = os.Create(overwriteFname)
		if openErr != nil {
			err = fmt.Errorf("open file: %s: %v\n", overwriteFname, openErr)
			return
		}
		ex.OverwriteFhandle = overwriteListFp
		ex.OverwriteWriter = bufio.NewWriter(overwriteListFp)
	}
	return
}

func (ex *FileExporter) WriteToFailedWriter(localFilePath, uploadFileKey string, err error) {
	if ex.FailureWriter != nil {
		ex.FailureLock.Lock()
		ex.FailureWriter.WriteString(fmt.Sprintf("%s\t%s\t%v\n", localFilePath, uploadFileKey, err))
		ex.FailureWriter.Flush()
		ex.FailureLock.Unlock()
	}
}

func (ex *FileExporter) WriteToSuccessWriter(localFilePath, uploadFileKey string) {
	if ex.SuccessWriter != nil {
		ex.SuccessLock.Lock()
		ex.SuccessWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
		ex.SuccessWriter.Flush()
		ex.SuccessLock.Unlock()
	}
}

func (ex *FileExporter) WriteToOverwriter(localFilePath, uploadFileKey string) {
	if ex.OverwriteWriter != nil {
		ex.OverwriteLock.Lock()
		ex.OverwriteWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
		ex.OverwriteWriter.Flush()
		ex.OverwriteLock.Unlock()
	}
}

func (ex *FileExporter) Close() {
	ex.SuccessFhandle.Close()
	ex.FailureFhandle.Close()
	ex.OverwriteFhandle.Close()
}

func (ex *FileExporter) FlushWriter() {

	//flush
	exportWriters := []*bufio.Writer{
		ex.SuccessWriter,
		ex.FailureWriter,
		ex.OverwriteWriter,
	}

	for _, writer := range exportWriters {
		if writer != nil {
			writer.Flush()
		}
	}
}

var (
	currentFileCount  int64
	successFileCount  int64
	notOverwriteCount int64
	failureFileCount  int64
	skippedFileCount  int64
)

// QiniuUpload
func QiniuUpload(threadCount int, uploadConfig *UploadConfig, exporter *FileExporter) {
	var upSettings = storage.Settings{
		Workers:   16,
		ChunkSize: 4 * 1024 * 1024,
		TryTimes:  3,
	}
	timeStart := time.Now()
	//create job id
	jobId := uploadConfig.JobId()
	QShellRootPath := RootPath()
	if QShellRootPath == "" {
		logs.Error("Empty root path")
		os.Exit(STATUS_HALT)
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
		os.Exit(STATUS_HALT)
	}

	//leveldb folder
	leveldbFileName := filepath.Join(storePath, jobId+".ldb")
	ldb, err := leveldb.OpenFile(leveldbFileName, nil)
	if err != nil {
		logs.Error("Open leveldb `%s` failed due to %s", leveldbFileName, err)
		os.Exit(STATUS_HALT)
	}
	defer ldb.Close()

	//open cache list file
	cacheResultFileHandle, err := os.Open(cacheResultName)
	if err != nil {
		logs.Error("Open list file `%s` failed due to %s", cacheResultName, err)
		os.Exit(STATUS_HALT)
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

	bm := GetBucketManager()
	//scan lines and upload
	for bScanner.Scan() {
		line := bScanner.Text()
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			logs.Error("Invalid cache line `%s`", line)
			continue
		}

		localFileRelativePath := items[0]
		currentFileCount += 1

		//check skip local file or folder
		if skip, prefix := uploadConfig.HitByPathPrefixes(localFileRelativePath); skip {
			logs.Informational("Skip by path prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, prefix := uploadConfig.HitByFilePrefixes(localFileRelativePath); skip {
			logs.Informational("Skip by file prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, fixedStr := uploadConfig.HitByFixesString(localFileRelativePath); skip {
			logs.Informational("Skip by fixed string `%s` for local file path `%s`", fixedStr, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, suffix := uploadConfig.HitBySuffixes(localFileRelativePath); skip {
			logs.Informational("Skip by suffix `%s` for local file `%s`", suffix, localFileRelativePath)
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
			logs.Error("Error stat local file `%s` due to `%s`", localFilePath, statErr)
			continue
		}

		localFileSize := localFileStat.Size()
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

		logs.Informational("Uploading file %s => %s : %s", localFilePath, uploadConfig.Bucket, uploadFileKey)

		//start to upload
		upWaitGroup.Add(1)
		uploadTasks <- func() {
			defer upWaitGroup.Done()

			policy := storage.PutPolicy{}
			policy.Scope = uploadConfig.Bucket
			if uploadConfig.Overwrite {
				policy.Scope = fmt.Sprintf("%s:%s", uploadConfig.Bucket, uploadFileKey)
				policy.InsertOnly = 0
			}

			policy.FileType = uploadConfig.FileType

			policy.Expires = 7 * 24 * 3600
			upToken := policy.UploadToken(bm.GetMac())

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
	exporter.FlushWriter()
	exporter.Close()

	logs.Informational("-------------Upload Result--------------")
	logs.Informational("%20s%10d", "Total:", totalFileCount)
	logs.Informational("%20s%10d", "Success:", successFileCount)
	logs.Informational("%20s%10d", "Failure:", failureFileCount)
	logs.Informational("%20s%10d", "NotOverwrite:", notOverwriteCount)
	logs.Informational("%20s%10d", "Skipped:", skippedFileCount)
	logs.Informational("%20s%15s", "Duration:", time.Since(timeStart))
	logs.Info("----------------------------------------")
	fmt.Println("\nSee upload log at path", uploadConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(STATUS_ERROR)
	} else {
		os.Exit(STATUS_OK)
	}
}

func prepareCacheFileList(cacheResultName, cacheCountName, srcDir string, rescanLocal bool) (totalFileCount int64, cacheErr error) {
	//cache file
	cacheTempName := fmt.Sprintf("%s.temp", cacheResultName)

	rescanLocalDir := false
	if _, statErr := os.Stat(cacheResultName); statErr == nil {
		//file exists
		rescanLocalDir = rescanLocal
	} else {
		rescanLocalDir = true
	}

	if rescanLocalDir {
		logs.Informational("Listing local sync dir, this can take a long time for big directory, please wait patiently")
		totalFileCount, cacheErr = DirCache(srcDir, cacheTempName)
		if cacheErr != nil {
			return
		}

		if rErr := os.Rename(cacheTempName, cacheResultName); rErr != nil {
			logs.Error("Rename the temp cached file error", rErr)
			cacheErr = rErr
			return
		}
		//write the total count to local file
		if cFp, cErr := os.Create(cacheCountName); cErr == nil {
			func() {
				defer cFp.Close()
				uploadInfo := UploadInfo{
					TotalFileCount: totalFileCount,
				}
				uploadInfoBytes, mErr := json.Marshal(&uploadInfo)
				if mErr == nil {
					if _, wErr := cFp.Write(uploadInfoBytes); wErr != nil {
						logs.Warning("Write local cached count file error %s", cErr)
					} else {
						cFp.Close()
					}
				}
			}()
		} else {
			logs.Error("Open local cached count file error %s,", cErr)
		}
	} else {
		logs.Informational("Use the last cached local sync dir file list")
		//read from local cache
		if rFp, rErr := os.Open(cacheCountName); rErr == nil {
			func() {
				defer rFp.Close()
				uploadInfo := UploadInfo{}
				decoder := json.NewDecoder(rFp)
				if dErr := decoder.Decode(&uploadInfo); dErr == nil {
					totalFileCount = uploadInfo.TotalFileCount
				}
			}()
		} else {
			logs.Warning("Open local cached count file error %s,", rErr)
			totalFileCount = GetFileLineCount(cacheResultName)
		}
	}

	return
}

func checkFileNeedToUpload(bm *BucketManager, uploadConfig *UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions,
	ldbKey, localFilePath, uploadFileKey string, localFileLastModified, localFileSize int64) (needToUpload bool) {
	//default to upload
	needToUpload = true

	//check before to upload
	if uploadConfig.CheckExists {
		rsEntry, sErr := bm.Stat(uploadConfig.Bucket, uploadFileKey)
		if sErr != nil {
			if _, ok := sErr.(*storage.ErrorInfo); !ok {
				//not logic error, should be network error
				logs.Error("Get file `%s` stat error, %s", uploadFileKey, sErr)
				atomic.AddInt64(&failureFileCount, 1)
				needToUpload = false
			}
			return
		}
		ldbValue := fmt.Sprintf("%d", localFileLastModified)
		if uploadConfig.CheckHash {
			//compare hash
			localEtag, cErr := GetEtag(localFilePath)
			if cErr != nil {
				logs.Error("File `%s` calc local hash failed, %s", uploadFileKey, cErr)
				atomic.AddInt64(&failureFileCount, 1)
				needToUpload = false
				return
			}
			if rsEntry.Hash == localEtag {
				logs.Informational("File `%s` exists in bucket, hash match, ignore this upload", uploadFileKey)
				atomic.AddInt64(&skippedFileCount, 1)
				putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
				if putErr != nil {
					logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
				}
				needToUpload = false
				return
			}
			if !uploadConfig.Overwrite {
				logs.Warning("Skip upload of unmatch hash file `%s` because `overwrite` is false", localFilePath)
				atomic.AddInt64(&notOverwriteCount, 1)
				needToUpload = false
				return
			}
			logs.Informational("File `%s` exists in bucket, but hash not match, go to upload", uploadFileKey)
		} else {
			if uploadConfig.CheckSize {
				if rsEntry.Fsize == localFileSize {
					logs.Info("File `%s` exists in bucket, size match, ignore this upload", uploadFileKey)
					atomic.AddInt64(&skippedFileCount, 1)
					putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
					if putErr != nil {
						logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
					}
					needToUpload = false
					return
				}
				if !uploadConfig.Overwrite {
					logs.Warning("Skip upload of unmatch size file `%s` because `overwrite` is false", localFilePath)
					atomic.AddInt64(&notOverwriteCount, 1)
					needToUpload = false
					return
				}
				logs.Info("File `%s` exists in bucket, but size not match, go to upload", uploadFileKey)
			} else {
				logs.Info("File `%s` exists in bucket, no hash or size check, ignore this upload", uploadFileKey)
				atomic.AddInt64(&skippedFileCount, 1)
				putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
				if putErr != nil {
					logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
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
			logs.Informational("Skip by local leveldb log for file `%s`", localFilePath)
			atomic.AddInt64(&skippedFileCount, 1)
			needToUpload = false
		} else {
			if !uploadConfig.Overwrite {
				//no overwrite set
				logs.Warning("Skip upload of changed file `%s` because `overwrite` is false", localFilePath)
				atomic.AddInt64(&notOverwriteCount, 1)
				needToUpload = false
			}
		}

		return
	}
	return
}

func formUploadFile(uploadConfig *UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string,
	upToken string, localFilePath, uploadFileKey string, localFileLastModified int64, exporter *FileExporter) {

	uploader := storage.NewFormUploader(nil)
	putRet := storage.PutRet{}

	err := uploader.PutFile(context.Background(), &putRet, upToken, uploadFileKey, localFilePath, nil)
	if err != nil {
		atomic.AddInt64(&failureFileCount, 1)
		logs.Error("Form upload file `%s` => `%s` failed due to nerror `%v`", localFilePath, uploadFileKey, err)
		exporter.WriteToFailedWriter(localFilePath, uploadFileKey, err)
	} else {
		atomic.AddInt64(&successFileCount, 1)
		logs.Informational("Upload file `%s` => `%s : %s` success", localFilePath, uploadConfig.Bucket, uploadFileKey)
		putErr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFileLastModified)), ldbWOpt)
		if putErr != nil {
			logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
		}
		//delete on success
		if uploadConfig.DeleteOnSuccess {
			deleteErr := os.Remove(localFilePath)
			if deleteErr != nil {
				logs.Error("Delete `%s` on upload success error due to `%s`", localFilePath, deleteErr)
			} else {
				logs.Info("Delete `%s` on upload success done", localFilePath)
			}
		}

		exporter.WriteToSuccessWriter(localFilePath, uploadFileKey)
		if uploadConfig.Overwrite {
			exporter.WriteToOverwriter(localFilePath, uploadFileKey)
		}
	}
}

var progressRecorder = NewProgressRecorder("")

func resumableUploadFile(uploadConfig *UploadConfig, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string, upToken string, storePath, localFilePath, uploadFileKey string, localFileLastModified int64, exporter *FileExporter) {

	uploader := storage.NewResumeUploader(nil)
	//params
	putRet := storage.PutRet{}
	putExtra := storage.RputExtra{}

	var progressFilePath string
	if !uploadConfig.DisableResume {
		//progress file
		progressFileKey := Md5Hex(fmt.Sprintf("%s:%s|%s:%s", uploadConfig.SrcDir,
			uploadConfig.Bucket, localFilePath, uploadFileKey))
		progressFilePath = filepath.Join(storePath, fmt.Sprintf("%s.progress", progressFileKey))
		progressRecorder.FilePath = progressFilePath
	}

	var notifyFunc = func(blkIdx, blkSize int, ret *storage.BlkputRet) {
		progressRecorder.BlkCtxs = append(progressRecorder.BlkCtxs, *ret)
		progressRecorder.Offset += int64(blkSize)
	}
	putExtra.Notify = notifyFunc

	//resumable upload
	err := uploader.PutFile(context.Background(), &putRet, upToken, uploadFileKey, localFilePath, &putExtra)
	if err != nil {
		atomic.AddInt64(&failureFileCount, 1)
		logs.Error("Resumable upload file `%s` => `%s` failed due to nerror `%v`", localFilePath, uploadFileKey, err)
		exporter.WriteToFailedWriter(localFilePath, uploadFileKey, err)
	} else {
		if progressFilePath != "" {
			os.Remove(progressFilePath)
		}

		atomic.AddInt64(&successFileCount, 1)
		logs.Informational("Upload file `%s` => `%s` success", localFilePath, uploadFileKey)
		putErr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFileLastModified)), ldbWOpt)
		if putErr != nil {
			logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
		}
		//delete on success
		if uploadConfig.DeleteOnSuccess {
			deleteErr := os.Remove(localFilePath)
			if deleteErr != nil {
				logs.Error("Delete `%s` on upload success error due to `%s`", localFilePath, deleteErr)
			} else {
				logs.Info("Delete `%s` on upload success done", localFilePath)
			}
		}
		exporter.WriteToSuccessWriter(localFilePath, uploadFileKey)
		if uploadConfig.Overwrite {
			exporter.WriteToOverwriter(localFilePath, uploadFileKey)
		}
	}
}
