package qshell

import (
	"bufio"
	"fmt"
	"github.com/qiniu/log"
	"github.com/yanunon/oss-go-api/oss"
	"os"
	"time"
)

type AliListBucket struct {
	DataCenter      string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Prefix          string
}

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
	prefixLen := len(this.Prefix)
	ossClient := oss.NewClient(this.DataCenter, this.AccessKeyId, this.AccessKeySecret, 0)
	maxRetryTimes := 5
	retryTimes := 1

	log.Info("Listing the oss bucket...")
	for {
		lbr, lbrErr := ossClient.GetBucket(this.Bucket, this.Prefix, marker, "", "")
		if lbrErr != nil {
			err = lbrErr
			log.Error("Parse list result error, ", "marker=[", marker, "]", lbrErr)
			if retryTimes <= maxRetryTimes {
				log.Debug("Retry marker=", marker, "] for ", retryTimes, "time...")
				retryTimes += 1
				continue
			} else {
				break
			}
		} else {
			retryTimes = 1
		}
		for _, object := range lbr.Contents {
			lmdTime, lmdPErr := time.Parse("2006-01-02T15:04:05.999Z", object.LastModified)
			if lmdPErr != nil {
				log.Error("Parse object last modified error, ", lmdPErr)
				lmdTime = time.Now()
			}
			bw.WriteString(fmt.Sprintln(fmt.Sprintf("%s\t%d\t%d", object.Key[prefixLen:], object.Size, lmdTime.UnixNano()/100)))
		}
		if !lbr.IsTruncated {
			break
		}
		marker = lbr.NextMarker
	}
	fErr := bw.Flush()
	if fErr != nil {
		log.Error("Write data to buffer writer failed", fErr)
		err = fErr
		return
	}
	return err
}
