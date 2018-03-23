package qshell

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
{
	"dest_dir"		:	"/Users/jemy/Backup",
	"bucket"		:	"test-bucket",
	"prefix"		:	"demo/",
	"suffixes"		: 	".png,.jpg",
	"range_bytes"   :   4096,
	"range_worker"  :   10
}
*/

const (
	MIN_DOWNLOAD_THREAD_COUNT = 1
	MAX_DOWNLOAD_THREAD_COUNT = 2000
)

const (
	DEFAULT_RANGE_BLOCK  = 4 * 1024 * 1024
	DEFAULT_RANGE_WORKER = 1
)

type DownloadConfig struct {
	DestDir  string `json:"dest_dir"`
	Bucket   string `json:"bucket"`
	Prefix   string `json:"prefix,omitempty"`
	Suffixes string `json:"suffixes,omitempty"`
	//down from cdn
	Referer   string `json:"referer,omitempty"`
	CdnDomain string `json:"cdn_domain,omitempty"`
	//log settings
	LogLevel  string `json:"log_level,omitempty"`
	LogFile   string `json:"log_file,omitempty"`
	LogRotate int    `json:"log_rotate,omitempty"`
	LogStdout bool   `json:"log_stdout,omitempty"`

	IsHostFileSpecified bool `json:"-"`

	//range download
	RangeBytes  int `json:"range_bytes,omitempty"`
	RangeWorker int `json:"range_worker,omitempty"`
}

