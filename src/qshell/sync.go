package qshell

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"qiniu/api.v6/auth/digest"
	rio "qiniu/api.v6/resumable/io"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
	"strconv"
	"strings"
	"time"
)

//range get and chunk upload

const (
	RETRY_MAX_TIMES = 5
	RETRY_INTERVAL  = time.Second * 1
	HTTP_TIMEOUT    = time.Second * 10
)

type PutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

type SyncProgress struct {
	BlkCtxs   []rio.BlkputRet `json:"blk_ctxs"`
	Offset    int64           `json:"offset"`
	TotalSize int64           `json:"total_size"`
}

func Sync(mac *digest.Mac, srcResUrl, bucket, key, upHostIp string) (putRet PutRet, err error) {
	if exists, cErr := checkExists(mac, bucket, key); cErr != nil {
		err = cErr
		return
	} else if exists {
		err = errors.New("File with same key` already exists in bucket")
		return
	}

	syncProgress := SyncProgress{}
	//create sync id
	syncId := Md5Hex(fmt.Sprintf("%s:%s:%s", srcResUrl, bucket, key))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "sync")
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		return
	}

	progressFile := filepath.Join(storePath, fmt.Sprintf("%s.progress", syncId))
	if statInfo, statErr := os.Stat(progressFile); statErr == nil {
		//check file last modified time, if older than one week, ignore
		if statInfo.ModTime().Add(time.Hour * 24 * 5).After(time.Now()) {
			//try read old progress
			progressFh, openErr := os.Open(progressFile)
			if openErr == nil {
				decoder := json.NewDecoder(progressFh)
				decoder.Decode(&syncProgress)
				progressFh.Close()
			}
		}
	}

	//check offset valid or not
	if syncProgress.Offset%BLOCK_SIZE != 0 {
		logs.Info("Invalid offset from progress file,", syncProgress.Offset)
		syncProgress.Offset = 0
		syncProgress.TotalSize = 0
		syncProgress.BlkCtxs = make([]rio.BlkputRet, 0)
	}

	//check offset and blk ctxs
	if syncProgress.Offset != 0 && syncProgress.BlkCtxs != nil {
		if int(syncProgress.Offset/BLOCK_SIZE) != len(syncProgress.BlkCtxs) {
			logs.Info("Invalid offset and block contexts")
			syncProgress.Offset = 0
			syncProgress.TotalSize = 0
			syncProgress.BlkCtxs = make([]rio.BlkputRet, 0)
		}
	}

	//check blk ctxs, when no progress found
	if syncProgress.Offset == 0 || syncProgress.BlkCtxs == nil {
		syncProgress.Offset = 0
		syncProgress.TotalSize = 0
		syncProgress.BlkCtxs = make([]rio.BlkputRet, 0)
	}

	//get total size
	totalSize, hErr := getRemoteFileLength(srcResUrl)
	if hErr != nil {
		err = hErr
		return
	}

	if totalSize != syncProgress.TotalSize {
		if syncProgress.TotalSize != 0 {
			logs.Warning("Remote file length changed, progress file out of date")
		}
		syncProgress.Offset = 0
		syncProgress.TotalSize = totalSize
		syncProgress.BlkCtxs = make([]rio.BlkputRet, 0)
	}

	//get total block count
	totalBlkCnt := 0
	if totalSize%BLOCK_SIZE == 0 {
		totalBlkCnt = int(totalSize / BLOCK_SIZE)
	} else {
		totalBlkCnt = int(totalSize/BLOCK_SIZE) + 1
	}

	//init the range offset
	rangeStartOffset := syncProgress.Offset
	fromBlkIndex := int(rangeStartOffset / BLOCK_SIZE)

	lastBlock := false

	//create upload token
	policy := rs.PutPolicy{Scope: bucket}
	//token is valid for one year
	policy.Expires = 3600 * 24 * 365
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	uptoken := policy.Token(mac)
	putClient := rio.NewClient(uptoken, upHostIp)

	//range get and mkblk upload
	for blkIndex := fromBlkIndex; blkIndex < totalBlkCnt; blkIndex++ {
		if blkIndex == totalBlkCnt-1 {
			lastBlock = true
		}

		syncPercent := fmt.Sprintf("%.2f", float64(blkIndex+1)*100.0/float64(totalBlkCnt))
		logs.Info(fmt.Sprintf("Syncing block %d [%s%%] ...", blkIndex, syncPercent))
		blkCtx, pErr := rangeMkblkPipe(srcResUrl, totalSize, rangeStartOffset, BLOCK_SIZE, lastBlock, putClient)
		if pErr != nil {
			logs.Error(pErr.Error())
			time.Sleep(RETRY_INTERVAL)

			for retryTimes := 1; retryTimes <= RETRY_MAX_TIMES; retryTimes++ {
				logs.Info("Retrying %d time range & mkblk block [%d]", retryTimes, blkIndex)
				blkCtx, pErr = rangeMkblkPipe(srcResUrl, totalSize, rangeStartOffset, BLOCK_SIZE, lastBlock, putClient)
				if pErr != nil {
					logs.Error(pErr)
					//wait a interval and retry
					time.Sleep(RETRY_INTERVAL)
					continue
				} else {
					break
				}
			}
		}

		if pErr != nil {
			err = errors.New("Max retry reached and range & mkblk still failed, check your network")
			return
		}

		//advance range offset
		rangeStartOffset += BLOCK_SIZE

		syncProgress.BlkCtxs = append(syncProgress.BlkCtxs, blkCtx)
		syncProgress.Offset = rangeStartOffset

		rErr := recordProgress(progressFile, syncProgress)
		if rErr != nil {
			logs.Info(rErr.Error())
		}
	}

	//make file
	putExtra := rio.PutExtra{
		Progresses: syncProgress.BlkCtxs,
	}
	mkErr := rio.Mkfile(putClient, nil, &putRet, key, true, totalSize, &putExtra)
	if mkErr != nil {
		err = fmt.Errorf("Mkfile error, %s", mkErr.Error())
		return
	}

	//delete progress file
	os.Remove(progressFile)

	return
}

