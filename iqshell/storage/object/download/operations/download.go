package operations

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MIN_DOWNLOAD_THREAD_COUNT = 1    // 最小的下载线程数目
	MAX_DOWNLOAD_THREAD_COUNT = 2000 // 最大下载线程数目
)

var downloadTasks chan func()
var initDownOnce sync.Once

func doDownload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

type DownloadInfo struct {
	ConfigFile     string
	ThreadCount    int
	DownloadConfig config.DownloadConfig
}

// 【qdownload] 批量下载文件， 可以下载以前缀的文件，也可以下载一个文件列表
func Download(info DownloadInfo) {
	var downloadConfig storage.DownloadConfig

	cfh, oErr := os.Open(info.ConfigFile)
	if oErr != nil {
		log.ErrorF("open file: %s: %v", info.ConfigFile, oErr)
		os.Exit(1)
	}
	content, rErr := ioutil.ReadAll(cfh)
	if rErr != nil {
		log.ErrorF("read configFile content: %v", rErr)
		os.Exit(1)
	}

	// remove windows utf-8 BOM
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	uErr := json.Unmarshal(content, &downloadConfig)

	if uErr != nil {
		log.ErrorF("decode configFile content: %v", uErr)
		os.Exit(1)
	}

	destFileInfo, err := os.Stat(downloadConfig.DestDir)
	if err != nil {
		log.ErrorF("stat %s: %v", downloadConfig.DestDir, err)
		os.Exit(1)
	}

	if !destFileInfo.IsDir() {
		log.Error("Download dest dir should be a directory")
		os.Exit(data.STATUS_HALT)
	}

	if info.ThreadCount < MIN_DOWNLOAD_THREAD_COUNT || info.ThreadCount > MAX_DOWNLOAD_THREAD_COUNT {
		log.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			MIN_DOWNLOAD_THREAD_COUNT, MAX_DOWNLOAD_THREAD_COUNT)

		if info.ThreadCount < MIN_DOWNLOAD_THREAD_COUNT {
			info.ThreadCount = MIN_DOWNLOAD_THREAD_COUNT
		} else if info.ThreadCount > MAX_DOWNLOAD_THREAD_COUNT {
			info.ThreadCount = MAX_DOWNLOAD_THREAD_COUNT
		}
	}

	rootPath := info.DownloadConfig.RecordRoot
	if rootPath == "" {
		rootPath = workspace.GetWorkspace()
	}
	if rootPath == "" {
		log.Error(alert.CannotEmpty("root path", ""))
		os.Exit(1)
	}

	downConfig := &info.DownloadConfig
	timeStart := time.Now()
	//create job id
	jobId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", downConfig.DestDir, downConfig.Bucket, downConfig.KeyFile))

	//local storage path
	storePath := filepath.Join(rootPath, "qdownload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		log.ErrorF("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		os.Exit(data.STATUS_ERROR)
	}

	//init log settings
	defaultLogFile := filepath.Join(storePath, fmt.Sprintf("%s.log", jobId))
	//init log level
	logLevel := log.LevelInfo
	logRotate := 1
	if downConfig.LogRotate > 0 {
		logRotate = downConfig.LogRotate
	}
	switch downConfig.LogLevel {
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarning
	case "error":
		logLevel = log.LevelError
	default:
		logLevel = log.LevelInfo
	}

	//init log writer
	if downConfig.LogFile == "" {
		//set default log file
		downConfig.LogFile = defaultLogFile
	}

	//Todo: 处理 stdout 不输出日志
	if !downConfig.LogStdout {
		// logs.GetBeeLogger().DelLogger(logs.AdapterConsole)
	}
	downConfig.Init()
	//open log file
	fmt.Println("Writing download log to file", downConfig.LogFile)

	//daily rotate
	logCfg := log.Config{
		Filename: downConfig.LogFile,
		Level:    int(logLevel),
		Daily:    true,
		MaxDays:  logRotate,
	}
	log.LoadFileLogger(logCfg)
	// logs.SetLogger(logs.AdapterFile, logCfg.ToJson())
	fmt.Println()

	downloadDomain := downConfig.DownloadDomain()
	if downloadDomain == "" {
		domainOfBucket, dErr := bucket.DomainOfBucket(downConfig.Bucket)
		if dErr != nil {
			log.Error("get domains of bucket: ", dErr)
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
		log.Error("Open resume record leveldb error", openErr)
		os.Exit(data.STATUS_ERROR)
	}
	defer resumeLevelDb.Close()
	//sync underlying writes from the OS buffer cache
	//through to actual disk
	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}
	if downConfig.KeyFile != "" {
		fmt.Println("Batch stat file info, this may take a long time, please wait...")
		generateMiddleFile(downConfig, jobListFileName)
	} else {
		//list bucket, prepare file list to download
		log.InfoF("Listing bucket `%s` by prefix `%s`", downConfig.Bucket, downConfig.Prefix)
		bucket.ListToFile(bucket.ListToFileApiInfo{
			ListApiInfo: bucket.ListApiInfo{
				Bucket:            downConfig.Bucket,
				Prefix:            downConfig.Prefix,
				Marker:            "",
				Delimiter:         "",
				StartTime:         time.Time{},
				EndTime:           time.Time{},
				Suffixes:          nil,
				MaxRetry:          20,
				StopWhenListError: false,
			},
			FilePath:    "",
			AppendMode:  false,
			Readable:    false,
		}, func(marker string, err error) {
			log.ErrorF("marker: %s", marker)
			log.ErrorF("list bucket Error: %v", err)
		})
	}

	//init wait group
	downWaitGroup := sync.WaitGroup{}

	initDownOnce.Do(func() {
		downloadTasks = make(chan func(), info.ThreadCount)
		for i := 0; i < info.ThreadCount; i++ {
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
		log.Error("Open list file error", openErr)
		os.Exit(data.STATUS_ERROR)
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
					log.InfoF("Skip download `%s`, suffix filter not match", fileKey)
					continue
				}
			}

			fileSize, pErr := strconv.ParseInt(items[1], 10, 64)
			if pErr != nil {
				log.ErrorF("Invalid list line", line)
				continue
			}

			fileHash := items[2]

			fileMtime, pErr := strconv.ParseInt(items[3], 10, 64)
			if pErr != nil {
				log.Error("Invalid list line", line)
				continue
			}

			var fileUrl string
			if downConfig.Public {
				fileUrl = download.PublicUrl(download.PublicUrlApiInfo{
					BucketDomain: downloadDomain,
					Key:          fileKey,
					UseHttps:     downConfig.UseHttps,
				})
			} else {
				fileUrl = download.PrivateUrl(download.PrivateUrlApiInfo{
					BucketDomain: downloadDomain,
					Key:          fileKey,
					UseHttps:     downConfig.UseHttps,
				})
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
						log.InfoF("Local file `%s` exists, same as in bucket, download skip", localAbsFilePath)
						existsFileCount += 1
						continue
					} else {
						//somthing changed, must download a new file
						log.InfoF("Local file `%s` exists, but remote file changed, go to download", localAbsFilePath)
						downNewFile = true
					}
				} else {
					// 数据库中不存在
					if downConfig.CheckHash {
						// 无法验证信息，重新下载
						downNewFile = true
						log.InfoF("Local file `%s` exists, but can't find file info from db, go to download", localAbsFilePath)
						continue
					}

					// 不验证 hash 仅仅验证 size, size 相同则认为 文件不变
					if localFileInfo.Size() != fileSize {
						log.InfoF("Local file `%s` exists, size not the same as in bucket, go to download", localAbsFilePath)
						downNewFile = true
					} else {
						//treat the local file not changed, write to leveldb, though may not accurate
						//nothing to do
						log.WarningF("Local file `%s` exists with same size as `%s`, treat it not changed", localAbsFilePath, fileKey)
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
									log.ErrorF("Rename temp file `%s` to final file `%s` error", localAbsFilePathTemp, localAbsFilePath,
										renameErr)
								}
								continue
							}
						} else {
							log.InfoF("Local tmp file `%s` exists, but remote file changed, go to download", localAbsFilePathTemp)
							downNewFile = true
						}
					} else {
						//log tmp file exists, but no record in leveldb, download a new file
						log.InfoF("Local tmp file `%s` exists, but no record in leveldb ,go to download", localAbsFilePathTemp)
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

	log.InfoF("-------Download Result-------")
	log.InfoF("%10s%10d", "Total:", totalFileCount)
	log.InfoF("%10s%10d", "Skipped:", skipBySuffixes)
	log.InfoF("%10s%10d", "Exists:", existsFileCount)
	log.InfoF("%10s%10d", "Success:", successFileCount)
	log.InfoF("%10s%10d", "Update:", updateFileCount)
	log.InfoF("%10s%10d", "Failure:", failureFileCount)
	log.InfoF("%10s%15s", "Duration:", time.Since(timeStart))
	log.InfoF("-----------------------------")
	fmt.Println("\nSee download log at path", downConfig.LogFile)

	if failureFileCount > 0 {
		os.Exit(data.STATUS_ERROR)
	}
}

