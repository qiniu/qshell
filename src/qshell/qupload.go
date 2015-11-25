package qshell

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
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

Valid values for zone are [aws,nb,bc]
*/

const (
	DEFAULT_PUT_THRESHOLD   int64 = 10 * 1024 * 1024 //100MB
	MIN_UPLOAD_THREAD_COUNT int64 = 1
	MAX_UPLOAD_THREAD_COUNT int64 = 100
)

type UploadConfig struct {
	//basic config
	SrcDir    string `json:"src_dir"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`

	//optional config
	PutThreshold     int64  `json:"put_threshold,omitempty"`
	KeyPrefix        string `json:"key_prefix,omitempty"`
	IgnoreDir        bool   `json:"ignore_dir,omitempty"`
	Overwrite        bool   `json:"overwrite,omitempty"`
	CheckExists      bool   `json:"check_exists,omitempty"`
	SkipFilePrefixes string `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes string `json:"skip_path_prefixes,omitempty"`
	SkipSuffixes     string `json:"skip_suffixes,omitempty"`

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

func QiniuUpload(threadCount int, uploadConfigFile string) {
	timeStart := time.Now()
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

	//make SrcDir the full path
	uploadConfig.SrcDir, _ = filepath.Abs(uploadConfig.SrcDir)

	dirCache := DirCache{}
	currentUser, err := user.Current()
	if err != nil {
		log.Error("Failed to get current user", err)
		return
	}

	pathSep := string(os.PathSeparator)
	//create job id
	jobId := Md5Hex(fmt.Sprintf("%s:%s", uploadConfig.SrcDir, uploadConfig.Bucket))

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

	log.Info("Listing local sync dir, this can take a long time, please wait paitently...")
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
			log.Error(fmt.Sprintf("Invalid cache line `%s'", line))
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
		localFlmd, _ := strconv.Atoi(items[2])
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
			log.Error(fmt.Sprintf("Error stat local file `%s' due to `%s'", localFilePath, err))
			continue
		}

		fsize := fstat.Size()
		ldbKey := fmt.Sprintf("%s => %s", localFilePath, uploadFileKey)

		log.Info(fmt.Sprintf("Uploading %s (%d/%d, %.1f%%) ...", ldbKey, currentFileCount, totalFileCount,
			float32(currentFileCount)*100/float32(totalFileCount)))

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
			flmd, _ := strconv.Atoi(string(ldbFlmd))
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
				policy.Scope = uploadConfig.Bucket + ":" + uploadFileKey
				policy.InsertOnly = 0
			}
			policy.Expires = 30 * 24 * 3600
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
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, pErr.Err))
					} else {
						log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFilePath, uploadFileKey, err))
					}
				} else {
					os.Remove(progressFpath)
					atomic.AddInt64(&successFileCount, 1)
					perr := ldb.Put([]byte(ldbKey), []byte("Y"), &ldbWOpt)
					if perr != nil {
						log.Error(fmt.Sprintf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr))
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

	log.Info("-------Upload Result-------")
	log.Info("Total:\t", currentFileCount)
	log.Info("Success:\t", successFileCount)
	log.Info("Failure:\t", failureFileCount)
	log.Info("Skipped:\t", skippedFileCount)
	log.Info("Duration:\t", time.Since(timeStart))
	log.Info("-------------------------")

}