func rangeMkblkPipe(srcResUrl string, totalSize, rangeStartOffset, rangeBlockSize int64, lastBlock bool,
	putClient rpc.Client) (putRet rio.BlkputRet, err error) {
	//range get
	dReq, dReqErr := http.NewRequest("GET", srcResUrl, nil)
	if dReqErr != nil {
		err = fmt.Errorf("New request error, %s", dReqErr.Error())
		return
	}

	//proxyURL, _ := url.Parse("http://localhost:8888")

	//set range header
	rangeEndOffset := rangeStartOffset + rangeBlockSize - 1
	if lastBlock {
		rangeEndOffset = totalSize - 1
	}

	dReq.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rangeStartOffset, rangeEndOffset))

	//set client properties
	client := http.DefaultClient
	client.Timeout = time.Duration(HTTP_TIMEOUT)
	//client.Transport = &http.Transport{
	//	Proxy: http.ProxyURL(proxyURL),
	//}

	client.CheckRedirect = func(rReq *http.Request, rVias []*http.Request) (err error) {
		rReq.Header.Add("Range", dReq.Header.Get("Range"))
		return nil
	}

	//get response
	dResp, dRespErr := client.Do(dReq)
	if dRespErr != nil {
		err = fmt.Errorf("Get response error, %s", dRespErr.Error())
		return
	}
	defer dResp.Body.Close()

	//fmt.Println("-------------------")
	//fmt.Println(dResp.StatusCode)
	//for k, v := range dResp.Header {
	//	fmt.Println(k, ":", strings.Join(v, ","))
	//}

	//status error
	if dResp.StatusCode/100 != 2 {
		err = fmt.Errorf("Get resource error, %s", dResp.Status)
		return
	}

	//if not support range, go back and err
	if dResp.Header.Get("Content-Range") == "" {
		err = errors.New("Remote server not support range")
		return
	}

	//parse content-range
	contentRange := dResp.Header.Get("Content-Range")
	rangeSize, _ := parseContentRange(contentRange)

	//check ranged block size
	if !lastBlock && rangeSize != rangeBlockSize {
		err = errors.New("Block read error, only the last range block can has bytes less than <RangeBlockSize>")
		return
	}

	//read content
	buffer := bytes.NewBuffer(nil)
	cpCnt, cpErr := io.Copy(buffer, dResp.Body)
	if cpErr != nil || cpCnt != rangeSize {
		err = errors.New("Read range block response error, not fully read")
		return
	}

	//mkblk
	blkPutRet := rio.BlkputRet{}
	blockSize := int(rangeSize)
	blockDataReader := bytes.NewReader(buffer.Bytes())
	blockDataSize := buffer.Len()

	mkErr := rio.Mkblock(putClient, nil, &blkPutRet, blockSize, blockDataReader, blockDataSize)
	if mkErr != nil {
		err = fmt.Errorf("Mkblk error, %s", mkErr.Error())
		return
	}

	putRet = blkPutRet

	return
}

func recordProgress(progressFile string, syncProgress SyncProgress) (err error) {
	fh, openErr := os.Create(progressFile)
	if openErr != nil {
		err = fmt.Errorf("Open progress file %s error, %s", progressFile, openErr.Error())
		return
	}
	defer fh.Close()

	jsonBytes, mErr := json.Marshal(&syncProgress)
	if mErr != nil {
		err = fmt.Errorf("Marshal sync progress error, %s", mErr.Error())
		return
	}

	_, wErr := fh.Write(jsonBytes)
	if wErr != nil {
		err = fmt.Errorf("Write sync progress error, %s", wErr.Error())
	}

	return
}

func getRemoteFileLength(srcResUrl string) (totalSize int64, err error) {
	resp, respErr := http.Head(srcResUrl)
	if respErr != nil {
		err = fmt.Errorf("New head request failed, %s", respErr.Error())
		return
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		err = errors.New("Head request with no Content-Length found error")
		return
	}

	totalSize, _ = strconv.ParseInt(contentLength, 10, 64)

	return
}

func checkExists(mac *digest.Mac, bucket, key string) (exists bool, err error) {
	client := rs.NewMac(mac)
	entry, sErr := client.Stat(nil, bucket, key)
	if sErr != nil {
		if v, ok := sErr.(*rpc.ErrorInfo); !ok {
			err = fmt.Errorf("Check file exists error, %s", sErr.Error())
			return
		} else {
			if v.Code != 612 {
				err = fmt.Errorf("Check file exists error, %s", v.Err)
				return
			} else {
				exists = false
				return
			}
		}
	}

	if entry.Hash != "" {
		exists = true
	}

	return
}

//Content-Range: bytes 25538640-25538647/25538648
func parseContentRange(contentRange string) (rangeSize, totalSize int64) {
	contentRangeItems := strings.Split(contentRange, " ")
	sizeItems := strings.Split(contentRangeItems[1], "/")

	rangePartItems := strings.Split(sizeItems[0], "-")
	totalSize, _ = strconv.ParseInt(sizeItems[1], 10, 64)

	fromOffset, _ := strconv.ParseInt(rangePartItems[0], 10, 64)
	toOffset, _ := strconv.ParseInt(rangePartItems[1], 10, 64)

	rangeSize = toOffset - fromOffset + 1

	return
}