func doDownload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func QiniuDownload(threadCount int, downConfig *DownloadConfig) {
	timeStart := time.Now()
	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", downConfig.DestDir, downConfig.Bucket))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "qdownload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		os.Exit(STATUS_ERROR)
	}

	//init log settings
	defaultLogFile := filepath.Join(storePath, fmt.Sprintf("%s.log", jobId))
	//init log level
	logLevel := logs.LevelInfo
	logRotate := 1
	if downConfig.LogRotate > 0 {
		logRotate = downConfig.LogRotate
	}
	switch downConfig.LogLevel {
	case "debug":
		logLevel = logs.LevelDebug
	case "info":
		logLevel = logs.LevelInfo
	case "warn":
		logLevel = logs.LevelWarning
	case "error":
		logLevel = logs.LevelError
	default:
		logLevel = logs.LevelInfo
	}

	//init log writer
	if downConfig.LogFile == "" {
		//set default log file
		downConfig.LogFile = defaultLogFile
	}

	if !downConfig.LogStdout {
		logs.GetBeeLogger().DelLogger(logs.AdapterConsole)
	}
	//open log file
	fmt.Println("Writing download log to file", downConfig.LogFile)

	//daily rotate
	logCfg := BeeLogConfig{
		Filename: downConfig.LogFile,
		Level:    logLevel,
		Daily:    true,
		MaxDays:  logRotate,
	}
	logs.SetLogger(logs.AdapterFile, logCfg.ToJson())
	fmt.Println()

	logs.Info("Load account from %s", filepath.Join(QShellRootPath, ".qshell/account.json"))
	account, gErr := GetAccount()
	if gErr != nil {
		fmt.Println("Get account error,", gErr)
		os.Exit(STATUS_ERROR)
	}
	mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}

	var domainOfBucket string

	//get bucket zone info
	bucketInfo, gErr := GetBucketInfo(&mac, downConfig.Bucket)
	if gErr != nil {
		logs.Error("Get bucket region info error,", gErr)
		os.Exit(STATUS_ERROR)
	}
	//get domains of bucket
	domainsOfBucket, gErr := GetDomainsOfBucket(&mac, downConfig.Bucket)
	if gErr != nil {
		logs.Error("Get domains of bucket error,", gErr)
		os.Exit(STATUS_ERROR)
	}

	if len(domainsOfBucket) == 0 {
		logs.Error("No domains found for bucket", downConfig.Bucket)
		os.Exit(STATUS_ERROR)
	}

	for _, domain := range domainsOfBucket {
		if !strings.HasPrefix(domain, ".") {
			domainOfBucket = domain
			break
		}
	}

	if !downConfig.IsHostFileSpecified {
		//set up host
		SetZone(bucketInfo.Region)
	}

	//set proxy
	ioProxyAddress := conf.IO_HOST
	//check whether cdn domain is set
	if downConfig.CdnDomain != "" {
		ioProxyAddress = downConfig.CdnDomain
	}

	//trim http and https prefix
	ioProxyAddress = strings.TrimPrefix(ioProxyAddress, "http://")
	ioProxyAddress = strings.TrimPrefix(ioProxyAddress, "https://")
	if downConfig.CdnDomain != "" {
		domainOfBucket = ioProxyAddress
	}

	if domainOfBucket == "" {
		logs.Error("No domains found to download files")
		os.Exit(STATUS_ERROR)
	}

	jobListFileName := filepath.Join(storePath, fmt.Sprintf("%s.list", jobId))
	resumeFile := filepath.Join(storePath, fmt.Sprintf("%s.ldb", jobId))
	resumeLevelDb, openErr := leveldb.OpenFile(resumeFile, nil)
	if openErr != nil {
		logs.Error("Open resume record leveldb error", openErr)
		os.Exit(STATUS_ERROR)
	}
	defer resumeLevelDb.Close()
	//sync underlying writes from the OS buffer cache
	//through to actual disk
	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	//list bucket, prepare file list to download
	logs.Info("Listing bucket `%s` by prefix `%s`", downConfig.Bucket, downConfig.Prefix)
	listErr := ListBucket(&mac, downConfig.Bucket, downConfig.Prefix, "", jobListFileName)
	if listErr != nil {
		logs.Error("List bucket error", listErr)
		os.Exit(STATUS_ERROR)
	}

	//init wait group
	var downWaitGroup = sync.WaitGroup{}
	var downloadTasks chan func()
	var initDownOnce sync.Once

	initDownOnce.Do(func() {
		downloadTasks = make(chan func(), threadCount)
		for i := 0; i < threadCount; i++ {
			go doDownload(downloadTasks)
		}
	})

	//init counters
	var totalFileCount int64
	var currentFileCount int64
	var existsFileCount int64
	var updateFileCount int64
	var successFileCount int64
	var failureFileCount int64
	var skipBySuffixes int64

	totalFileCount = GetFileLineCount(jobListFileName)

	//open prepared file list to download files
	listFp, openErr := os.Open(jobListFileName)
	if openErr != nil {
		logs.Error("Open list file error", openErr)
		os.Exit(STATUS_ERROR)
	}
	defer listFp.Close()

	listScanner := bufio.NewScanner(listFp)
	listScanner.Split(bufio.ScanLines)
	//key, fsize, etag, lmd, mime, enduser

	downSuffixes := strings.Split(downConfig.Suffixes, ",")
	filterSuffixes := make([]string, 0, len(downSuffixes))

	for _, suffix := range downSuffixes {
		if strings.TrimSpace(suffix) != "" {
			filterSuffixes = append(filterSuffixes, suffix)
		}
	}

	for listScanner.Scan() {
		currentFileCount += 1
		line := strings.TrimSpace(listScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) >= 4 {
			fileKey := items[0]

			if len(filterSuffixes) > 0 {
				//filter files by suffixes
				var goAhead bool
				for _, suffix := range filterSuffixes {
					if strings.HasSuffix(fileKey, suffix) {
						goAhead = true
						break
					}
				}

				if !goAhead {
					skipBySuffixes += 1
					logs.Info("Skip download `%s`, suffix filter not match", fileKey)
					continue
				}
			}

			remoteFileSize, pErr := strconv.ParseInt(items[1], 10, 64)
			if pErr != nil {
				logs.Error("Invalid list line", line)
				continue
			}

			remoteFileMtime, pErr := strconv.ParseInt(items[3], 10, 64)
			if pErr != nil {
				logs.Error("Invalid list line", line)
				continue
			}

			fileUrl := makePrivateDownloadLink(&mac, domainOfBucket, ioProxyAddress, fileKey)

			//progress
			if totalFileCount != 0 {
				fmt.Printf("Downloading %s [%d/%d, %.1f%%] ...\n", fileKey, currentFileCount, totalFileCount,
					float32(currentFileCount)*100/float32(totalFileCount))
			} else {
				fmt.Printf("Downloading %s ...\n", fileKey)
			}
			//check whether log file exists
			localFilePath := filepath.Join(downConfig.DestDir, fileKey)
			localAbsFilePath, _ := filepath.Abs(localFilePath)

			//check whether to download new files
			rKey := localAbsFilePath
			rVal := fmt.Sprintf("%d|%d", remoteFileMtime, remoteFileSize)

			var updateOldFile bool

			localFileInfo, statErr := os.Stat(localFilePath)
			if statErr == nil {
				//log file exists, check whether have updates
				oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePath), nil)
				if notFoundErr == nil {
					//if exists
					oldFileInfoItems := strings.Split(string(oldFileInfo), "|")
					oldFileLmd, _ := strconv.ParseInt(oldFileInfoItems[0], 10, 64)
					//oldFileSize, _ := strconv.ParseInt(oldFileInfoItems[1], 10, 64)
					if oldFileLmd == remoteFileMtime && localFileInfo.Size() == remoteFileSize {
						//nothing change, ignore
						logs.Info("Local file `%s` exists, same as in bucket, download skip", localAbsFilePath)
						existsFileCount += 1
						continue
					} else {
						//somthing changed, must download a new file
						logs.Info("Local file `%s` exists, but remote file changed, go to download", localAbsFilePath)
						updateOldFile = true
					}
				} else {
					if localFileInfo.Size() != remoteFileSize {
						logs.Info("Local file `%s` exists, size not the same as in bucket, go to download", localAbsFilePath)
						updateOldFile = true
					} else {
						//treat the local file not changed, write to leveldb, though may not accurate
						//nothing to do
						logs.Warning("Local file `%s` exists with same size as `%s`, treat it not changed", localAbsFilePath, fileKey)
						atomic.AddInt64(&existsFileCount, 1)
						continue
					}
				}
			}

			//set file info in leveldb
			resumeLevelDb.Put([]byte(rKey), []byte(rVal), &ldbWOpt)

			//download new
			downWaitGroup.Add(1)
			downloadTasks <- func() {
				defer downWaitGroup.Done()

				downErr := downloadFile(downConfig, fileKey, fileUrl, domainOfBucket, remoteFileSize)
				if downErr != nil {
					atomic.AddInt64(&failureFileCount, 1)
				} else {
					atomic.AddInt64(&successFileCount, 1)
					if updateOldFile {
						atomic.AddInt64(&updateFileCount, 1)
					}
				}
			}
		}
	}

	//wait for all tasks done
	downWaitGroup.Wait()

	logs.Info("-------Download Result-------")
	logs.Info("%10s%10d", "Total:", totalFileCount)
	logs.Info("%10s%10d", "Skipped:", skipBySuffixes)
	logs.Info("%10s%10d", "Exists:", existsFileCount)
	logs.Info("%10s%10d", "Success:", successFileCount)
	logs.Info("%10s%10d", "Update:", updateFileCount)
	logs.Info("%10s%10d", "Failure:", failureFileCount)
	logs.Info("%10s%15s", "Duration:", time.Since(timeStart))
	logs.Info("-----------------------------")
	fmt.Println("\nSee download log at path", downConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(STATUS_ERROR)
	}
}

