package qshell

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	fio "qiniu/api.v6/io"
	rio "qiniu/api.v6/resumable/io"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
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
}

var upSettings = rio.Settings{
	Workers:   16,
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  3,
}

var uploadTasks chan func()
var initUpOnce sync.Once

func doUpload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

var currentFileCount int64
var successFileCount int64
var notOverwriteCount int64
var failureFileCount int64
var skippedFileCount int64

// FileExporter
type FileExporter struct {
	SuccessFname    string
	SuccessLock     sync.RWMutex
	SuccessWriter   *bufio.Writer
	FailureFname    string
	FailureLock     sync.RWMutex
	FailureWriter   *bufio.Writer
	OverwriteFname  string
	OverwriteLock   sync.RWMutex
	OverwriteWriter *bufio.Writer
}

// QiniuUpload
func QiniuUpload(threadCount int, uploadConfig *UploadConfig, exporter *FileExporter) {
	timeStart := time.Now()
	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "qupload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		os.Exit(STATUS_HALT)
	}

	defaultLogFile := filepath.Join(storePath, fmt.Sprintf("%s.log", jobId))
	//init log level
	logLevel := logs.LevelInformational
	logRotate := 1
	if uploadConfig.LogRotate > 0 {
		logRotate = uploadConfig.LogRotate
	}
	switch uploadConfig.LogLevel {
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

	//init log writer
	if uploadConfig.LogFile == "" {
		//set default log file
		uploadConfig.LogFile = defaultLogFile
	}

	if !uploadConfig.LogStdout {
		logs.GetBeeLogger().DelLogger(logs.AdapterConsole)
	}
	//open log file
	fmt.Println("Writing upload log to file", uploadConfig.LogFile)

	//daily rotate
	logCfg := BeeLogConfig{
		Filename: uploadConfig.LogFile,
		Level:    logLevel,
		Daily:    true,
		MaxDays:  logRotate,
	}
	logs.SetLogger(logs.AdapterFile, logCfg.ToJson())
	fmt.Println()

	//init file list writer
	var successListFp *os.File
	var failureListFp *os.File
	var overwriteListFp *os.File
	var openErr error

	var successListWriter *bufio.Writer
	var failureListWriter *bufio.Writer
	var overwriteListWriter *bufio.Writer
	if exporter.SuccessFname != "" {
		successListFp, openErr = os.Create(exporter.SuccessFname)
		if openErr != nil {
			logs.Error("Open success list file error, %s", openErr)
		} else {
			defer successListFp.Close()
			successListWriter = bufio.NewWriter(successListFp)
			exporter.SuccessWriter = successListWriter
		}
	}

	if exporter.FailureFname != "" {
		failureListFp, openErr = os.Create(exporter.FailureFname)
		if openErr != nil {
			logs.Error("Open fail list file error, %s", openErr)
		} else {
			defer failureListFp.Close()
			failureListWriter = bufio.NewWriter(failureListFp)
			exporter.FailureWriter = failureListWriter
		}
	}

	if exporter.OverwriteFname != "" {
		overwriteListFp, openErr = os.Create(exporter.OverwriteFname)
		if openErr != nil {
			logs.Error("Open overwrite list file error, %s", openErr)
		} else {
			defer overwriteListFp.Close()
			overwriteListWriter = bufio.NewWriter(overwriteListFp)
			exporter.OverwriteWriter = overwriteListWriter
		}
	}

	//global up settings
	logs.Info("Load account from %s", filepath.Join(QShellRootPath, ".qshell/account.json"))
	account, gErr := GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(STATUS_HALT)
	}
	mac := digest.Mac{AccessKey: account.AccessKey, SecretKey: []byte(account.SecretKey)}
	//get bucket zone info
	bucketInfo, gErr := GetBucketInfo(&mac, uploadConfig.Bucket)
	if gErr != nil {
		logs.Error("Get bucket region info error,", gErr)
		os.Exit(STATUS_HALT)
	}

	//set up host
	SetZone(bucketInfo.Region)

	//chunk upload threshold
	putThreshold := DEFAULT_PUT_THRESHOLD
	if uploadConfig.PutThreshold > 0 {
		putThreshold = uploadConfig.PutThreshold
	}

	//use host if not empty, overwrite the default config
	if uploadConfig.UpHost != "" {
		conf.UP_HOST = strings.TrimSuffix(uploadConfig.UpHost, "/")
	}
	//set resume upload settings
	rio.SetSettings(&upSettings)

	//make SrcDir the full path
	uploadConfig.SrcDir, _ = filepath.Abs(uploadConfig.SrcDir)

	//find the local file list, by specified or by config
	var cacheResultName string
	var cacheCountName string
	var totalFileCount int64
	var cacheErr error
	_, localFileStatErr := os.Stat(uploadConfig.FileList)
	if uploadConfig.FileList != "" && localFileStatErr == nil {
		//use specified file list
		cacheResultName = uploadConfig.FileList
		totalFileCount = GetFileLineCount(cacheResultName)
	} else {
		cacheResultName = filepath.Join(storePath, fmt.Sprintf("%s.cache", jobId))
		cacheCountName = filepath.Join(storePath, fmt.Sprintf("%s.count", jobId))
		totalFileCount, cacheErr = prepareCacheFileList(cacheResultName, cacheCountName,
			uploadConfig.SrcDir, uploadConfig.RescanLocal)
		if cacheErr != nil {
			os.Exit(STATUS_HALT)
		}
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
	bScanner.Split(bufio.ScanLines)

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

	//check bind net interface card
	var transport *http.Transport
	var rsClient rs.Client
	if uploadConfig.BindNicIp != "" {
		transport = &http.Transport{
			Dial: (&net.Dialer{
				LocalAddr: &net.TCPAddr{
					IP: net.ParseIP(uploadConfig.BindNicIp),
				},
			}).Dial,
		}
	}

	if transport != nil {
		rsClient = rs.NewMacEx(&mac, transport, "")
	} else {
		rsClient = rs.NewMac(&mac)
	}

	//check remote rs ip bind
	if uploadConfig.BindRsIp != "" {
		rsClient.Conn.BindRemoteIp = uploadConfig.BindRsIp
	}

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
		if skip, prefix := hitByPathPrefixes(uploadConfig.SkipPathPrefixes, localFileRelativePath); skip {
			logs.Informational("Skip by path prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, prefix := hitByFilePrefixes(uploadConfig.SkipFilePrefixes, localFileRelativePath); skip {
			logs.Informational("Skip by file prefix `%s` for local file path `%s`", prefix, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, fixedStr := hitByFixesString(uploadConfig.SkipFixedStrings, localFileRelativePath); skip {
			logs.Informational("Skip by fixed string `%s` for local file path `%s`", fixedStr, localFileRelativePath)
			skippedFileCount += 1
			continue
		}

		if skip, suffix := hitBySuffixes(uploadConfig.SkipSuffixes, localFileRelativePath); skip {
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
		needToUpload, isOverwrite := checkFileNeedToUpload(uploadConfig, &rsClient, ldb, &ldbWOpt, ldbKey, localFilePath,
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

			policy := rs.PutPolicy{}
			policy.Scope = uploadConfig.Bucket
			if uploadConfig.Overwrite {
				policy.Scope = fmt.Sprintf("%s:%s", uploadConfig.Bucket, uploadFileKey)
				policy.InsertOnly = 0
			}

			policy.FileType = uploadConfig.FileType

			policy.Expires = 7 * 24 * 3600
			upToken := policy.Token(&mac)

			if localFileSize > putThreshold {
				resumableUploadFile(uploadConfig, transport, ldb, &ldbWOpt, ldbKey, upToken, storePath,
					localFilePath, uploadFileKey, localFileLastModified, isOverwrite, exporter)
			} else {
				formUploadFile(uploadConfig, transport, ldb, &ldbWOpt, ldbKey, upToken,
					localFilePath, uploadFileKey, localFileLastModified, isOverwrite, exporter)
			}
		}
	}

	upWaitGroup.Wait()

	//flush
	exportWriters := []*bufio.Writer{
		exporter.SuccessWriter,
		exporter.FailureWriter,
		exporter.OverwriteWriter,
	}

	for _, writer := range exportWriters {
		if writer != nil {
			writer.Flush()
		}
	}

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

func hitByPathPrefixes(pathPrefixesStr, localFileRelativePath string) (hit bool, pathPrefix string) {
	if pathPrefixesStr != "" {
		//unpack skip prefix
		pathPrefixes := strings.Split(pathPrefixesStr, ",")
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

func hitByFilePrefixes(filePrefixesStr, localFileRelativePath string) (hit bool, filePrefix string) {
	if filePrefixesStr != "" {
		//unpack skip prefix
		filePrefixes := strings.Split(filePrefixesStr, ",")
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

func hitByFixesString(fixedStringsStr, localFileRelativePath string) (hit bool, hitFixedStr string) {
	if fixedStringsStr != "" {
		//unpack fixed strings
		fixedStrings := strings.Split(fixedStringsStr, ",")
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

func hitBySuffixes(suffixesStr, localFileRelativePath string) (hit bool, hitSuffix string) {
	if suffixesStr != "" {
		suffixes := strings.Split(suffixesStr, ",")
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

func checkFileNeedToUpload(uploadConfig *UploadConfig, rsClient *rs.Client, ldb *leveldb.DB, ldbWOpt *opt.WriteOptions,
	ldbKey, localFilePath, uploadFileKey string, localFileLastModified, localFileSize int64) (needToUpload, isOverwrite bool) {
	//default to upload
	needToUpload = true

	//check before to upload
	if uploadConfig.CheckExists {
		rsEntry, checkErr := rsClient.Stat(nil, uploadConfig.Bucket, uploadFileKey)
		if checkErr == nil {
			ldbValue := fmt.Sprintf("%d", localFileLastModified)
			if uploadConfig.CheckHash {
				//compare hash
				localEtag, cErr := GetEtag(localFilePath)
				if cErr != nil {
					logs.Error("File `%s` calc local hash failed, %s", uploadFileKey, cErr)
					atomic.AddInt64(&failureFileCount, 1)
					needToUpload = false
				}
				if rsEntry.Hash == localEtag {
					logs.Informational("File `%s` exists in bucket, hash match, ignore this upload", uploadFileKey)
					atomic.AddInt64(&skippedFileCount, 1)
					putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
					if putErr != nil {
						logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
					}
					needToUpload = false
				} else {
					if !uploadConfig.Overwrite {
						logs.Warning("Skip upload of unmatch hash file `%s` because `overwrite` is false", localFilePath)
						atomic.AddInt64(&notOverwriteCount, 1)
						needToUpload = false
					} else {
						isOverwrite = true
						logs.Informational("File `%s` exists in bucket, but hash not match, go to upload", uploadFileKey)
					}
				}
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
					} else {
						if !uploadConfig.Overwrite {
							logs.Warning("Skip upload of unmatch size file `%s` because `overwrite` is false", localFilePath)
							atomic.AddInt64(&notOverwriteCount, 1)
							needToUpload = false
						} else {
							isOverwrite = true
							logs.Info("File `%s` exists in bucket, but size not match, go to upload", uploadFileKey)
						}
					}
				} else {
					logs.Info("File `%s` exists in bucket, no hash or size check, ignore this upload", uploadFileKey)
					atomic.AddInt64(&skippedFileCount, 1)
					putErr := ldb.Put([]byte(ldbKey), []byte(ldbValue), ldbWOpt)
					if putErr != nil {
						logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, putErr)
					}
					needToUpload = false
				}
			}
		} else {
			if _, ok := checkErr.(*rpc.ErrorInfo); !ok {
				//not logic error, should be network error
				logs.Error("Get file `%s` stat error, %s", uploadFileKey, checkErr)
				atomic.AddInt64(&failureFileCount, 1)
				needToUpload = false
			}
		}
	} else {
		//check leveldb
		ldbFlmd, err := ldb.Get([]byte(ldbKey), nil)
		flmd, _ := strconv.ParseInt(string(ldbFlmd), 10, 64)
		//not exist, return ErrNotFound
		//check last modified

		if err == nil {
			if localFileLastModified == flmd {
				logs.Informational("Skip by local leveldb log for file `%s`", localFilePath)
				atomic.AddInt64(&skippedFileCount, 1)
				needToUpload = false
			} else {
				if !uploadConfig.Overwrite {
					//no overwrite set
					logs.Warning("Skip upload of changed file `%s` because `overwrite` is false",
						localFilePath)
					atomic.AddInt64(&notOverwriteCount, 1)
					needToUpload = false
				} else {
					isOverwrite = true
				}
			}
		}
	}
	return
}

func formUploadFile(uploadConfig *UploadConfig, transport *http.Transport,
	ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string, upToken string,
	localFilePath, uploadFileKey string, localFileLastModified int64, isOverwrite bool, exporter *FileExporter) {
	var putClient rpc.Client
	if transport != nil {
		putClient = rpc.NewClientEx(transport, uploadConfig.BindUpIp)
	} else {
		putClient = rpc.NewClient(uploadConfig.BindUpIp)
	}

	putRet := fio.PutRet{}
	putExtra := fio.PutExtra{
		CheckCrc: 1,
	}

	err := fio.PutFile(putClient, nil, &putRet, upToken, uploadFileKey, localFilePath, &putExtra)
	if err != nil {
		atomic.AddInt64(&failureFileCount, 1)
		var errMsg string
		if pErr, ok := err.(*rpc.ErrorInfo); ok {
			errMsg = pErr.Err
		} else {
			errMsg = err.Error()
		}
		logs.Error("Form upload file `%s` => `%s` failed due to nerror `%s`", localFilePath, uploadFileKey, errMsg)
		if exporter.FailureWriter != nil {
			exporter.FailureLock.Lock()
			exporter.FailureWriter.WriteString(fmt.Sprintf("%s\t%s\t%s\n", localFilePath, uploadFileKey, errMsg))
			exporter.FailureWriter.Flush()
			exporter.FailureLock.Unlock()
		}
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

		if exporter.SuccessWriter != nil {
			exporter.SuccessLock.Lock()
			exporter.SuccessWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
			exporter.SuccessWriter.Flush()
			exporter.SuccessLock.Unlock()
		}
		if isOverwrite && exporter.OverwriteWriter != nil {
			exporter.OverwriteLock.Lock()
			exporter.OverwriteWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
			exporter.OverwriteWriter.Flush()
			exporter.OverwriteLock.Unlock()
		}
	}
}

func resumableUploadFile(uploadConfig *UploadConfig, transport *http.Transport,
	ldb *leveldb.DB, ldbWOpt *opt.WriteOptions, ldbKey string, upToken string, storePath,
	localFilePath, uploadFileKey string, localFileLastModified int64, isOverwrite bool, exporter *FileExporter) {
	var putClient rpc.Client
	if transport != nil {
		putClient = rio.NewClientEx(upToken, transport, uploadConfig.BindUpIp)
	} else {
		putClient = rio.NewClient(upToken, uploadConfig.BindUpIp)
	}

	//params
	putRet := rio.PutRet{}
	putExtra := rio.PutExtra{}

	//progress file
	progressFileKey := Md5Hex(fmt.Sprintf("%s:%s|%s:%s", uploadConfig.SrcDir,
		uploadConfig.Bucket, localFilePath, uploadFileKey))
	progressFilePath := filepath.Join(storePath, fmt.Sprintf("%s.progress", progressFileKey))
	putExtra.ProgressFile = progressFilePath

	//resumable upload
	err := rio.PutFile(putClient, nil, &putRet, uploadFileKey, localFilePath, &putExtra)
	if err != nil {
		os.Remove(progressFilePath)
		atomic.AddInt64(&failureFileCount, 1)
		var errMsg string
		if pErr, ok := err.(*rpc.ErrorInfo); ok {
			errMsg = pErr.Err
		} else {
			errMsg = err.Error()
		}
		logs.Error("Resumable upload file `%s` => `%s` failed due to nerror `%s`", localFilePath, uploadFileKey, errMsg)
		if exporter.FailureWriter != nil {
			exporter.FailureLock.Lock()
			exporter.FailureWriter.WriteString(fmt.Sprintf("%s\t%s\t%s\n", localFilePath, uploadFileKey, errMsg))
			exporter.FailureWriter.Flush()
			exporter.FailureLock.Unlock()
		}
	} else {
		os.Remove(progressFilePath)
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

		if exporter.SuccessWriter != nil {
			exporter.SuccessLock.Lock()
			exporter.SuccessWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
			exporter.SuccessWriter.Flush()
			exporter.SuccessLock.Unlock()
		}
		if isOverwrite && exporter.OverwriteWriter != nil {
			exporter.OverwriteLock.Lock()
			exporter.OverwriteWriter.WriteString(fmt.Sprintf("%s\t%s\n", localFilePath, uploadFileKey))
			exporter.OverwriteWriter.Flush()
			exporter.OverwriteLock.Unlock()
		}
	}
}
