package upload

import (
	"bytes"
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
	resumeV2MinChunkSize = 1024 * 1024
	resumeV2MaxPart      = 10000
	httpTimeout          = time.Second * 60
)

type conveyor struct {
	cfg *storage.Config
}

func newConveyorUploader(cfg *storage.Config) Uploader {
	return &conveyor{
		cfg: cfg,
	}
}

func (c *conveyor) upload(info *ApiInfo) (ret *ApiResult, err *data.CodeError) {
	log.DebugF("conveyor upload:%s => [%s:%s]", info.FilePath, info.ToBucket, info.SaveKey)

	// 检查 Host
	if len(info.UpHost) == 0 {
		acc, gErr := workspace.GetAccount()
		if gErr != nil {
			err = data.NewEmptyError().AppendDescF("sync get account error:%v", gErr)
			return
		}
		info.UpHost, err = getUpHost(c.cfg, acc.AccessKey, info.ToBucket)
		if err != nil {
			err = data.NewEmptyError().AppendDescF("sync get up host error:%v", err)
			return
		}
	} else {
		info.UpHost = utils.Endpoint(c.cfg.UseHTTPS, info.UpHost)
	}

	ctx := workspace.GetContext()
	progressFile, fErr := ProgressFileFromUrl(info.FilePath, info.ToBucket, info.SaveKey)
	if fErr != nil {
		err = fErr
		return
	}
	recorder := api.NewProgressRecorder(progressFile)
	if rErr := recorder.Recover(); rErr != nil {
		log.WarningF("sync progress recover error:%v", rErr)
	}
	recorder.CheckValid(info.LocalFileSize, 0, info.UseResumeV2)
	recorder.TotalSize = info.LocalFileSize

	if info.Progress != nil {
		info.Progress.SetFileSize(info.LocalFileSize)
		info.Progress.Start()
	}

	uploader := api.NewResume(api.ResumeInfo{
		UpHost:        info.UpHost,
		Bucket:        info.ToBucket,
		TokenProvider: info.TokenProvider,
		Key:           info.SaveKey,
		Recorder:      recorder,
		Cfg:           nil,
	}, info.UseResumeV2)
	uploader = api.NewRetryResume(uploader, info.TryTimes, info.TryInterval)

	// 1. 初始化服务
	err = uploader.InitServer(ctx)
	if err != nil {
		return
	}

	// 2. 上传文件分片
	var blockSize = info.ChunkSize
	if blockSize < resumeV2MinChunkSize {
		blockSize = int64(data.BLOCK_SIZE)
	}

	if info.UseResumeV2 {
		// 检查块大小是否满足实际需求
		maxParts := int64(resumeV2MaxPart)
		if blockSize*maxParts < info.LocalFileSize {
			blockSize = (info.LocalFileSize + maxParts - 1) / maxParts
		}
	}

	totalBlkCnt := storage.BlockCount(info.LocalFileSize) //range get and mkblk upload
	rangeStartOffset := recorder.Offset                   //init the range offset
	fromBlkIndex := int(rangeStartOffset / data.BLOCK_SIZE)

	if info.Progress != nil {
		info.Progress.SendSize(rangeStartOffset)
	}

	var bf *bytes.Buffer
	for blkIndex := fromBlkIndex; blkIndex < totalBlkCnt; blkIndex++ {
		log.DebugF("") // 此处仅为日志换行
		log.DebugF("Syncing block %d ...", blkIndex)

		// 2.1 获取上传数据
		var retryTimes int
		for {
			bf, err = getRange(info.FilePath, info.LocalFileSize, rangeStartOffset, blockSize)
			if err != nil && retryTimes >= info.TryTimes {
				err = data.NewEmptyError().AppendDesc(strings.Join([]string{"sync Get range block data failed: ", err.Error()}, ""))
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
		} else {
			if info.Progress != nil {
				info.Progress.SendSize(int64(len(dataBytes)))
			}
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
		err = data.NewEmptyError().AppendDescF("sync complete error:%v", err)
		return
	} else {
		if info.Progress != nil {
			info.Progress.End()
		}
	}

	//delete progress file
	if rErr := os.Remove(progressFile); rErr != nil {
		log.WarningF("sync remove record progress error:%v", rErr)
	}

	return
}

func ProgressFileFromUrl(srcResUrl, bucket, key string) (progressFile string, err *data.CodeError) {

	//create sync id
	syncId := utils.Md5Hex(fmt.Sprintf("%s:%s:%s", srcResUrl, bucket, key))

	//local storage path
	QShellRootPath := workspace.GetWorkspace()
	if QShellRootPath == "" {
		err = data.NewEmptyError().AppendDescF("empty root path\n")
		return
	}
	storePath := filepath.Join(QShellRootPath, ".qshell", "sync")
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		err = data.NewEmptyError().AppendDescF("sync Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		return
	}

	progressFile = filepath.Join(storePath, fmt.Sprintf("%s.progress", syncId))
	return
}

func getRange(srcResUrl string, totalSize, rangeStartOffset, rangeBlockSize int64) (buffer *bytes.Buffer, err *data.CodeError) {
	//range get
	dReq, dReqErr := http.NewRequest("GET", srcResUrl, nil)
	if dReqErr != nil {
		err = data.NewEmptyError().AppendDescF("New request error, %s", dReqErr.Error())
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
	client.Timeout = httpTimeout
	//client.Transport = &http.Transport{
	//	Proxy: http.ProxyURL(proxyURL),
	//}

	client.CheckRedirect = func(rReq *http.Request, rVias []*http.Request) error {
		rReq.Header.Add("Range", dReq.Header.Get("Range"))
		return nil
	}

	//get response
	dResp, dRespErr := client.Do(dReq)
	if dRespErr != nil {
		err = data.NewEmptyError().AppendDescF("Get response error, %s", dRespErr.Error())
		return
	}
	defer dResp.Body.Close()

	//status error
	if dResp.StatusCode/100 != 2 {
		err = data.NewEmptyError().AppendDescF("Get resource error, %s", dResp.Status)
		return
	}

	//if not support range, go back and err
	if dResp.Header.Get("Content-Range") == "" {
		err = data.NewEmptyError().AppendDesc("sync Remote server not support range")
		return
	}

	//parse content-range
	contentRange := dResp.Header.Get("Content-Range")
	rangeSize, _ := parseContentRange(contentRange)

	//check ranged block size
	if rangeSize != (rangeEndOffset - rangeStartOffset + 1) {
		err = data.NewEmptyError().AppendDesc("sync Block read error, only the last range block can has bytes less than <RangeBlockSize>")
		return
	}

	//read content
	buffer = bytes.NewBuffer(nil)
	cpCnt, cpErr := io.Copy(buffer, dResp.Body)
	if cpErr != nil || cpCnt != rangeSize {
		err = data.NewEmptyError().AppendDescF("sync Read range block response error, not fully read:%v", cpErr)
		return
	}

	return buffer, nil
}

// Content-Range: bytes 25538640-25538647/25538648
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

func getUpHost(cfg *storage.Config, ak, bucket string) (upHost string, err *data.CodeError) {

	var zone *storage.Zone
	if cfg.Zone != nil {
		zone = cfg.Zone
	} else {
		if v, zoneErr := storage.GetZone(ak, bucket); zoneErr != nil {
			err = data.ConvertError(zoneErr)
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