type EntryPath struct {
	Bucket  string
	Key     string
	PutTime string
}

func generateMiddleFile(downloadConfig *config.DownloadConfig, jobListFileName string) {
	kFile, kErr := os.Open(downloadConfig.KeyFile)
	if kErr != nil {
		log.ErrorF("open KeyFile: %s: %v\n", downloadConfig.KeyFile, kErr)
		os.Exit(data.STATUS_ERROR)
	}
	defer kFile.Close()

	jobListFh, jErr := os.Create(jobListFileName)
	if jErr != nil {
		log.ErrorF("open jobListFileName: %s: %v\n", jobListFileName, kErr)
		os.Exit(data.STATUS_ERROR)
	}
	defer jobListFh.Close()

	scanner := bufio.NewScanner(kFile)
	entries := make([]batch.Operation, 0, downloadConfig.BatchNum)

	writeEntry := func() {
		bret, _ := batch.Some(entries)
		if len(bret) == len(entries) {
			for j, item := range bret {
				entry := entries[j].(object.StatusApiInfo)
				if item.Code != 200 || item.Error != "" {
					fmt.Fprintln(os.Stderr, entry.Key+"\t"+item.Error)
				} else {
					fmt.Fprintf(jobListFh, "%s\t%d\t%s\t%d\t%s\t%d\n", entry.Key, item.FSize, item.Hash, item.PutTime, item.MimeType, item.Type)
				}
			}
		}
		entries = entries[:0]
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(entries) < downloadConfig.BatchNum {
			entries = append(entries, object.StatusApiInfo{
				Bucket: downloadConfig.Bucket,
				Key:    line,
			})
		}
		if len(entries) == downloadConfig.BatchNum {
			writeEntry()
		}
	}
	// 最后一批数量小于分割的batchNum
	if len(entries) > 0 {
		writeEntry()
	}
}

