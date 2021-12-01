package ali

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego/logs"
	"os"
)

// 阿里空间字段
type AliListBucket struct {
	DataCenter      string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Prefix          string
}

// 列举阿里空间的文件列表
func (this *AliListBucket) ListBucket(listResultFile string) (err error) {
	//open result file
	fp, openErr := os.Create(listResultFile)
	if openErr != nil {
		err = openErr
		return
	}
	defer fp.Close()
	bw := bufio.NewWriter(fp)
	//list bucket by prefix
	marker := ""
	limit := 1000
	prefixLen := len(this.Prefix)
	ossClient, clientErr := oss.New(this.DataCenter, this.AccessKeyId, this.AccessKeySecret)
	if clientErr != nil {
		err = clientErr
		return
	}

	ossBucket, bucketErr := ossClient.Bucket(this.Bucket)
	if bucketErr != nil {
		err = bucketErr
		return
	}

	maxRetryTimes := 5
	retryTimes := 1

	logs.Info("Listing the oss bucket...")
	for {
		listOptions := []oss.Option{
			oss.MaxKeys(limit),
			oss.Prefix(this.Prefix),
			oss.Marker(marker),
		}

		lbr, lbrErr := ossBucket.ListObjects(listOptions...)
		if lbrErr != nil {
			err = lbrErr
			logs.Error("Parse list result error,", "marker=[", marker, "]", lbrErr)
			if retryTimes <= maxRetryTimes {
				logs.Warning("Retry marker=", marker, "] for", retryTimes, "time...")
				retryTimes += 1
				continue
			} else {
				break
			}
		} else {
			retryTimes = 1
		}
		for _, object := range lbr.Objects {
			lmdTime := object.LastModified
			bw.WriteString(fmt.Sprintln(fmt.Sprintf("%s\t%d\t%d", object.Key[prefixLen:], object.Size, lmdTime.UnixNano()/100)))
		}
		if !lbr.IsTruncated {
			break
		}
		marker = lbr.NextMarker
	}
	fErr := bw.Flush()
	if fErr != nil {
		logs.Error("Write data to buffer writer failed", fErr)
		err = fErr
		return
	}
	return err
}
