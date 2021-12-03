package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/output"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/tools"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"golang.org/x/text/encoding/simplifiedchinese"
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
	MIN_DOWNLOAD_THREAD_COUNT = 1    // 最小的下载线程数目
	MAX_DOWNLOAD_THREAD_COUNT = 2000 // 最大下载线程数目
)

// qdownload子命令用到的配置参数
type DownloadConfig struct {
	FileEncoding string `json:"file_encoding"`
	KeyFile      string `json:"key_file"`
	DestDir      string `json:"dest_dir"`
	Bucket       string `json:"bucket"`
	Prefix       string `json:"prefix,omitempty"`
	Suffixes     string `json:"suffixes,omitempty"`
	IoHost       string `json:"io_host,omitempty"`
	Public       bool   `json:"public,omitempty"`
	CheckHash    bool   `json:"check_hash"`
	//down from cdn
	Referer   string `json:"referer,omitempty"`
	CdnDomain string `json:"cdn_domain,omitempty"`
	UseHttps  bool   `json:"use_https,omitempty"`
	//log settings
	RecordRoot string `json:"record_root,omitempty"`
	LogLevel   string `json:"log_level,omitempty"`
	LogFile    string `json:"log_file,omitempty"`
	LogRotate  int    `json:"log_rotate,omitempty"`
	LogStdout  bool   `json:"log_stdout,omitempty"`

	batchNum int
}

func (d *DownloadConfig) init() {
	if d.batchNum <= 0 {
		d.batchNum = 1000
	}
}

// 获取一个存储空间的绑定的所有域名
func (d *DownloadConfig) DomainOfBucket(bm *BucketManager) (domain string, err error) {
	//get domains of bucket
	domainsOfBucket, gErr := bm.DomainsOfBucket(d.Bucket)
	if gErr != nil {
		err = fmt.Errorf("Get domains of bucket error: %v", gErr)
		return
	}

	if len(domainsOfBucket) == 0 {
		err = fmt.Errorf("No domains found for bucket: %s", d.Bucket)
		return
	}

	for _, d := range domainsOfBucket {
		if !strings.HasPrefix(d, ".") {
			domain = d
			break
		}
	}
	return
}

// 获取一个存储空间的下载域名， 默认使用用户配置的域名，如果没有就使用接口随机选择一个下载域名
func (d *DownloadConfig) DownloadDomain() (domain string) {
	if d.CdnDomain != "" {
		domain = d.CdnDomain
	} else if d.IoHost != "" {
		domain = d.IoHost
	}
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	return
}

// qdownload需要使用文件的大小等信息判断文件是否已经下载
// 所以批量下载需要使用listBucket接口产生的中间生成文件
// 对于直接使用文件列表批量下载的方式，为了产生这个中间文件，需要使用
// batchstat接口来获取文件信息
func (d *DownloadConfig) generateMiddileFile(bm *BucketManager, jobListFileName string) {
	kFile, kErr := os.Open(d.KeyFile)
	if kErr != nil {
		logs.Error("open KeyFile: %s: %v\n", d.KeyFile, kErr)
		os.Exit(config.STATUS_ERROR)
	}
	defer kFile.Close()

	jobListFh, jErr := os.Create(jobListFileName)
	if jErr != nil {
		logs.Error("open jobListFileName: %s: %v\n", jobListFileName, kErr)
		os.Exit(config.STATUS_ERROR)
	}
	defer jobListFh.Close()

	scanner := bufio.NewScanner(kFile)
	entries := make([]EntryPath, 0, d.batchNum)

	writeEntry := func() {
		bret, _ := bm.BatchStat(entries)
		if len(bret) == len(entries) {
			for j, entry := range entries {
				item := bret[j]
				if item.Code != 200 || item.Data.Error != "" {
					fmt.Fprintln(os.Stderr, entry.Key+"\t"+item.Data.Error)
				} else {
					fmt.Fprintf(jobListFh, "%s\t%d\t%s\t%d\t%s\t%d\n", entry.Key, item.Data.Fsize, item.Data.Hash, item.Data.PutTime, item.Data.MimeType, item.Data.Type)
				}
			}
		}
		entries = entries[:0]
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		entry := EntryPath{
			Bucket: d.Bucket,
			Key:    line,
		}
		if len(entries) < d.batchNum {
			entries = append(entries, entry)
		}
		if len(entries) == d.batchNum {
			writeEntry()
		}
	}
	// 最后一批数量小于分割的batchNum
	if len(entries) > 0 {
		writeEntry()
	}
}

var downloadTasks chan func()
var initDownOnce sync.Once

func doDownload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

