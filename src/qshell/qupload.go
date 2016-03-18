package qshell

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	"qiniu/log"
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
	"src_dir"				:	"/Users/jemy/Photos",
	"file_list"             :   "",
	"access_key"			:	"<Your AccessKey>",
	"secret_key"			:	"<Your SecretKey>",
	"bucket"				:	"test-bucket",
	"put_threshold"			:	10000000,
	"key_prefix"			:	"2014/12/01/",
	"ignore_dir"			:	false,
	"overwrite"				:	false,
	"check_exists"			:	true,
	"skip_file_prefixes"	:	"IMG_",
	"skip_path_prefixes"	:	"tmp/,bin/,obj/",
	"skip_suffixes"			:	".exe,.obj,.class",
	"up_host"				:	"http://upload.qiniu.com",
	"zone"					:	"bc",
	"bind_up_ip"			:	"",
	"bind_rs_ip"			:	"",
	"bind_nic_ip"			:	"",
	"rescan_local"			:	false

}

or the simplest one

{
	"src_dir" 		:	"/Users/jemy/Photos",
	"access_key" 	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"bucket"		:	"test-bucket",
}

Valid values for zone are [aws,nb,bc]
*/

const (
	DEFAULT_PUT_THRESHOLD   int64 = 10 * 1024 * 1024 //10MB
	MIN_UPLOAD_THREAD_COUNT int64 = 1
	MAX_UPLOAD_THREAD_COUNT int64 = 100
)

type UploadInfo struct {
	TotalFileCount int64 `json:"total_file_count"`
}

