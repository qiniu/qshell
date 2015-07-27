package qshell

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/log"
	"github.com/qiniu/rpc"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
{
	"dest_dir"		:	"/Users/jemy/Backup",
	"bucket"		:	"test-bucket",
	"domain"		:	"<Your bucket domain>",
	"access_key"	:	"<Your AccessKey>",
	"secret_key"	:	"<Your SecretKey>",
	"is_private"	:	false,
	"prefix"		:	"demo/"
}
*/

const (
	MIN_DOWNLOAD_THREAD_COUNT = 1
	MAX_DOWNLOAD_THREAD_COUNT = 100
)

type DownloadConfig struct {
	DestDir   string `json:"dest_dir"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	IsPrivate bool   `json:"is_private"`
	Prefix    string `json:"prefix,omitempty"`
}

func QiniuDownload(threadCount int, downloadConfigFile string) {
	cnfFp, err := os.Open(downloadConfigFile)
	if err != nil {
		log.Error("Open download config file", downloadConfigFile, "failed,", err)
		return
	}
	defer cnfFp.Close()
	cnfData, err := ioutil.ReadAll(cnfFp)
	if err != nil {
		log.Error("Read download config file error", err)
		return
	}
	downConfig := DownloadConfig{}
	cnfErr := json.Unmarshal(cnfData, &downConfig)
	if cnfErr != nil {
		log.Error("Parse download config error", err)
		return
	}
	cnfJson, _ := json.Marshal(&downConfig)
	jobId := fmt.Sprintf("%x", md5.Sum(cnfJson))
	jobListName := fmt.Sprintf("%s.list.txt", jobId)
	acct := Account{
		AccessKey: downConfig.AccessKey,
		SecretKey: downConfig.SecretKey,
	}
	bLister := ListBucket{
		Account: acct,
	}
	log.Debug("List bucket...")
	listErr := bLister.List(downConfig.Bucket, downConfig.Prefix, jobListName)
	if listErr != nil {
		log.Error("List bucket error", listErr)
		return
	}
	listFp, openErr := os.Open(jobListName)
	if openErr != nil {
		log.Error("Open list file error", openErr)
		return
	}
	defer listFp.Close()
	listScanner := bufio.NewScanner(listFp)
	listScanner.Split(bufio.ScanLines)
	downWorkGroup := sync.WaitGroup{}
	downCounter := 0

	threadThresold := threadCount + 1
	for listScanner.Scan() {
		downCounter += 1
		if downCounter%threadThresold == 0 {
			downWorkGroup.Wait()
		}
		line := strings.TrimSpace(listScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) > 2 {
			fileKey := items[0]
			fileSize, _ := strconv.ParseInt(items[1], 10, 64)
			//not backup yet
			if !checkLocalDuplicate(downConfig.DestDir, fileKey, fileSize) {
				downWorkGroup.Add(1)
				go func() {
					defer downWorkGroup.Done()
					downloadFile(downConfig, fileKey)
				}()
			}
		}
	}
	downWorkGroup.Wait()
	fmt.Println("All downloaded!")
}

func checkLocalDuplicate(destDir string, fileKey string, fileSize int64) bool {
	dup := false
	filePath := filepath.Join(destDir, fileKey)
	fStat, statErr := os.Stat(filePath)
	if statErr == nil {
		//exist, check file size
		localFileSize := fStat.Size()
		if localFileSize == fileSize {
			dup = true
		}
	}
	return dup
}

func downloadFile(downConfig DownloadConfig, fileKey string) {
	localFilePath := filepath.Join(downConfig.DestDir, fileKey)
	ldx := strings.LastIndex(localFilePath, string(os.PathSeparator))
	if ldx != -1 {
		localFileDir := localFilePath[:ldx]
		err := os.MkdirAll(localFileDir, 0775)
		if err != nil {
			log.Error("MkdirAll failed for", localFileDir)
			return
		}
	}
	fmt.Println("Downloading", fileKey, "=>", localFilePath, "...")
	downUrl := strings.Join([]string{downConfig.Domain, fileKey}, "/")
	if downConfig.IsPrivate {
		now := time.Now().Add(time.Second * 3600 * 24)
		downUrl = fmt.Sprintf("%s?e=%d", downUrl, now.Unix())
		mac := digest.Mac{downConfig.AccessKey, []byte(downConfig.SecretKey)}
		token := digest.Sign(&mac, []byte(downUrl))
		downUrl = fmt.Sprintf("%s&token=%s", downUrl, token)
	}
	resp, respErr := rpc.DefaultClient.Get(nil, downUrl)
	if respErr != nil {
		log.Error("Download", fileKey, "failed by url", downUrl)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		localFp, openErr := os.OpenFile(localFilePath, os.O_CREATE|os.O_WRONLY, 0666)
		if openErr != nil {
			log.Error("Open local file", localFilePath, "failed")
			return
		}
		defer localFp.Close()
		_, err := io.Copy(localFp, resp.Body)
		if err != nil {
			log.Error("Download", fileKey, "failed", err)
		}
	} else {
		log.Error("Download", fileKey, "failed by url", downUrl)
	}
}
