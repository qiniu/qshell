package qshell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rsf"
	"qiniu/log"
	"qiniu/rpc"
)

type ListBucket struct {
	Account
}

/*
*@param bucket
*@param prefix
*@param marker
*@param listResultFile
*@return listError
 */
func (this *ListBucket) List(bucket, prefix, marker, listResultFile string) (retErr error) {
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
				log.Errorf("Failed to open list result file `%s`", listResultFile)
				return
			}
		} else {
			listResultFh, openErr = os.Create(listResultFile)
			if openErr != nil {
				retErr = openErr
				log.Errorf("Failed to open list result file `%s`", listResultFile)
				return
			}
		}
	}
	defer listResultFh.Close()
	bWriter := bufio.NewWriter(listResultFh)

	mac := digest.Mac{this.AccessKey, []byte(this.SecretKey)}

	//get zone info
	bucketInfo, gErr := GetBucketInfo(&mac, bucket)
	if gErr != nil {
		retErr = gErr
		log.Errorf("Failed to get region info of bucket `%s`", bucket)
		return
	}

	//set zone
	SetZone(bucketInfo.Region)

	//init
	client := rsf.New(&mac)
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
					log.Errorf("List error for marker `%s`, %s", marker, v.Err)
				} else {
					log.Errorf("List error for marker `%s`, %s", marker, listErr)
				}
				if retryTimes <= maxRetryTimes {
					log.Debugf("Retry list for marker `%s` for %d times", marker, retryTimes)
					retryTimes += 1
					continue
				} else {
					log.Errorf("List failed too many times for marker `%s`", marker)
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
			lineData := fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%s\r\n", entry.Key, entry.Fsize, entry.Hash, entry.PutTime, entry.MimeType, entry.EndUser)
			_, wErr := bWriter.WriteString(lineData)
			if wErr != nil {
				log.Errorf("Write line data `%s` to list result file failed.", lineData)
			}
		}

		//flush
		fErr := bWriter.Flush()
		if fErr != nil {
			log.Error("Flush data to list result file error", listErr)
		}
	}

	return
}
