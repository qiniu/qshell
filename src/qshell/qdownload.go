package qshell

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
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
}
*/

const (
	MIN_DOWNLOAD_THREAD_COUNT = 1
	MAX_DOWNLOAD_THREAD_COUNT = 2000
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
	LogLevel     string `json:"log_level,omitempty"`
	LogFile      string `json:"log_file,omitempty"`
	LogRotate    int    `json:"log_rotate,omitempty"`
	LogStdout    bool   `json:"log_stdout,omitempty"`
	FileEncoding string `json:"file_encoding,omitempty"`

	IsHostFileSpecified bool `json:"-"`
}

var downloadTasks chan func()
var initDownOnce sync.Once

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
	storePath := filepath.Join(QShellRootPath, ".qshell", QAccountName, "qdownload", jobId)
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

	for _, d := range domainsOfBucket {
		if !strings.HasPrefix(d.Domain, ".") {
			domainOfBucket = d.Domain
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
	downWaitGroup := sync.WaitGroup{}

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

			fileSize, pErr := strconv.ParseInt(items[1], 10, 64)
			if pErr != nil {
				logs.Error("Invalid list line", line)
				continue
			}

			fileMtime, pErr := strconv.ParseInt(items[3], 10, 64)
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
			localFilePathTemp := fmt.Sprintf("%s.tmp", localFilePath)
			//make the absolute path
			localAbsFilePath, _ := filepath.Abs(localFilePath)
			localAbsFilePathTemp, _ := filepath.Abs(localFilePathTemp)

			//create the path to check
			localFilePathToCheck := localAbsFilePath
			localFilePathTempToCheck := localAbsFilePathTemp

			//add check for gbk file encoding for windows
			if strings.ToLower(downConfig.FileEncoding) == "gbk" {
				localFilePathToCheck, _ = utf82GBK(localAbsFilePath)
				localFilePathTempToCheck, _ = utf82GBK(localAbsFilePathTemp)
			}

			localFileInfo, statErr := os.Stat(localFilePathToCheck)
			var downNewFile bool
			var fromBytes int64

			if statErr == nil {
				//log file exists, check whether have updates
				oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePathToCheck), nil)
				if notFoundErr == nil {
					//if exists
					oldFileInfoItems := strings.Split(string(oldFileInfo), "|")
					oldFileLmd, _ := strconv.ParseInt(oldFileInfoItems[0], 10, 64)
					//oldFileSize, _ := strconv.ParseInt(oldFileInfoItems[1], 10, 64)
					if oldFileLmd == fileMtime && localFileInfo.Size() == fileSize {
						//nothing change, ignore
						logs.Info("Local file `%s` exists, same as in bucket, download skip", localAbsFilePath)
						existsFileCount += 1
						continue
					} else {
						//somthing changed, must download a new file
						logs.Info("Local file `%s` exists, but remote file changed, go to download", localAbsFilePath)
						downNewFile = true
					}
				} else {
					if localFileInfo.Size() != fileSize {
						logs.Info("Local file `%s` exists, size not the same as in bucket, go to download", localAbsFilePath)
						downNewFile = true
					} else {
						//treat the local file not changed, write to leveldb, though may not accurate
						//nothing to do
						logs.Warning("Local file `%s` exists with same size as `%s`, treat it not changed", localAbsFilePath, fileKey)
						atomic.AddInt64(&existsFileCount, 1)
						continue
					}
				}
			} else {
				//check whether tmp file exists
				localTmpFileInfo, statErr := os.Stat(localFilePathTempToCheck)
				if statErr == nil {
					//if tmp file exists, check whether last modify changed
					oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePathToCheck), nil)
					if notFoundErr == nil {
						//if exists
						oldFileInfoItems := strings.Split(string(oldFileInfo), "|")
						oldFileLmd, _ := strconv.ParseInt(oldFileInfoItems[0], 10, 64)
						//oldFileSize, _ := strconv.ParseInt(oldFileInfoItems[1], 10, 64)
						if oldFileLmd == fileMtime {
							//tmp file exists, file not changed, use range to download
							if localTmpFileInfo.Size() < fileSize {
								fromBytes = localTmpFileInfo.Size()
							} else {
								//rename it
								renameErr := os.Rename(localFilePathTempToCheck, localFilePathToCheck)
								if renameErr != nil {
									logs.Error("Rename temp file `%s` to final file `%s` error", localAbsFilePathTemp, localAbsFilePath,
										renameErr)
								}
								continue
							}
						} else {
							logs.Info("Local tmp file `%s` exists, but remote file changed, go to download", localAbsFilePathTemp)
							downNewFile = true
						}
					} else {
						//log tmp file exists, but no record in leveldb, download a new file
						logs.Info("Local tmp file `%s` exists, but no record in leveldb ,go to download", localAbsFilePathTemp)
						downNewFile = true
					}
				} else {
					//no file exists, donwload a new file
					downNewFile = true
				}
			}

			//set file info in leveldb
			rKey := localFilePathToCheck
			rVal := fmt.Sprintf("%d|%d", fileMtime, fileSize)
			resumeLevelDb.Put([]byte(rKey), []byte(rVal), &ldbWOpt)

			//download new
			downWaitGroup.Add(1)
			downloadTasks <- func() {
				defer downWaitGroup.Done()

				downErr := downloadFile(downConfig, fileKey, fileUrl, domainOfBucket, fileSize, fromBytes)
				if downErr != nil {
					atomic.AddInt64(&failureFileCount, 1)
				} else {
					atomic.AddInt64(&successFileCount, 1)
					if !downNewFile {
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
func downloadFile(downConfig *DownloadConfig, fileKey, fileUrl, domainOfBucket string, fileSize int64,
	fromBytes int64) (err error) {
	startDown := time.Now().Unix()
	destDir := downConfig.DestDir
	localFilePath := filepath.Join(destDir, fileKey)
	localFileDir := filepath.Dir(localFilePath)
	localFilePathTemp := fmt.Sprintf("%s.tmp", localFilePath)

	//make the absolute path
	localAbsFilePath, _ := filepath.Abs(localFilePath)
	localAbsFilePathTemp, _ := filepath.Abs(localFilePathTemp)
	localAbsFileDir, _ := filepath.Abs(localFileDir)

	localFilePathTarget := localAbsFilePath
	localFilePathTempTarget := localAbsFilePathTemp
	localFileDirTarget := localAbsFileDir

	//add check for gbk file encoding for windows
	if strings.ToLower(downConfig.FileEncoding) == "gbk" {
		localFilePathTarget, _ = utf82GBK(localAbsFilePath)
		localFilePathTempTarget, _ = utf82GBK(localAbsFilePathTemp)
		localFileDirTarget, _ = utf82GBK(localAbsFileDir)
	}

	mkdirErr := os.MkdirAll(localFileDirTarget, 0775)
	if mkdirErr != nil {
		err = mkdirErr
		logs.Error("MkdirAll failed for", localFileDir, mkdirErr)
		return
	}

	logs.Info("Downloading", fileKey, "=>", localAbsFilePath)
	//new request
	req, reqErr := http.NewRequest("GET", fileUrl, nil)
	if reqErr != nil {
		err = reqErr
		logs.Info("New request", fileKey, "failed by url", fileUrl, reqErr)
		return
	}
	//set host
	req.Host = domainOfBucket
	if downConfig.Referer != "" {
		req.Header.Add("Referer", downConfig.Referer)
	}

	if fromBytes != 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", fromBytes))
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = respErr
		logs.Info("Download", fileKey, "failed by url", fileUrl, respErr)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		var localFp *os.File
		var openErr error
		if fromBytes != 0 {
			localFp, openErr = os.OpenFile(localFilePathTempTarget, os.O_APPEND|os.O_WRONLY, 0655)
		} else {
			localFp, openErr = os.Create(localFilePathTempTarget)
		}

		if openErr != nil {
			err = openErr
			logs.Error("Open local file", localAbsFilePathTemp, "failed", openErr)
			return
		}

		cpCnt, cpErr := io.Copy(localFp, resp.Body)
		if cpErr != nil {
			err = cpErr
			localFp.Close()
			logs.Error("Download", fileKey, "failed", cpErr)
			return
		}
		localFp.Close()

		endDown := time.Now().Unix()
		avgSpeed := fmt.Sprintf("%.2fKB/s", float64(cpCnt)/float64(endDown-startDown)/1024)

		//move temp file to log file
		renameErr := os.Rename(localFilePathTempTarget, localFilePathTarget)
		if renameErr != nil {
			err = renameErr
			logs.Error("Rename temp file to final log file error", renameErr)
			return
		}
		logs.Info("Download", fileKey, "=>", localAbsFilePath, "success", avgSpeed)
	} else {
		err = errors.New("download failed")
		logs.Info("Download", fileKey, "failed by url", fileUrl, resp.Status)
		return
	}
	return
}

func utf82GBK(text string) (string, error) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	return gbkEncoder.String(text)
}
