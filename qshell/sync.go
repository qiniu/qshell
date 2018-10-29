package qshell

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/storage"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

type ProgressRecorder struct {
	BlkCtxs      []storage.BlkputRet `json:"blk_ctxs"`
	Offset       int64               `json:"offset"`
	TotalSize    int64               `json:"total_size"`
	LastModified int                 `json:"last_modified"` // 上传文件的modification time
	FilePath     string              // 断点续传记录保存文件
}

func NewProgressRecorder(filePath string) *ProgressRecorder {
	p := new(ProgressRecorder)
	p.FilePath = filePath
	p.BlkCtxs = make([]storage.BlkputRet, 0)
	return p
}

func ProgressFileFromUrl(srcResUrl, bucket, key string) (progressFile string, err error) {

	//create sync id
	syncId := Md5Hex(fmt.Sprintf("%s:%s:%s", srcResUrl, bucket, key))

	//local storage path
	storePath := filepath.Join(QShellRootPath, ".qshell", "sync")
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		logs.Error("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		err = mkdirErr
		return
	}

	progressFile = filepath.Join(storePath, fmt.Sprintf("%s.progress", syncId))
	return
}

func (p *ProgressRecorder) Recover() (err error) {
	if statInfo, statErr := os.Stat(p.FilePath); statErr == nil {
		//check file last modified time, if older than one week, ignore
		if statInfo.ModTime().Add(time.Hour * 24 * 5).After(time.Now()) {
			//try read old progress
			progressFh, openErr := os.Open(p.FilePath)
			if openErr != nil {
				err = openErr
				return
			}
			decoder := json.NewDecoder(progressFh)
			decoder.Decode(p)
			progressFh.Close()
		}
	}
	return
}

func (p *ProgressRecorder) RecoverFromUrl(srcResUrl, bucket, key string) (err error) {

	progressFile, pErr := ProgressFileFromUrl(srcResUrl, bucket, key)
	if err != nil {
		err = pErr
		return
	}
	p.FilePath = progressFile
	return p.Recover()
}

func (p *ProgressRecorder) Reset() {
	p.Offset = 0
	p.TotalSize = 0
	p.BlkCtxs = make([]storage.BlkputRet, 0)
}

func (p *ProgressRecorder) CheckValid(fileSize int64, lastModified int) {

	//check offset valid or not
	if p.Offset%BLOCK_SIZE != 0 {
		logs.Info("Invalid offset from progress file,", p.Offset)
		p.Reset()
		return
	}

	//check offset and blk ctxs
	if p.Offset != 0 && p.BlkCtxs != nil {
		if int(p.Offset/BLOCK_SIZE) != len(p.BlkCtxs) {
			logs.Info("Invalid offset and block contexts")
			p.Reset()
			return
		}
	}

	//check blk ctxs, when no progress found
	if p.Offset == 0 || p.BlkCtxs == nil {
		p.Reset()
		return
	}
	if fileSize != p.TotalSize {
		if p.TotalSize != 0 {
			logs.Warning("Remote file length changed, progress file out of date")
		}
		p.Offset = 0
		p.TotalSize = fileSize
		p.BlkCtxs = make([]storage.BlkputRet, 0)
		return
	}
	if len(p.BlkCtxs) > 0 {
		if lastModified != 0 && p.LastModified != lastModified {
			p.Reset()
		}
	}
}