//file key -> mtime
func downloadFile(downConfig *config.DownloadConfig, fileKey, fileUrl, domainOfBucket string, fileSize int64,
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
		log.Error("MkdirAll failed for", localFileDir, mkdirErr)
		return
	}

	log.Info("Downloading", fileKey, "=>", localAbsFilePath)
	//new request
	req, reqErr := http.NewRequest("GET", fileUrl, nil)
	if reqErr != nil {
		err = reqErr
		log.Info("New request", fileKey, "failed by url", fileUrl, reqErr)
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
		log.Info("Download", fileKey, "failed by url", fileUrl, respErr)
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
			log.Error("Open local file", localAbsFilePathTemp, "failed", openErr)
			return
		}

		cpCnt, cpErr := io.Copy(localFp, resp.Body)
		if cpErr != nil {
			err = cpErr
			localFp.Close()
			log.Error("Download", fileKey, "failed", cpErr)
			return
		}
		localFp.Close()

		if downConfig.CheckHash {
			log.Info("Download", fileKey, " check hash")

			hashFile, openTempFileErr := os.Open(localFilePathTempTarget)
			if openTempFileErr != nil {
				err = openTempFileErr
				return
			}

			downloadFileHash := ""
			if fileHash != "" && utils.IsSignByEtagV2(fileHash) {
				bucketManager, gErr := bucket.GetBucketManager()
				if gErr != nil {
					err = gErr
					return
				}
				stat, errs := bucketManager.Stat(downConfig.Bucket, fileKey)
				if errs == nil {
					downloadFileHash, err = utils.EtagV2(hashFile, stat.Parts)
				}
				log.Info("Download", fileKey, " v2 local hash:", downloadFileHash, " server hash:", fileHash)
			} else {
				downloadFileHash, err = utils.EtagV1(hashFile)
				log.Info("Download", fileKey, " v1 local hash:", downloadFileHash, " server hash:", fileHash)
			}
			if err != nil {
				log.Error("get file hash error", err)
				return
			}

			if downloadFileHash != fileHash {
				err = errors.New(fileKey + ": except hash:" + fileHash + " but is:" + downloadFileHash)
				log.Error("file error", err)
				return
			}
		}

		endDown := time.Now().Unix()
		avgSpeed := fmt.Sprintf("%.2fKB/s", float64(cpCnt)/float64(endDown-startDown)/1024)

		//move temp file to log file
		renameErr := os.Rename(localFilePathTempTarget, localFilePathTarget)
		if renameErr != nil {
			err = renameErr
			log.Error("Rename temp file to final log file error", renameErr)
			return
		}
		log.Info("Download", fileKey, "=>", localAbsFilePath, "success", avgSpeed)
	} else {
		err = errors.New("download failed")
		log.Info("Download", fileKey, "failed by url", fileUrl, resp.Status)
		return
	}
	return
}

func utf82GBK(text string) (string, error) {
	var gbkEncoder = simplifiedchinese.GBK.NewEncoder()
	return gbkEncoder.String(text)
}
