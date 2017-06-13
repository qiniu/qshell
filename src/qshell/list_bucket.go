package qshell

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rsf"
	"qiniu/rpc"
)

/*
*@param bucket
*@param prefix
*@param marker
*@param listResultFile
*@return listError
 */
func ListBucket(mac *digest.Mac, bucket, prefix, marker, listResultFile string) (retErr error) {
	var listResultFh *os.File
	if listResultFile == "stdout" {
		listResultFh = os.Stdout
	} else {
		var openErr error
		//if marker not empty, continue the list
		if marker != "" {
			listResultFh, openErr = os.OpenFile(listResultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if openErr != nil {
				retErr = openErr
				logs.Error("Failed to open list result file `%s`", listResultFile)
				return
			}
		} else {
			listResultFh, openErr = os.Create(listResultFile)
			if openErr != nil {
				retErr = openErr
				logs.Error("Failed to open list result file `%s`", listResultFile)
				return
			}
		}
	}
	defer listResultFh.Close()
	bWriter := bufio.NewWriter(listResultFh)

	//get zone info
	bucketInfo, gErr := GetBucketInfo(mac, bucket)
	if gErr != nil {
		retErr = gErr
		logs.Error("Failed to get region info of bucket `%s`, %s", bucket, gErr)
		return
	}

	//set zone
	SetZone(bucketInfo.Region)

	//init
	client := rsf.New(mac)
	limit := 1000
	run := true
	maxRetryTimes := 5
	retryTimes := 1

	//start to list
	for run {
		entries, markerOut, listErr := client.ListPrefix(nil, bucket, prefix, marker, limit)
		if listErr != nil {
			if listErr == io.EOF {
				run = false
			} else {
				if v, ok := listErr.(*rpc.ErrorInfo); ok {
					logs.Error("List error for marker `%s`, %s", marker, v.Err)
				} else {
					logs.Error("List error for marker `%s`, %s", marker, listErr)
				}
				if retryTimes <= maxRetryTimes {
					logs.Warning("Retry list for marker `%s` for %d times", marker, retryTimes)
					retryTimes += 1
					continue
				} else {
					logs.Error("List failed too many times for marker `%s`", marker)
					break
				}
			}
		} else {
			retryTimes = 1
			if markerOut == "" {
				run = false
			} else {
				marker = markerOut
			}
		}

		//append entries
		for _, entry := range entries {
			lineData := fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%s\r\n",
				entry.Key, entry.Fsize, entry.Hash, entry.PutTime, entry.MimeType, entry.FileType, entry.EndUser)
			_, wErr := bWriter.WriteString(lineData)
			if wErr != nil {
				logs.Error("Write line data `%s` to list result file failed.", lineData)
			}
		}

		//flush
		fErr := bWriter.Flush()
		if fErr != nil {
			logs.Error("Flush data to list result file error", listErr)
		}
	}

	return
}
