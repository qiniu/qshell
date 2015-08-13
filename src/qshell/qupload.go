package qshell

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	fio "github.com/qiniu/api.v6/io"
	rio "github.com/qiniu/api.v6/resumable/io"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/log"
	"github.com/qiniu/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

/*
Config file like:

{
	"up_host"		:	"http://upload.qiniu.com",
	"src_dir" 		:	"/Users/jemy/Photos",
	"access_key" 	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"bucket"		:	"test-bucket",
	"ignore_dir"	:	false,
	"key_prefix"	:	"2014/12/01/",
	"overwrite"		:	false,
	"check_exists"	:	true
}

or without up_host and key_prefix and ignore_dir and check_exists

{
	"src_dir" 		:	"/Users/jemy/Photos",
	"access_key" 	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"bucket"		:	"test-bucket",
}
*/

const (
	PUT_THRESHOLD           int64 = 10 * 1 << 20
	MIN_UPLOAD_THREAD_COUNT int64 = 1
	MAX_UPLOAD_THREAD_COUNT int64 = 100
)

type UploadConfig struct {
	SrcDir       string `json:"src_dir"`
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	Bucket       string `json:"bucket"`
	UpHost       string `json:"up_host,omitempty"`
	KeyPrefix    string `json:"key_prefix,omitempty"`
	IgnoreDir    bool   `json:"ignore_dir,omitempty"`
	Overwrite    bool   `json:"overwrite,omitempty"`
	CheckExists  bool   `json:"check_exists,omitempty"`
	SkipPrefixes string `json:"skip_prefixes,omitempty"`
	SkipSuffixes string `json:"skip_suffixes,omitempty"`
}

var upSettings = rio.Settings{
	ChunkSize: 1 * 1024 * 1024,
	TryTimes:  5,
}

