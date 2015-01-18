package qshell

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	fio "github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

/*
Config file like:

{
	"src_dir" 		:	"/Users/jemy/Photos",
	"access_key" 	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"bucket"		:	"test-bucket",
	"ignore_dir"	:	false,
	"key_prefix"	:	"2014/12/01/"
}

or without key_prefix and ignore_dir

{
	"src_dir" 		:	"/Users/jemy/Photos",
	"access_key" 	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"bucket"		:	"test-bucket",
}
*/

const (
	PUT_THRESHOLD int64 = 2 << 30
)

type UploadConfig struct {
	SrcDir    string `json:"src_dir"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	KeyPrefix string `json:"key_prefix,omitempty"`
	IgnoreDir bool   `json:"ignore_dir,omitempty"`
}

func QiniuUpload(putThreshold int64, uploadConfigFile string) {
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
	config, _ := json.Marshal(&uploadConfig)
	md5Sum := md5.Sum(config)
	storePath := fmt.Sprintf("%s/.qshell/qupload/%x", currentUser.HomeDir, md5Sum)
	err = os.MkdirAll(storePath, 0775)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to mkdir `%s' due to `%s'", storePath, err))
		return
	}
	cacheFileName := fmt.Sprintf("%s/%x.cache", storePath, md5Sum)
	leveldbFileName := fmt.Sprintf("%s/%x.ldb", storePath, md5Sum)
	listFileName := fmt.Sprintf("%s/%x.list", storePath, md5Sum)
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
	currentFileCount := 0
	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	for bScanner.Scan() {
		line := strings.TrimSpace(bScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) > 1 {
			localFname := items[0]
			uploadFname := localFname
			if uploadConfig.IgnoreDir {
				if i := strings.LastIndex(uploadFname, string(os.PathSeparator)); i != -1 {
					uploadFname = uploadFname[i+1:]
				}
			}
			if uploadConfig.KeyPrefix != "" {
				uploadFname = strings.Join([]string{uploadConfig.KeyPrefix, uploadFname}, "")
			}
			localFnameFull := strings.Join([]string{uploadConfig.SrcDir, localFname}, string(os.PathSeparator))
			//check leveldb
			currentFileCount += 1
			ldbKey := fmt.Sprintf("%s => %s", localFnameFull, uploadFname)
			log.Debug(fmt.Sprintf("Checking %s ...", ldbKey))
			_, err := ldb.Get([]byte(ldbKey), nil)
			//not exist, return ErrNotFound
			if err == nil {
				continue
			}
			fmt.Print("\033[2K\r")
			fmt.Printf("Uploading %s (%d/%d, %.0f%%) ...", ldbKey, currentFileCount, totalFileCount,
				float32(currentFileCount)*100/float32(totalFileCount))
			os.Stdout.Sync()
			fstat, err := os.Stat(localFnameFull)
			if err != nil {
				log.Error(fmt.Sprintf("Error stat local file `%s' due to `%s'", localFnameFull, err))
				continue
			}
			fsize := fstat.Size()
			mac := digest.Mac{uploadConfig.AccessKey, []byte(uploadConfig.SecretKey)}
			policy := rs.PutPolicy{}
			policy.Scope = uploadConfig.Bucket
			policy.Expires = 24 * 3600
			uptoken := policy.Token(&mac)
			if fsize > putThreshold {
				putRet := rio.PutRet{}
				err := rio.PutFile(nil, &putRet, uptoken, uploadFname, localFnameFull, nil)
				if err != nil {
					log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFnameFull, uploadFname, err))
				} else {
					perr := ldb.Put([]byte(ldbKey), []byte("Y"), &ldbWOpt)
					if perr != nil {
						log.Error(fmt.Sprintf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr))
					}
				}
			} else {
				putRet := fio.PutRet{}
				err := fio.PutFile(nil, &putRet, uptoken, uploadFname, localFnameFull, nil)
				if err != nil {
					log.Error(fmt.Sprintf("Put file `%s' => `%s' failed due to `%s'", localFnameFull, uploadFname, err))
				} else {
					perr := ldb.Put([]byte(ldbKey), []byte("Y"), &ldbWOpt)
					if perr != nil {
						log.Error(fmt.Sprintf("Put key `%s' into leveldb error due to `%s'", ldbKey, perr))
					}
				}
			}
		} else {
			log.Error(fmt.Sprintf("Error cache line `%s'", line))
		}
	}
	fmt.Println()
	fmt.Println("Upload done!")
	//list bucket
	acct := Account{
		uploadConfig.AccessKey,
		uploadConfig.SecretKey,
	}
	bucketLister := ListBucket{
		Account: acct,
	}
	fmt.Println("Listing bucket...")
	bucketLister.List(uploadConfig.Bucket, uploadConfig.KeyPrefix, listFileName)
	//check data integrity
	fmt.Println("Checking data integrity...")
	CheckQrsync(cacheFileName, listFileName, uploadConfig.IgnoreDir, uploadConfig.KeyPrefix)
}