// 【qdownload] 批量下载文件， 可以下载以前缀的文件，也可以下载一个文件列表
func QiniuDownload(threadCount int, downConfig *DownloadConfig) {
	QShellRootPath := downConfig.RecordRoot
	if QShellRootPath == "" {
		QShellRootPath = config.RootPath()
	}
	if QShellRootPath == "" {
		fmt.Fprintf(os.Stderr, "empty root path\n")
		os.Exit(1)
	}
	timeStart := time.Now()
	//create job id
	jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", downConfig.DestDir, downConfig.Bucket, downConfig.KeyFile))

	//local storage path
	storePath := filepath.Join(QShellRootPath, "qdownload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		os.Exit(config.STATUS_ERROR)
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
	downConfig.init()
	//open log file
	fmt.Println("Writing download log to file", downConfig.LogFile)

	//daily rotate
	logCfg := output.BeeLogConfig{
		Filename: downConfig.LogFile,
		Level:    logLevel,
		Daily:    true,
		MaxDays:  logRotate,
	}
	logs.SetLogger(logs.AdapterFile, logCfg.ToJson())
	fmt.Println()

	bm := GetBucketManager()

	downloadDomain := downConfig.DownloadDomain()
	if downloadDomain == "" {
		domainOfBucket, dErr := downConfig.DomainOfBucket(bm)
		if dErr != nil {
			logs.Error("get domains of bucket: ", dErr)
			os.Exit(1)
		}
		downloadDomain = domainOfBucket
	}
	if downloadDomain == "" {
		panic("download domain cannot be empty")
	}

	jobListFileName := filepath.Join(storePath, fmt.Sprintf("%s.list", jobId))
	resumeFile := filepath.Join(storePath, fmt.Sprintf("%s.ldb", jobId))
	resumeLevelDb, openErr := leveldb.OpenFile(resumeFile, nil)
	if openErr != nil {
		logs.Error("Open resume record leveldb error", openErr)
		os.Exit(config.STATUS_ERROR)
	}
	defer resumeLevelDb.Close()
	//sync underlying writes from the OS buffer cache
	//through to actual disk
	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}
	if downConfig.KeyFile != "" {
		fmt.Println("Batch stat file info, this may take a long time, please wait...")
		downConfig.generateMiddileFile(bm, jobListFileName)
	} else {
		//list bucket, prepare file list to download
		logs.Info("Listing bucket `%s` by prefix `%s`", downConfig.Bucket, downConfig.Prefix)
		listErr := bm.ListFiles(downConfig.Bucket, downConfig.Prefix, "", jobListFileName)
		if listErr != nil {
			logs.Error("List bucket error", listErr)
			os.Exit(config.STATUS_ERROR)
		}
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

	totalFileCount = utils.GetFileLineCount(jobListFileName)

	//open prepared file list to download files
	listFp, openErr := os.Open(jobListFileName)
	if openErr != nil {
		logs.Error("Open list file error", openErr)
		os.Exit(config.STATUS_ERROR)
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

			fileHash := items[2]

			fileMtime, pErr := strconv.ParseInt(items[3], 10, 64)
			if pErr != nil {
				logs.Error("Invalid list line", line)
				continue
			}

			var fileUrl string
			if downConfig.Public {
				fileUrl = bm.MakePublicDownloadLink(downloadDomain, fileKey, downConfig.UseHttps)
			} else {
				fileUrl = bm.MakePrivateDownloadLink(downloadDomain, fileKey, downConfig.UseHttps)
			}

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
				//local file exists, check whether have updates
				oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePath), nil)
				if notFoundErr == nil {
					//if exists
					oldFileInfoItems := strings.Split(string(oldFileInfo), "|")
					oldFileLmd, _ := strconv.ParseInt(oldFileInfoItems[0], 10, 64)
					//oldFileSize, _ := strconv.ParseInt(oldFileInfoItems[1], 10, 64)
					oldFileHash := ""
					if len(oldFileInfoItems) > 2 {
						oldFileHash = oldFileInfoItems[2]
					}

					if oldFileLmd == fileMtime && localFileInfo.Size() == fileSize &&
						(downConfig.CheckHash && (len(oldFileHash) == 0 || oldFileHash == fileHash)) {
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
					// 数据库中不存在
					if downConfig.CheckHash {
						// 无法验证信息，重新下载
						downNewFile = true
						logs.Info("Local file `%s` exists, but can't find file info from db, go to download", localAbsFilePath)
						continue
					}

					// 不验证 hash 仅仅验证 size, size 相同则认为 文件不变
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
					//no local file exists, download a new local file
					downNewFile = true
				}
			}

			//set file info in leveldb
			rKey := localFilePathToCheck
			rVal := fmt.Sprintf("%d|%d|%s", fileMtime, fileSize, fileHash)
			resumeLevelDb.Put([]byte(rKey), []byte(rVal), &ldbWOpt)

			//download new
			downWaitGroup.Add(1)
			downloadTasks <- func() {
				defer downWaitGroup.Done()

				downErr := downloadFile(downConfig, fileKey, fileUrl, downloadDomain, fileSize, fromBytes, fileHash)
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
		os.Exit(config.STATUS_ERROR)
	}
}

//file key -> mtime
func downloadFile(downConfig *DownloadConfig, fileKey, fileUrl, domainOfBucket string, fileSize int64,
	fromBytes int64, fileHash string) (err error) {
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

		if downConfig.CheckHash {
			logs.Info("Download", fileKey, " check hash")

			hashFile, openTempFileErr := os.Open(localFilePathTempTarget)
			if openTempFileErr != nil {
				err = openTempFileErr
				return
			}

			downloadFileHash := ""
			if fileHash != "" && tools.IsSignByEtagV2(fileHash) {
				bucketManager := GetBucketManager()
				stat, errs := bucketManager.Stat(downConfig.Bucket, fileKey)
				if errs == nil {
					downloadFileHash, err = tools.EtagV2(hashFile, stat.Parts)
				}
				logs.Info("Download", fileKey, " v2 local hash:", downloadFileHash, " server hash:", fileHash)
			} else {
				downloadFileHash, err = tools.EtagV1(hashFile)
				logs.Info("Download", fileKey, " v1 local hash:", downloadFileHash, " server hash:", fileHash)
			}
			if err != nil {
				logs.Error("get file hash error", err)
				return
			}

			if downloadFileHash != fileHash {
				err = errors.New(fileKey + ": except hash:" + fileHash + " but is:" + downloadFileHash)
				logs.Error("file error", err)
				return
			}
		}

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