func QiniuUpload(threadCount int, uploadConfigFile string) {
	fp, err := os.Open(uploadConfigFile)
	if err != nil {
		log.Error(fmt.Sprintf("Open upload config file `%s' error due to `%s'", uploadConfigFile, err))
		return
	}
	defer fp.Close()
	configData, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Error(fmt.Sprintf("Read upload config file `%s' error due to `%s'", uploadConfigFile, err))
		return
	}
	var uploadConfig UploadConfig
	err = json.Unmarshal(configData, &uploadConfig)
	if err != nil {
		log.Error(fmt.Sprintf("Parse upload config file `%s' errror due to `%s'", uploadConfigFile, err))
		return
	}
	if _, err := os.Stat(uploadConfig.SrcDir); err != nil {
		log.Error("Upload config error for parameter `SrcDir`,", err)
		return
	}
	dirCache := DirCache{}
	currentUser, err := user.Current()
	if err != nil {
		log.Error("Failed to get current user", err)
		return
	}

	pathSep := string(os.PathSeparator)
	//create job id
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(strings.TrimSuffix(uploadConfig.SrcDir, pathSep) + ":" + uploadConfig.Bucket))
	jobId := fmt.Sprintf("%x", md5Hasher.Sum(nil))

	//local storage path
	storePath := filepath.Join(currentUser.HomeDir, ".qshell", "qupload", jobId)
	err = os.MkdirAll(storePath, 0775)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to mkdir `%s' due to `%s'", storePath, err))
		return
	}

	//cache file
	cacheFileName := filepath.Join(storePath, jobId+".cache")
	//leveldb folder
	leveldbFileName := filepath.Join(storePath, jobId+".ldb")

	totalFileCount := dirCache.Cache(uploadConfig.SrcDir, cacheFileName)
	ldb, err := leveldb.OpenFile(leveldbFileName, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Open leveldb `%s' failed due to `%s'", leveldbFileName, err))
		return
	}
	defer ldb.Close()
	//sync
	ufp, err := os.Open(cacheFileName)
	if err != nil {
		log.Error(fmt.Sprintf("Open cache file `%s' failed due to `%s'", cacheFileName, err))
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

	//use host if not empty
	if uploadConfig.UpHost != "" {
		conf.UP_HOST = uploadConfig.UpHost
	}
	//set resume upload settings
	rio.SetSettings(&upSettings)
	mac := digest.Mac{uploadConfig.AccessKey, []byte(uploadConfig.SecretKey)}

	for bScanner.Scan() {
		line := strings.TrimSpace(bScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			log.Error(fmt.Sprintf("Invalid cache line `%s'", line))
			continue
		}

		localFname := items[0]
		currentFileCount += 1

		skip := false
		//check skip local file or folder
		if uploadConfig.SkipPrefixes != "" {
			//unpack skip prefix
			skipPrefixes := strings.Split(uploadConfig.SkipPrefixes, ",")
			for _, prefix := range skipPrefixes {
				if strings.HasPrefix(localFname, prefix) {
					log.Info(fmt.Sprintf("Skip by prefix '%s' for local file %s", prefix, localFname))
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
				if strings.HasSuffix(localFname, suffix) {
					log.Info(fmt.Sprintf("Skip by suffix '%s' for local file %s", suffix, localFname))
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
		localFlmd, _ := strconv.Atoi(items[2])
		uploadFileKey := localFname

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

		localFilePath := filepath.Join(uploadConfig.SrcDir, localFname)
		fstat, err := os.Stat(localFilePath)
		if err != nil {
			log.Error(fmt.Sprintf("Error stat local file `%s' due to `%s'", localFilePath, err))
			return
		}

		fsize := fstat.Size()
		ldbKey := fmt.Sprintf("%s => %s", localFilePath, uploadFileKey)

		log.Info(fmt.Sprintf("Uploading %s (%d/%d, %.1f%%) ...", ldbKey, currentFileCount, totalFileCount,
			float32(currentFileCount)*100/float32(totalFileCount)))

		rsClient := rs.New(&mac)
		log.Debug(fmt.Sprintf("Checking %s ...", ldbKey))

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
					atomic.AddInt64(&successFileCount, 1)
					log.Info(fmt.Sprintf("File %s already exists in bucket, ignore this upload", uploadFileKey))
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
			flmd, _ := strconv.Atoi(string(ldbFlmd))
			//not exist, return ErrNotFound
			//check last modified

			if err == nil && localFlmd == flmd {
				atomic.AddInt64(&successFileCount, 1)
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
				policy.Scope = uploadConfig.Bucket + ":" + uploadFileKey
				policy.InsertOnly = 0
			}
			policy.Expires = 24 * 3600
			uptoken := policy.Token(&mac)
			if fsize > PUT_THRESHOLD {
				putRet := rio.PutRet{}
				err := rio.PutFile(nil, &putRet, uptoken, uploadFileKey, localFilePath, nil)
				if err != nil {
					atomic.AddInt64(&failureFileCount, 1)
					if pErr, ok := err.(*rpc.ErrorInfo); ok {
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, pErr.Err))
					} else {
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, err))
					}
				} else {
					atomic.AddInt64(&successFileCount, 1)
					perr := ldb.Put([]byte(ldbKey), []byte("Y"), &ldbWOpt)
					if perr != nil {
						log.Error(fmt.Sprintf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr))
					}
				}
			} else {
				putRet := fio.PutRet{}
				err := fio.PutFile(nil, &putRet, uptoken, uploadFileKey, localFilePath, nil)
				if err != nil {
					atomic.AddInt64(&failureFileCount, 1)
					if pErr, ok := err.(*rpc.ErrorInfo); ok {
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, pErr.Err))
					} else {
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, err))
					}
				} else {
					atomic.AddInt64(&successFileCount, 1)
					perr := ldb.Put([]byte(ldbKey), []byte(strconv.Itoa(localFlmd)), &ldbWOpt)
					if perr != nil {
						log.Error(fmt.Sprintf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr))
					}
				}
			}
		}()

	}
	upWorkGroup.Wait()

	log.Info("-------Upload Done-------")
	log.Info("Total:\t", currentFileCount)
	log.Info("Success:\t", successFileCount)
	log.Info("Failure:\t", failureFileCount)
	log.Info("Skipped:\t", skippedFileCount)
	log.Info("-------------------------")

}
