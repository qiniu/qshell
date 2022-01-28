package upload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/api"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	resumeV2MaxPart = 10000
	httpTimeout     = time.Second * 10
)

type conveyor struct {
	cfg *storage.Config
}

func newConveyorUploader(cfg *storage.Config) Uploader {
	return &conveyor{
		cfg: cfg,
	}
}

func (c *conveyor) upload(info ApiInfo) (ret ApiResult, err error) {

	// 检查 Host
	if len(info.UpHost) == 0 {
		acc, gErr := workspace.GetAccount()
		if gErr != nil {
			err = fmt.Errorf("sync get account error:%v", gErr)
			return
		}
		info.UpHost, err = getUpHost(c.cfg, acc.AccessKey, info.ToBucket)
		if err != nil {
			err = fmt.Errorf("sync get up host error:%v", err)
			return
		}
	}

	ctx := workspace.GetContext()
	progressFile, fErr := ProgressFileFromUrl(info.FilePath, info.ToBucket, info.SaveKey)
	if fErr != nil {
		err = fErr
		return
	}
	recorder := api.NewProgressRecorder(progressFile)
	if rErr := recorder.Recover(); rErr != nil {
		log.WarningF("sync progress recover error:", rErr)
	}
	recorder.CheckValid(info.FileSize, 0, info.UseResumeV2)
	recorder.TotalSize = info.FileSize

	uploader := api.NewResume(api.ResumeInfo{
		UpHost:   info.UpHost,
		Bucket:   info.ToBucket,
		UpToken:  info.TokenProvider(),
		Key:      info.SaveKey,
		Recorder: recorder,
		Cfg:      nil,
	}, info.UseResumeV2)
	uploader = api.NewRetryResume(uploader, info.TryTimes, info.TryInterval)

	// 1. 初始化服务
	err = uploader.InitServer(ctx)
	if err != nil {
		return
	}

	// 2. 上传文件分片
	var blockSize = int64(data.BLOCK_SIZE)
	if info.UseResumeV2 {
		// 检查块大小是否满足实际需求
		maxParts := int64(resumeV2MaxPart)
		if blockSize*maxParts < info.FileSize {
			blockSize = (info.FileSize + maxParts - 1) / maxParts
		}
	}

	totalBlkCnt := storage.BlockCount(info.FileSize) //range get and mkblk upload
	rangeStartOffset := recorder.Offset              //init the range offset
	fromBlkIndex := int(rangeStartOffset / data.BLOCK_SIZE)
	var bf *bytes.Buffer
	for blkIndex := fromBlkIndex; blkIndex < totalBlkCnt; blkIndex++ {

		syncPercent := fmt.Sprintf("%.2f", float64(blkIndex+1)*100.0/float64(totalBlkCnt))
		log.DebugF(fmt.Sprintf("Syncing block %d [%s%%] ...", blkIndex, syncPercent))

		// 2.1 获取上传数据
		var retryTimes int
		for {
			bf, err = getRange(info.FilePath, info.FileSize, rangeStartOffset, blockSize)
			if err != nil && retryTimes >= info.TryTimes {
				err = errors.New(strings.Join([]string{"sync Get range block data failed: ", err.Error()}, ""))
				return
			}
			if err == nil {
				break
			}
			time.Sleep(info.TryInterval)
			log.DebugF("sync Retrying %d time get range for block [%d] for error:%v", retryTimes, blkIndex, err)
			retryTimes++
		}
		dataBytes := bf.Bytes()

		// 2.2 上传数据到云存储
		err = uploader.UploadBlock(ctx, 0, dataBytes)
		if err != nil {
			return
		}

		//advance range offset
		rangeStartOffset += data.BLOCK_SIZE
		if sErr := recorder.RecordProgress(); sErr != nil {
			log.WarningF("sync save record progress error:%v", sErr)
		}
	}

	// 3. 合并文件
	err = uploader.Complete(ctx, &ret)
	if err != nil {
		err = fmt.Errorf("sync complete error:%v", err)
		return
	}

	//delete progress file
	if rErr := os.Remove(progressFile); rErr != nil {
		log.WarningF("sync remove record progress error:%v", rErr)
	}

	return
}

func ProgressFileFromUrl(srcResUrl, bucket, key string) (progressFile string, err error) {

	//create sync id
	syncId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", srcResUrl, bucket, key))

	//local storage path
	QShellRootPath := workspace.GetWorkspace()
	if QShellRootPath == "" {
		err = fmt.Errorf("empty root path\n")
		return
	}
	storePath := filepath.Join(QShellRootPath, ".qshell", "sync")
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		err = fmt.Errorf("sync Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		return
	}

	progressFile = filepath.Join(storePath, fmt.Sprintf("%s.progress", syncId))
	return
}

func getRange(srcResUrl string, totalSize, rangeStartOffset, rangeBlockSize int64) (data *bytes.Buffer, err error) {
	//range get
	dReq, dReqErr := http.NewRequest("GET", srcResUrl, nil)
	if dReqErr != nil {
		err = fmt.Errorf("New request error, %s", dReqErr.Error())
		return
	}

	//set range header
	rangeEndOffset := rangeStartOffset + rangeBlockSize - 1
	if rangeEndOffset >= totalSize {
		rangeEndOffset = totalSize - 1
	}

	dReq.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rangeStartOffset, rangeEndOffset))

	//set client properties
	client := http.DefaultClient
	client.Timeout = time.Duration(httpTimeout)
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

	//status error
	if dResp.StatusCode/100 != 2 {
		err = fmt.Errorf("Get resource error, %s", dResp.Status)
		return
	}

	//if not support range, go back and err
	if dResp.Header.Get("Content-Range") == "" {
		err = errors.New("sync Remote server not support range")
		return
	}

	//parse content-range
	contentRange := dResp.Header.Get("Content-Range")
	rangeSize, _ := parseContentRange(contentRange)

	//check ranged block size
	if rangeSize != (rangeEndOffset - rangeStartOffset + 1) {
		err = errors.New("sync Block read error, only the last range block can has bytes less than <RangeBlockSize>")
		return
	}

	//read content
	buffer := bytes.NewBuffer(nil)
	cpCnt, cpErr := io.Copy(buffer, dResp.Body)
	if cpErr != nil || cpCnt != rangeSize {
		err = errors.New("sync Read range block response error, not fully read")
		return
	}

	return buffer, nil
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

func getUpHost(cfg *storage.Config, ak, bucket string) (upHost string, err error) {

	var zone *storage.Zone
	if cfg.Zone != nil {
		zone = cfg.Zone
	} else {
		if v, zoneErr := storage.GetZone(ak, bucket); zoneErr != nil {
			err = zoneErr
			return
		} else {
			zone = v
		}
	}

	scheme := "http://"
	if cfg.UseHTTPS {
		scheme = "https://"
	}

	host := zone.SrcUpHosts[0]
	if cfg.UseCdnDomains {
		host = zone.CdnUpHosts[0]
	}

	upHost = fmt.Sprintf("%s%s", scheme, host)
	return
}
