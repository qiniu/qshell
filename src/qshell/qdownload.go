package qshell

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	"qiniu/log"
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
	"access_key"	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"prefix"		:	"demo/",
	"suffix"		: 	".mp4",
	"referer"		: 	""
}
*/

const (
	MIN_DOWNLOAD_THREAD_COUNT = 1
	MAX_DOWNLOAD_THREAD_COUNT = 100
)

type DownloadConfig struct {
	DestDir   string `json:"dest_dir"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Prefix    string `json:"prefix,omitempty"`
	Suffix    string `json:"suffix,omitempty"`
	Referer   string `json:"referer,omitemtpy"`
}

var downloadTasks chan func()
var initDownOnce sync.Once

func doDownload(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func QiniuDownload(threadCount int, downloadConfigFile string) {
	timeStart := time.Now()
	cnfFp, openErr := os.Open(downloadConfigFile)
	if openErr != nil {
		log.Error("Open download config file error,", openErr)
		return
	}
	defer cnfFp.Close()
	cnfData, rErr := ioutil.ReadAll(cnfFp)
	if rErr != nil {
		log.Error("Read download config file error,", rErr)
		return
	}
	downConfig := DownloadConfig{}
	cnfErr := json.Unmarshal(cnfData, &downConfig)
	if cnfErr != nil {
		log.Error("Parse download config error,", cnfErr)
		return
	}

	//check dest dir
	destFileInfo, statErr := os.Stat(downConfig.DestDir)
	if statErr != nil {
		log.Error("Invalid dest dir,", statErr)
		return
	}

	if !destFileInfo.IsDir() {
		log.Error("Dest dir should be a directory")
		return
	}

	mac := digest.Mac{downConfig.AccessKey, []byte(downConfig.SecretKey)}
	//get bucket zone info
	bucketInfo, gErr := GetBucketInfo(&mac, downConfig.Bucket)
	if gErr != nil {
		log.Error("Get bucket region info error,", gErr)
		return
	}
	//get domains of bucket
	domainsOfBucket, gErr := GetDomainsOfBucket(&mac, downConfig.Bucket)
	if gErr != nil {
		log.Error("Get domains of bucket error,", gErr)
		return
	}

	if len(domainsOfBucket) == 0 {
		log.Error("No domains found for bucket", downConfig.Bucket)
		return
	}

	domainOfBucket := domainsOfBucket[0]

	//set up host
	SetZone(bucketInfo.Region)
	ioProxyAddress := conf.IO_HOST

	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", downConfig.DestDir, downConfig.Bucket))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "qdownload", jobId)
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		log.Errorf("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		return
	}

	jobListFileName := filepath.Join(storePath, fmt.Sprintf("%s.list", jobId))
	resumeFile := filepath.Join(storePath, fmt.Sprintf("%s.ldb", jobId))
	resumeLevelDb, openErr := leveldb.OpenFile(resumeFile, nil)
	if openErr != nil {
		log.Error("Open resume record leveldb error", openErr)
		return
	}
	defer resumeLevelDb.Close()
	//sync underlying writes from the OS buffer cache
	//through to actual disk
	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	//list bucket, prepare file list to download
	log.Infof("Listing bucket `%s` by prefix `%s`", downConfig.Bucket, downConfig.Prefix)
	listErr := ListBucket(&mac, downConfig.Bucket, downConfig.Prefix, "", jobListFileName)
	if listErr != nil {
		log.Error("List bucket error", listErr)
		return
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

	totalFileCount = GetFileLineCount(jobListFileName)

	//open prepared file list to download files
	listFp, openErr := os.Open(jobListFileName)
	if openErr != nil {
		log.Error("Open list file error", openErr)
		return
	}
	defer listFp.Close()

	listScanner := bufio.NewScanner(listFp)
	listScanner.Split(bufio.ScanLines)
	//key, fsize, etag, lmd, mime, enduser
	for listScanner.Scan() {
		currentFileCount += 1
		line := strings.TrimSpace(listScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) >= 4 {
			fileKey := items[0]

			if downConfig.Suffix != "" && !strings.HasSuffix(fileKey, downConfig.Suffix) {
				//skip files by suffix specified
				continue
			}

			fileSize, pErr := strconv.ParseInt(items[1], 10, 64)
			if pErr != nil {
				log.Errorf("Invalid list line", line)
				continue
			}

			fileMtime, pErr := strconv.ParseInt(items[3], 10, 64)
			if pErr != nil {
				log.Errorf("Invalid list line", line)
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
			localFilePathTmp := fmt.Sprintf("%s.tmp", localFilePath)
			localFileInfo, statErr := os.Stat(localFilePath)

			var downNewLog bool
			var fromBytes int64

			if statErr == nil {
				//log file exists, check whether have updates
				oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePath), nil)
				if notFoundErr == nil {
					//if exists
					oldFileInfoItems := strings.Split(string(oldFileInfo), "|")
					oldFileLmd, _ := strconv.ParseInt(oldFileInfoItems[0], 10, 64)
					//oldFileSize, _ := strconv.ParseInt(oldFileInfoItems[1], 10, 64)
					if oldFileLmd == fileMtime && localFileInfo.Size() == fileSize {
						//nothing change, ignore
						existsFileCount += 1
						continue
					} else {
						//somthing changed, must download a new file
						downNewLog = true
					}
				} else {
					if localFileInfo.Size() != fileSize {
						downNewLog = true
					} else {
						//treat the local file not changed, write to leveldb, though may not accurate
						//nothing to do
						atomic.AddInt64(&existsFileCount, 1)
						continue
					}
				}
			} else {
				//check whether tmp file exists
				localTmpFileInfo, statErr := os.Stat(localFilePathTmp)
				if statErr == nil {
					//if tmp file exists, check whether last modify changed
					oldFileInfo, notFoundErr := resumeLevelDb.Get([]byte(localFilePath), nil)
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
								renameErr := os.Rename(localFilePathTmp, localFilePath)
								if renameErr != nil {
									log.Error("Rename temp file to final log file error", renameErr)
								}
								continue
							}
						} else {
							downNewLog = true
						}
					} else {
						downNewLog = true
					}
				} else {
					//no log file exists, donwload a new log file
					downNewLog = true
				}
			}

			//set file info in leveldb
			rKey := localAbsFilePath
			rVal := fmt.Sprintf("%d|%d", fileMtime, fileSize)
			resumeLevelDb.Put([]byte(rKey), []byte(rVal), &ldbWOpt)

			//download new
			downWaitGroup.Add(1)
			downloadTasks <- func() {
				defer downWaitGroup.Done()

				downErr := downloadFile(downConfig.DestDir, fileKey, fileUrl, domainOfBucket, fileSize, fromBytes)
				if downErr != nil {
					atomic.AddInt64(&failureFileCount, 1)
				} else {
					atomic.AddInt64(&successFileCount, 1)
					if !downNewLog {
						atomic.AddInt64(&updateFileCount, 1)
					}
				}
			}
		}
	}

	//wait for all tasks done
	downWaitGroup.Wait()

	log.Info("-------Download Result-------")
	log.Infof("%10s%10d\n", "Total:", totalFileCount)
	log.Infof("%10s%10d\n", "Exists:", existsFileCount)
	log.Infof("%10s%10d\n", "Success:", successFileCount)
	log.Infof("%10s%10d\n", "Update:", updateFileCount)
	log.Infof("%10s%10d\n", "Failure:", failureFileCount)
	log.Infof("%10s%15s\n", "Duration:", time.Since(timeStart))
	log.Info("-----------------------------")
}

/*
@param ioProxyHost - like http://iovip.qbox.me
*/
func makePrivateDownloadLink(mac *digest.Mac, domainOfBucket, ioProxyAddress, fileKey string) (fileUrl string) {
	publicUrl := fmt.Sprintf("http://%s/%s", domainOfBucket, fileKey)
	deadline := time.Now().Add(time.Hour * 24 * 30).Unix()
	privateUrl := PrivateUrl(mac, publicUrl, deadline)

	//replace the io proxy host
	fileUrl = strings.Replace(privateUrl, fmt.Sprintf("http://%s", domainOfBucket), ioProxyAddress, -1)
	return
}

//file key -> mtime
func downloadFile(destDir, fileName, fileUrl, domainsOfBucket string, fileSize int64, fromBytes int64) (err error) {
	startDown := time.Now().Unix()
	localFilePath := filepath.Join(destDir, fileName)
	localFileDir := filepath.Dir(localFilePath)
	localFilePathTmp := fmt.Sprintf("%s.tmp", localFilePath)

	mkdirErr := os.MkdirAll(localFileDir, 0775)
	if mkdirErr != nil {
		err = mkdirErr
		log.Error("MkdirAll failed for", localFileDir, mkdirErr)
		return
	}

	log.Info("Downloading", fileName, "=>", localFilePath)
	//new request
	req, reqErr := http.NewRequest("GET", fileUrl, nil)
	if reqErr != nil {
		err = reqErr
		log.Info("New request", fileName, "failed by url", fileUrl, reqErr)
		return
	}
	//set host
	req.Host = domainsOfBucket

	if fromBytes != 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", fromBytes))
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = respErr
		log.Info("Download", fileName, "failed by url", fileUrl, respErr)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		var localFp *os.File
		var openErr error
		if fromBytes != 0 {
			localFp, openErr = os.OpenFile(localFilePathTmp, os.O_APPEND|os.O_WRONLY, 0655)
		} else {
			localFp, openErr = os.Create(localFilePathTmp)
		}

		if openErr != nil {
			err = openErr
			log.Error("Open local file", localFilePathTmp, "failed", openErr)
			return
		}

		cpCnt, cpErr := io.Copy(localFp, resp.Body)
		if cpErr != nil {
			err = cpErr
			localFp.Close()
			log.Error("Download", fileName, "failed", cpErr)
			return
		}
		localFp.Close()

		endDown := time.Now().Unix()
		avgSpeed := fmt.Sprintf("%.2fKB/s", float64(cpCnt)/float64(endDown-startDown)/1024)

		//move temp file to log file
		renameErr := os.Rename(localFilePathTmp, localFilePath)
		if renameErr != nil {
			err = renameErr
			log.Error("Rename temp file to final log file error", renameErr)
			return
		}
		log.Info("Download", fileName, "=>", localFilePath, "success", avgSpeed)
	} else {
		err = errors.New("download failed")
		log.Info("Download", fileName, "failed by url", fileUrl, resp.Status)
		return
	}
	return
}