type UploadConfig struct {
	//basic config
	SrcDir    string `json:"src_dir"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`

	//optional config
	FileList         string `json:"file_list,omitempty"`
	PutThreshold     int64  `json:"put_threshold,omitempty"`
	KeyPrefix        string `json:"key_prefix,omitempty"`
	IgnoreDir        bool   `json:"ignore_dir,omitempty"`
	Overwrite        bool   `json:"overwrite,omitempty"`
	CheckExists      bool   `json:"check_exists,omitempty"`
	SkipFilePrefixes string `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes string `json:"skip_path_prefixes,omitempty"`
	SkipSuffixes     string `json:"skip_suffixes,omitempty"`
	RescanLocal      bool   `json:"rescan_local,omitempty"`

	//advanced config
	Zone   string `json:"zone,omitempty"`
	UpHost string `json:"up_host,omitempty"`

	BindUpIp string `json:"bind_up_ip,omitempty"`
	BindRsIp string `json:"bind_rs_ip,omitempty"`

	//local network interface card config
	BindNicIp string `json:"bind_nic_ip,omitempty"`
}

var upSettings = rio.Settings{
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  7,
}

func QiniuUpload(threadCount int, uploadConfig *UploadConfig) {
	timeStart := time.Now()

	//make SrcDir the full path
	uploadConfig.SrcDir, _ = filepath.Abs(uploadConfig.SrcDir)

	dirCache := DirCache{}

	pathSep := string(os.PathSeparator)
	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket))

	//local storage path
	storePath := filepath.Join(".qshell", "qupload", jobId)
	if err := os.MkdirAll(storePath, 0775); err != nil {
		log.Errorf("Failed to mkdir `%s' due to `%s'", storePath, err)
		return
	}

	//find the local file list, by specified or by config
	var cacheResultName string
	var totalFileCount int64
	_, fStatErr := os.Stat(uploadConfig.FileList)
	if uploadConfig.FileList != "" && fStatErr == nil {
		//use specified file list
		cacheResultName = uploadConfig.FileList
		totalFileCount = getFileLineCount(cacheResultName)
	} else {
		//cache file
		cacheResultName = filepath.Join(storePath, jobId+".cache")
		cacheTempName := filepath.Join(storePath, jobId+".cache.temp")
		cacheCountName := filepath.Join(storePath, jobId+".count")

		rescanLocalDir := false
		if _, statErr := os.Stat(cacheResultName); statErr == nil {
			//file already exists
			rescanLocalDir = uploadConfig.RescanLocal
		} else {
			rescanLocalDir = true
		}

		if rescanLocalDir {
			fmt.Println("Listing local sync dir, this can take a long time, please wait patiently...")
			totalFileCount = dirCache.Cache(uploadConfig.SrcDir, cacheTempName)
			if rErr := os.Remove(cacheResultName); rErr != nil {
				log.Debug("Remove the old cached file error", rErr)
			}
			if rErr := os.Rename(cacheTempName, cacheResultName); rErr != nil {
				fmt.Println("Rename the temp cached file error", rErr)
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
							log.Errorf("Write local cached count file error %s", cErr)
						} else {
							cFp.Close()
						}
					}
				}()
			} else {
				log.Errorf("Open local cached count file error %s", cErr)
			}
		} else {
			fmt.Println("Use the last cached local sync dir file list ...")
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
				log.Warnf("Open local cached count file error %s", rErr)
			}
		}
	}

	//leveldb folder
	leveldbFileName := filepath.Join(storePath, jobId+".ldb")
	ldb, err := leveldb.OpenFile(leveldbFileName, nil)
	if err != nil {
		log.Errorf("Open leveldb `%s' failed due to `%s'", leveldbFileName, err)
		return
	}
	defer ldb.Close()
	//sync
	ufp, err := os.Open(cacheResultName)
	if err != nil {
		log.Errorf("Open list file `%s' failed due to `%s'", cacheResultName, err)
		return
	}
	defer ufp.Close()
	bScanner := bufio.NewScanner(ufp)
	bScanner.Split(bufio.ScanLines)

	var currentFileCount int64 = 0
	var successFileCount int64 = 0
	var failureFileCount int64 = 0
	var skippedFileCount int64 = 0

	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	upWorkGroup := sync.WaitGroup{}
	upCounter := 0
	threadThreshold := threadCount + 1

	//chunk upload threshold
	putThreshold := DEFAULT_PUT_THRESHOLD
	if uploadConfig.PutThreshold > 0 {
		putThreshold = uploadConfig.PutThreshold
	}

	//check zone, default nb
	switch uploadConfig.Zone {
	case ZoneAWS:
		SetZone(ZoneAWSConfig)
	case ZoneBC:
		SetZone(ZoneBCConfig)
	default:
		SetZone(ZoneNBConfig)
	}

	//use host if not empty, overwrite the default config
	if uploadConfig.UpHost != "" {
		conf.UP_HOST = uploadConfig.UpHost
	}
	//set resume upload settings
	rio.SetSettings(&upSettings)
	mac := digest.Mac{uploadConfig.AccessKey, []byte(uploadConfig.SecretKey)}

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
		line := strings.TrimSpace(bScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			log.Errorf("Invalid cache line `%s'", line)
			continue
		}

		localFpath := items[0]
		currentFileCount += 1

		skip := false
		//check skip local file or folder
		if uploadConfig.SkipPathPrefixes != "" {
			//unpack skip prefix
			skipPathPrefixes := strings.Split(uploadConfig.SkipPathPrefixes, ",")
			for _, prefix := range skipPathPrefixes {
				if strings.HasPrefix(localFpath, strings.TrimSpace(prefix)) {
					log.Debug(fmt.Sprintf("Skip by path prefix '%s' for local file %s",
						strings.TrimSpace(prefix), localFpath))
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
				localFname := filepath.Base(localFpath)
				if strings.HasPrefix(localFname, strings.TrimSpace(prefix)) {
					log.Debug(fmt.Sprintf("Skip by file prefix '%s' for local file %s",
						strings.TrimSpace(prefix), localFpath))
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
				if strings.HasSuffix(localFpath, strings.TrimSpace(suffix)) {
					log.Debug(fmt.Sprintf("Skip by suffix '%s' for local file %s",
						strings.TrimSpace(suffix), localFpath))
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
		uploadFileKey := localFpath

		if uploadConfig.IgnoreDir {
			if i := strings.LastIndex(uploadFileKey, pathSep); i != -1 {
				uploadFileKey = uploadFileKey[i+1:]
			}
		}
		if uploadConfig.KeyPrefix != "" {
			uploadFileKey = strings.Join([]string{uploadConfig.KeyPrefix, uploadFileKey}, "")
		}
		//convert \ to / under windows
		if runtime.GOOS == "windows" {
			uploadFileKey = strings.Replace(uploadFileKey, "\\", "/", -1)
		}

		localFilePath := filepath.Join(uploadConfig.SrcDir, localFpath)
		fstat, err := os.Stat(localFilePath)
		if err != nil {
			log.Errorf("Error stat local file `%s' due to `%s'", localFilePath, err)
			continue
		}

		fsize := fstat.Size()
		ldbKey := fmt.Sprintf("%s => %s", localFilePath, uploadFileKey)

		if totalFileCount != 0 {
			fmt.Println(fmt.Sprintf("Uploading %s [%d/%d, %.1f%%] ...", ldbKey, currentFileCount, totalFileCount,
				float32(currentFileCount)*100/float32(totalFileCount)))
		} else {
			fmt.Println(fmt.Sprintf("Uploading %s ...", ldbKey))
		}

		//check exists
		if uploadConfig.CheckExists {
			rsEntry, checkErr := rsClient.Stat(nil, uploadConfig.Bucket, uploadFileKey)
			if checkErr == nil {
				//compare hash
				localEtag, cErr := GetEtag(localFilePath)
				if cErr != nil {
					atomic.AddInt64(&failureFileCount, 1)
					log.Error("Calc local file hash failed,", cErr)
					continue
				}
				if rsEntry.Hash == localEtag {
					atomic.AddInt64(&skippedFileCount, 1)
					log.Debug(fmt.Sprintf("File %s already exists in bucket, ignore this upload", uploadFileKey))
					continue
				}
			} else {
				if _, ok := checkErr.(*rpc.ErrorInfo); !ok {
					//not logic error, should be network error
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

			if err == nil && localFlmd == flmd {
				log.Debug("Skip by local log for file", localFpath)
				atomic.AddInt64(&skippedFileCount, 1)
				continue
			}
		}

		//worker
		upCounter += 1
		if upCounter%threadThreshold == 0 {
			upWorkGroup.Wait()
		}
		upWorkGroup.Add(1)

		//start to upload
		go func() {
			defer upWorkGroup.Done()

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
				progressFkey := Md5Hex(fmt.Sprintf("%s:%s|%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket, localFpath, uploadFileKey))
				progressFname := fmt.Sprintf("%s.progress", progressFkey)
				progressFpath := filepath.Join(storePath, progressFname)
				putExtra.ProgressFile = progressFpath

				err := rio.PutFile(putClient, nil, &putRet, uploadFileKey, localFilePath, &putExtra)
				if err != nil {
					atomic.AddInt64(&failureFileCount, 1)
					if pErr, ok := err.(*rpc.ErrorInfo); ok {
						log.Errorf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, pErr.Err)
					} else {
						log.Errorf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, err)
					}
				} else {
					os.Remove(progressFpath)
					atomic.AddInt64(&successFileCount, 1)
					perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
					if perr != nil {
						log.Errorf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr)
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
						log.Errorf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, pErr.Err)
					} else {
						log.Errorf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, err)
					}
				} else {
					atomic.AddInt64(&successFileCount, 1)
					perr := ldb.Put([]byte(ldbKey), []byte(fmt.Sprintf("%d", localFlmd)), &ldbWOpt)
					if perr != nil {
						log.Errorf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr)
					}
				}
			}
		}()

	}
	upWorkGroup.Wait()

	fmt.Println()
	fmt.Println("-------Upload Result-------")
	fmt.Println("Total:   \t", currentFileCount)
	fmt.Println("Success: \t", successFileCount)
	fmt.Println("Failure: \t", failureFileCount)
	fmt.Println("Skipped: \t", skippedFileCount)
	fmt.Println("Duration:\t", time.Since(timeStart))
	fmt.Println("-------------------------")

}

func getFileLineCount(filePath string) (totalCount int64) {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		return
	}
	defer fp.Close()

	bScanner := bufio.NewScanner(fp)
	for bScanner.Scan() {
		totalCount += 1
	}
	return
}
