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
	"file_list"             :   "",
	"bucket"		:	"test-bucket",
	"put_threshold"		:	10000000,
	"key_prefix"		:	"2014/12/01/",
	"ignore_dir"		:	false,
	"overwrite"		:	false,
	"check_exists"		:	true,
	"skip_file_prefixes"	:	"IMG_",
	"skip_path_prefixes"	:	"tmp/,bin/,obj/",
	"skip_suffixes"		:	".exe,.obj,.class",
	"skip_fixed_strings"    :   ".svn,.git",
	"up_host"		:	"http://upload.qiniu.com",
	"bind_up_ip"		:	"",
	"bind_rs_ip"		:	"",
	"bind_nic_ip"		:	"",
	"rescan_local"		:	false
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
}

var upSettings = rio.Settings{
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

func QiniuUpload(threadCount int, uploadConfig *UploadConfig) {
	timeStart := time.Now()
	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "qupload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		return
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
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s", "level":%d, "daily":true, "maxdays":%d}`,
		uploadConfig.LogFile, logLevel, logRotate))
	fmt.Println()

	//global up settings
	logs.Info("Load account from %s", filepath.Join(QShellRootPath, ".qshell/account.json"))
	account, gErr := GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		return
	}
	mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}
	//get bucket zone info
	bucketInfo, gErr := GetBucketInfo(&mac, uploadConfig.Bucket)
	if gErr != nil {
		logs.Error("Get bucket region info error,", gErr)
		return
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
		conf.UP_HOST = uploadConfig.UpHost
	}
	//set resume upload settings
	rio.SetSettings(&upSettings)

	//make SrcDir the full path
	uploadConfig.SrcDir, _ = filepath.Abs(uploadConfig.SrcDir)

	//find the local file list, by specified or by config
	var cacheResultName string
	var totalFileCount int64
	var cacheErr error
	_, fStatErr := os.Stat(uploadConfig.FileList)
	if uploadConfig.FileList != "" && fStatErr == nil {
		//use specified file list
		cacheResultName = uploadConfig.FileList
		totalFileCount = GetFileLineCount(cacheResultName)
	} else {
		//cache file
		cacheResultName = filepath.Join(storePath, fmt.Sprintf("%s.cache", jobId))
		cacheTempName := filepath.Join(storePath, fmt.Sprintf("%s.cache.temp", jobId))
		cacheCountName := filepath.Join(storePath, fmt.Sprintf("%s.count", jobId))

		rescanLocalDir := false
		if _, statErr := os.Stat(cacheResultName); statErr == nil {
			//file exists
			rescanLocalDir = uploadConfig.RescanLocal
		} else {
			rescanLocalDir = true
		}

		if rescanLocalDir {
			logs.Informational("Listing local sync dir, this can take a long time for big directory, please wait patiently")
			totalFileCount, cacheErr = DirCache(uploadConfig.SrcDir, cacheTempName)
			if cacheErr != nil {
				return
			}

			if rErr := os.Rename(cacheTempName, cacheResultName); rErr != nil {
				logs.Error("Rename the temp cached file error", rErr)
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
	}

	//leveldb folder
	leveldbFileName := filepath.Join(storePath, jobId+".ldb")
	ldb, err := leveldb.OpenFile(leveldbFileName, nil)
	if err != nil {
		logs.Error("Open leveldb `%s` failed due to %s", leveldbFileName, err)
		return
	}
	defer ldb.Close()
	//sync
	cacheResultFileHandle, err := os.Open(cacheResultName)
	if err != nil {
		logs.Error("Open list file `%s` failed due to %s", cacheResultName, err)
		return
	}
	defer cacheResultFileHandle.Close()
	bScanner := bufio.NewScanner(cacheResultFileHandle)
	bScanner.Split(bufio.ScanLines)

	var currentFileCount int64
	var successFileCount int64
	var notOverwriteCount int64
	var failureFileCount int64
	var skippedFileCount int64

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

		localRelFpath := items[0]
		currentFileCount += 1

		skip := false
		//check skip local file or folder
		if uploadConfig.SkipPathPrefixes != "" {
			//unpack skip prefix
			skipPathPrefixes := strings.Split(uploadConfig.SkipPathPrefixes, ",")
			for _, prefix := range skipPathPrefixes {
				if strings.TrimSpace(prefix) == "" {
					continue
				}

				if strings.HasPrefix(localRelFpath, prefix) {
					logs.Informational("Skip by path prefix `%s` for local file path `%s`", prefix, localRelFpath)
					skip = true
					skippedFileCount += 1
					break
				}
			}

			if skip {
				continue
			}
		}

		if uploadConfig.SkipFilePrefixes != "" {
			//unpack skip prefix
			skipFilePrefixes := strings.Split(uploadConfig.SkipFilePrefixes, ",")
			for _, prefix := range skipFilePrefixes {
				if strings.TrimSpace(prefix) == "" {
					continue
				}

				localFname := filepath.Base(localRelFpath)
				if strings.HasPrefix(localFname, prefix) {
					logs.Informational("Skip by file prefix `%s` for local file path `%s`", prefix, localRelFpath)
					skip = true
					skippedFileCount += 1
					break
				}
			}

			if skip {
				continue
			}
		}

		if uploadConfig.SkipFixedStrings != "" {
			//unpack fixed strings
			skipFixedStrings := strings.Split(uploadConfig.SkipFixedStrings, ",")
			for _, fixedStr := range skipFixedStrings {
				if strings.TrimSpace(fixedStr) == "" {
					continue
				}

				if strings.Contains(localRelFpath, fixedStr) {
					logs.Informational("Skip by fixed string `%s` for local file path `%s`", fixedStr, localRelFpath)
					skip = true
					skippedFileCount += 1
					break
				}
			}

			if skip {
				continue
			}
		}

		if uploadConfig.SkipSuffixes != "" {
			skipSuffixes := strings.Split(uploadConfig.SkipSuffixes, ",")
			for _, suffix := range skipSuffixes {
				if strings.TrimSpace(suffix) == "" {
					continue
				}

				if strings.HasSuffix(localRelFpath, suffix) {
					logs.Informational("Skip by suffix `%s` for local file `%s`", suffix, localRelFpath)
					skip = true
					skippedFileCount += 1
					break
				}
			}

			if skip {
				continue
			}
		}

		//pack the upload file key
		localFlmd, _ := strconv.ParseInt(items[2], 10, 64)
		uploadFileKey := localRelFpath

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

		localFilePath := filepath.Join(uploadConfig.SrcDir, localRelFpath)
		fstat, err := os.Stat(localFilePath)
		if err != nil {
			failureFileCount += 1
			logs.Error("Error stat local file `%s` due to `%s`", localFilePath, err)
			continue
		}

		fsize := fstat.Size()
		ldbKey := fmt.Sprintf("%s => %s", localFilePath, uploadFileKey)

		if totalFileCount != 0 {
			fmt.Printf("Uploading %s [%d/%d, %.1f%%] ...\n", ldbKey, currentFileCount, totalFileCount,
				float32(currentFileCount)*100/float32(totalFileCount))
		} else {
			fmt.Printf("Uploading %s ...\n", ldbKey)
		}

		//check exists
		if uploadConfig.CheckExists {
			rsEntry, checkErr := rsClient.Stat(nil, uploadConfig.Bucket, uploadFileKey)
			if checkErr == nil {
				if uploadConfig.CheckHash {
					//compare hash
					localEtag, cErr := GetEtag(localFilePath)
					if cErr != nil {
						atomic.AddInt64(&failureFileCount, 1)
						logs.Error("File `%s` calc local hash failed, %s", uploadFileKey, cErr)
						continue
					}
					if rsEntry.Hash == localEtag {
						atomic.AddInt64(&skippedFileCount, 1)
						perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
						if perr != nil {
							logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, perr)
						}
						logs.Informational("File `%s` exists in bucket, hash match, ignore this upload", uploadFileKey)
						continue
					} else {
						logs.Informational("File `%s` exists in bucket, but hash not match, go to upload", uploadFileKey)
					}
				} else {
					if uploadConfig.CheckSize {
						if rsEntry.Fsize == fsize {
							atomic.AddInt64(&skippedFileCount, 1)
							perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
							if perr != nil {
								logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, perr)
							}
							logs.Info("File `%s` exists in bucket, size match, ignore this upload", uploadFileKey)
							continue
						} else {
							logs.Info("File `%s` exists in bucket, but size not match, go to upload", uploadFileKey)
						}
					} else {
						atomic.AddInt64(&skippedFileCount, 1)
						perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
						if perr != nil {
							logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, perr)
						}
						logs.Info("File `%s` exists in bucket, no hash or size check, ignore this upload", uploadFileKey)
						continue
					}
				}
			} else {
				if _, ok := checkErr.(*rpc.ErrorInfo); !ok {
					//not logic error, should be network error
					logs.Error("Get file `%s` stat error, %s", uploadFileKey, checkErr)
					atomic.AddInt64(&failureFileCount, 1)
					continue
				}
			}
		} else {
			//check leveldb
			ldbFlmd, err := ldb.Get([]byte(ldbKey), nil)
			flmd, _ := strconv.ParseInt(string(ldbFlmd), 10, 64)
			//not exist, return ErrNotFound
			//check last modified

			if err == nil {
				if localFlmd == flmd {
					logs.Informational("Skip by local leveldb log for file `%s`", localRelFpath)
					atomic.AddInt64(&skippedFileCount, 1)
					continue
				} else {
					if !uploadConfig.Overwrite {
						//no overwrite set
						logs.Warning("Skip upload of changed file `%s` because not set to overwrite", localRelFpath)
						atomic.AddInt64(&notOverwriteCount, 1)
						continue
					}
				}
			}
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
			policy.Expires = 7 * 24 * 3600
			uptoken := policy.Token(&mac)
			if fsize > putThreshold {
				var putClient rpc.Client
				if transport != nil {
					putClient = rio.NewClientEx(uptoken, transport, uploadConfig.BindUpIp)
				} else {
					putClient = rio.NewClient(uptoken, uploadConfig.BindUpIp)
				}

				putRet := rio.PutRet{}
				putExtra := rio.PutExtra{}
				progressFkey := Md5Hex(fmt.Sprintf("%s:%s|%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket, localRelFpath, uploadFileKey))
				progressFname := fmt.Sprintf("%s.progress", progressFkey)
				progressFpath := filepath.Join(storePath, progressFname)
				putExtra.ProgressFile = progressFpath

				err := rio.PutFile(putClient, nil, &putRet, uploadFileKey, localFilePath, &putExtra)
				if err != nil {
					os.Remove(progressFpath)
					atomic.AddInt64(&failureFileCount, 1)
					if pErr, ok := err.(*rpc.ErrorInfo); ok {
						logs.Error("Upload file `%s` => `%s` failed due to `%s`", localFilePath, uploadFileKey, pErr.Err)
					} else {
						logs.Error("Upload file `%s` => `%s` failed due to `%s`", localFilePath, uploadFileKey, err)
					}
				} else {
					os.Remove(progressFpath)
					atomic.AddInt64(&successFileCount, 1)
					logs.Informational("Upload file `%s` => `%s` success", localFilePath, uploadFileKey)
					perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
					if perr != nil {
						logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, perr)
					}
				}
			} else {
				var putClient rpc.Client
				if transport != nil {
					putClient = rpc.NewClientEx(transport, uploadConfig.BindUpIp)
				} else {
					putClient = rpc.NewClient(uploadConfig.BindUpIp)
				}

				putRet := fio.PutRet{}
				err := fio.PutFile(putClient, nil, &putRet, uptoken, uploadFileKey, localFilePath, nil)
				if err != nil {
					atomic.AddInt64(&failureFileCount, 1)
					if pErr, ok := err.(*rpc.ErrorInfo); ok {
						logs.Error("Upload file `%s` => `%s` failed due to `%s`", localFilePath, uploadFileKey, pErr.Err)
					} else {
						logs.Error("Upload file `%s` => `%s` failed due to `%s`", localFilePath, uploadFileKey, err)
					}
				} else {
					atomic.AddInt64(&successFileCount, 1)
					logs.Informational("Upload file `%s` => `%s` success", localFilePath, uploadFileKey)
					perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
					if perr != nil {
						logs.Error("Put key `%s` into leveldb error due to `%s`", ldbKey, perr)
					}
				}
			}
		}
	}

	upWaitGroup.Wait()

	logs.Informational("-------------Upload Result--------------")
	logs.Informational("%20s%10d", "Total:", totalFileCount)
	logs.Informational("%20s%10d", "Success:", successFileCount)
	logs.Informational("%20s%10d", "Failure:", failureFileCount)
	logs.Informational("%20s%10d", "NotOverwrite:", notOverwriteCount)
	logs.Informational("%20s%10d", "Skipped:", skippedFileCount)
	logs.Informational("%20s%15s", "Duration:", time.Since(timeStart))
	logs.Info("----------------------------------------")
	fmt.Println("\nSee upload log at path", uploadConfig.LogFile)

}