func (p *ProgressRecorder) RecordProgress() (err error) {
	fh, openErr := os.Create(p.FilePath)
	if openErr != nil {
		err = fmt.Errorf("Open progress file %s error, %s", p.FilePath, openErr.Error())
		return
	}
	defer fh.Close()

	jsonBytes, mErr := json.Marshal(p)
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

func (m *BucketManager) CheckExists(bucket, key string) (exists bool, err error) {
	entry, sErr := m.Stat(bucket, key)
	if sErr != nil {
		if v, ok := sErr.(*storage.ErrorInfo); !ok {
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

type SputRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

func (m *BucketManager) Sync(srcResUrl, bucket, key string) (putRet SputRet, err error) {

	exists, cErr := m.CheckExists(bucket, key)
	if cErr != nil {
		err = cErr
		return
	}
	if exists {
		err = errors.New("File with same key` already exists in bucket")
		return
	}
	//get total size
	totalSize, hErr := getRemoteFileLength(srcResUrl)
	if hErr != nil {
		err = hErr
		return
	}
	progressFile, fErr := ProgressFileFromUrl(srcResUrl, bucket, key)
	if err != nil {
		err = fErr
		return
	}
	syncProgress := NewProgressRecorder(progressFile)
	syncProgress.RecoverFromUrl(srcResUrl, bucket, key)
	syncProgress.CheckValid(totalSize, 0)

	//get total block count
	fmt.Println("totalSize: ", totalSize)
	totalBlkCnt := storage.BlockCount(totalSize)

	//init the range offset
	rangeStartOffset := syncProgress.Offset
	fromBlkIndex := int(rangeStartOffset / BLOCK_SIZE)

	lastBlock := false

	//create upload token
	policy := storage.PutPolicy{Scope: bucket}
	//token is valid for one year
	policy.Expires = 3600 * 24 * 365
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	uptoken := policy.UploadToken(m.GetMac())
	ctx := context.Background()

	resumeUploader := NewResumeUploader(nil)
	ak, bucket, err := getAkBucketFromUploadToken(uptoken)
	if err != nil {
		return
	}
	upHost, eErr := resumeUploader.UpHost(ak, bucket)
	if err != nil {
		err = eErr
		return
	}
	//range get and mkblk upload
	var bf *bytes.Buffer
	var blockSize = BLOCK_SIZE
	for blkIndex := fromBlkIndex; blkIndex < totalBlkCnt; blkIndex++ {
		if blkIndex == totalBlkCnt-1 {
			lastBlock = true
		}

		syncPercent := fmt.Sprintf("%.2f", float64(blkIndex+1)*100.0/float64(totalBlkCnt))
		logs.Info(fmt.Sprintf("Syncing block %d [%s%%] ...", blkIndex, syncPercent))

		var blkCtx storage.BlkputRet
		var retryTimes int
		var rErr error
		for {
			bf, rErr = getRange(srcResUrl, totalSize, rangeStartOffset, BLOCK_SIZE, lastBlock)
			if rErr != nil && retryTimes >= RETRY_MAX_TIMES {
				err = errors.New(strings.Join([]string{"Get range block data failed: ", rErr.Error()}, ""))
				return
			}
			if rErr == nil {
				break
			}
			logs.Error(rErr.Error())
			time.Sleep(RETRY_INTERVAL)
			logs.Info("Retrying %d time get range for block [%d]", retryTimes, blkIndex)
			retryTimes++
		}
		data := bf.Bytes()
		if lastBlock {
			blockSize = len(data)
		}
		retryTimes = 0
		for {
			pErr := resumeUploader.Mkblk(ctx, uptoken, upHost, &blkCtx, blockSize, bytes.NewReader(data), len(data))
			if pErr != nil && retryTimes >= RETRY_MAX_TIMES {
				err = pErr
				return
			}
			if pErr == nil {
				break
			}
			logs.Error(pErr.Error())
			time.Sleep(RETRY_INTERVAL)

			logs.Info("Retrying %d time mkblk for block [%d]", retryTimes, blkIndex)
			retryTimes++
		}
		//advance range offset
		rangeStartOffset += BLOCK_SIZE

		syncProgress.BlkCtxs = append(syncProgress.BlkCtxs, blkCtx)
		syncProgress.Offset = rangeStartOffset

		sErr := syncProgress.RecordProgress()
		if sErr != nil {
			logs.Info(rErr.Error())
		}
	}

	//make file
	putExtra := storage.RputExtra{
		Progresses: syncProgress.BlkCtxs,
	}
	mkErr := resumeUploader.Mkfile(ctx, uptoken, upHost, &putRet, key, true, totalSize, &putExtra)
	if mkErr != nil {
		err = fmt.Errorf("Mkfile error, %s", mkErr.Error())
		return
	}

	//delete progress file
	os.Remove(progressFile)

	return
}

func getRange(srcResUrl string, totalSize, rangeStartOffset, rangeBlockSize int64, lastBlock bool) (data *bytes.Buffer, err error) {
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

	return buffer, nil
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