/*
@param ioProxyHost - like http://iovip.qbox.me
*/
func makePrivateDownloadLink(mac *digest.Mac, domainOfBucket, ioProxyAddress, fileKey string) (fileUrl string) {
	publicUrl := fmt.Sprintf("http://%s/%s", domainOfBucket, url.PathEscape(fileKey))
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	privateUrl, _ := PrivateUrl(mac, publicUrl, deadline)

	//replace the io proxy host
	fileUrl = strings.Replace(privateUrl, domainOfBucket, ioProxyAddress, -1)
	return
}

//file key -> mtime
func downloadFile(downConfig *DownloadConfig, fileName, fileUrl, domainOfBucket string, remoteFileSize int64) (err error) {
	startDown := time.Now().Unix()
	destDir := downConfig.DestDir
	localFilePath := filepath.Join(destDir, fileName)
	localFileDir := filepath.Dir(localFilePath)
	localFilePathTmp := fmt.Sprintf("%s.tmp", localFilePath)

	mkdirErr := os.MkdirAll(localFileDir, 0775)
	if mkdirErr != nil {
		err = mkdirErr
		logs.Error("MkdirAll failed for", localFileDir, mkdirErr)
		return
	}

	logs.Info("Downloading", fileName, "=>", localFilePath)
	//download worker
	var downWaitGroup = sync.WaitGroup{}
	var downloadTasks chan func()
	var initDownOnce sync.Once

	initDownOnce.Do(func() {
		downloadTasks = make(chan func(), downConfig.RangeWorker)
		for i := 0; i < downConfig.RangeWorker; i++ {
			go doDownload(downloadTasks)
		}
	})

	//create temp file
	localTempFile, createErr := os.Create(localFilePathTmp)
	if createErr != nil {
		err = fmt.Errorf("Create temp file for", fileName, "failed when opening,", createErr)
		return
	}
	defer localTempFile.Close()

	//try range download
	rangeBlockSize := int64(downConfig.RangeBytes)
	var totalTasks int64 = remoteFileSize / rangeBlockSize
	if remoteFileSize%int64(downConfig.RangeBytes) != 0 {
		totalTasks += 1
	}

	downErrs := sync.Map{}

	var i int64
	for i = 0; i < totalTasks; i++ {
		rangeIndex := i
		var rangeStart int64 = rangeIndex * rangeBlockSize
		var rangeEnd int64 = (rangeIndex+1)*rangeBlockSize - 1

		if rangeEnd >= remoteFileSize {
			rangeEnd = remoteFileSize - 1
		}

		downWaitGroup.Add(1)
		downloadTasks <- func() {
			defer downWaitGroup.Done()

			//try each download
			downErr := func() (err error) {
				req, reqErr := http.NewRequest("GET", fileUrl, nil)
				if reqErr != nil {
					err = reqErr
					return
				}

				//set host
				req.Host = domainOfBucket
				//set referer
				if downConfig.Referer != "" {
					req.Header.Add("Referer", downConfig.Referer)
				}
				//set range
				req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
				resp, respErr := http.DefaultClient.Do(req)
				if respErr != nil {
					err = respErr
					return
				}
				defer resp.Body.Close()

				//read into memory
				bodyBytes, readErr := ioutil.ReadAll(resp.Body)
				if readErr != nil {
					err = readErr
					return
				}

				//Clients of WriteAt can execute parallel WriteAt calls on the same destination
				//if the ranges do not overlap.
				wLen, wErr := localTempFile.WriteAt(bodyBytes, rangeStart)
				if wErr != nil {
					err = wErr
					return
				}

				expectLen := rangeEnd - rangeStart + 1
				if wLen != int(expectLen) {
					err = fmt.Errorf("Write range block %d not full error, expect %d, but %d bytes",
						rangeIndex, expectLen, wLen)
					return
				}

				return
			}()

			if downErr != nil {
				downErrs.Store(rangeIndex, downErr)
			}
		}
	}

	downWaitGroup.Wait()

	//check down errs
	for i = 0; i < totalTasks; i++ {
		downErrVal, _ := downErrs.Load(i)
		if downErr, _ := downErrVal.(error); downErr != nil {
			err = downErr
			logs.Error("Download temp file for %s => %s error, %s", fileName, localFilePath, downErr)
			return
		}
	}

	//sync
	fErr := localTempFile.Sync()
	if fErr != nil {
		err = fErr
		logs.Error("Sync temp file for %s => %s error, %s", fileName, localFilePath, fErr)
		return
	}

	cErr := localTempFile.Close()
	if cErr != nil {
		err = cErr
		logs.Error("Close temp file for %s => %s error, %s", fileName, localFilePath, fErr)
		return
	}

	rErr := os.Rename(localFilePathTmp, localFilePath)
	if rErr != nil {
		err = rErr
		logs.Error("Rename temp file for %s => %s error, %s", fileName, localFilePath, rErr)
	}

	endDown := time.Now().Unix()
	avgSpeed := fmt.Sprintf("%.2fKB/s", float64(remoteFileSize)/float64(endDown-startDown)/1024)
	logs.Info("Download", fileName, "=>", localFilePath, "success", avgSpeed)

	return
}
