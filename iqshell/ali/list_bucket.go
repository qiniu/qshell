package ali

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type ListBucketInfo struct {
	DataCenter string
	Bucket     string
	AccessKey  string
	SecretKey  string
	Prefix     string
	SaveToFile string
}

func (info *ListBucketInfo) Check() *data.CodeError {
	if len(info.DataCenter) == 0 {
		return alert.CannotEmptyError("DataCenter", "")
	}
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.AccessKey) == 0 {
		return alert.CannotEmptyError("AccessKeyId", "")
	}
	if len(info.SecretKey) == 0 {
		return alert.CannotEmptyError("AccessKeySecret", "")
	}
	if len(info.SaveToFile) == 0 {
		return alert.CannotEmptyError("ListBucketResultFile", "")
	}
	return nil
}

// ListBucket
// 列举阿里空间中的文件列表
func ListBucket(cfg *iqshell.Config, info ListBucketInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	// open result file
	fp, err := os.Create(info.SaveToFile)
	if err != nil {
		log.Error("create file error:", err)
		data.SetCmdStatusError()
		return
	}
	defer func(fp *os.File) {
		err := fp.Close()
		if err != nil {
			log.Error("file close error:", err)
			data.SetCmdStatusError()
		}
	}(fp)

	bw := bufio.NewWriter(fp)
	ossClient, err := oss.New(info.DataCenter, info.AccessKey, info.SecretKey)
	if err != nil {
		log.Error("create oss client error:", err)
		data.SetCmdStatusError()
		return
	}

	ossBucket, err := ossClient.Bucket(info.Bucket)
	if err != nil {
		log.Error("create oss bucket error:", err)
		data.SetCmdStatusError()
		return
	}

	log.Info("Listing the oss bucket...")

	var (
		marker        = ""
		limit         = 1000
		retryTimes    = 1
		maxRetryTimes = 5
		prefixLen     = len(info.Prefix)
	)
	for {
		lbr, lErr := ossBucket.ListObjects(oss.MaxKeys(limit), oss.Prefix(info.Prefix), oss.Marker(marker))
		if lErr != nil {
			log.Error("Parse list result error,", "marker=[", marker, "]", lErr)
			if retryTimes <= maxRetryTimes {
				log.Warning("Retry marker=", marker, "] for", retryTimes, "time...")
				retryTimes += 1
				continue
			} else {
				data.SetCmdStatusError()
				break
			}
		} else {
			retryTimes = 1
		}

		for _, object := range lbr.Objects {
			lmdTime := object.LastModified
			if _, e := bw.WriteString(fmt.Sprintf("%s\t%d\t%d\n", object.Key[prefixLen:], object.Size, lmdTime.UnixNano()/100)); e != nil {
				log.ErrorF("write result to file:%s error:%v", info.SaveToFile, e)
			}
		}

		if !lbr.IsTruncated {
			break
		}

		marker = lbr.NextMarker
	}

	fErr := bw.Flush()
	if fErr != nil {
		log.Error("Write data to buffer writer failed", fErr)
		data.SetCmdStatusError()
		return
	}

	log.Info("List bucket done!")
}
