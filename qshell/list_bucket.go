package qshell

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rsf"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"io"
	"os"
	"time"
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
	if listResultFile == "" {
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
		defer listResultFh.Close()
	}
	bWriter := bufio.NewWriter(listResultFh)

	//init
	client := rsf.New(mac)
	limit := 1000
	run := true
	maxRetryTimes := 5
	retryTimes := 1

	//start to list
	for run {
		entries, markerOut, listErr := client.ListPrefix(nil, bucket, prefix, marker, limit)
		limit = 1000
		if listErr != nil {
			limit = 1
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

func ListBucket2(mac *qbox.Mac, bucket, prefix, marker, listResultFile, delimiter string, startDate, endDate time.Time) (retErr error) {
	var listResultFh *os.File
	var listAll bool

	if listResultFile == "" {
		listResultFh = os.Stdout
	} else {
		var openErr error
		listResultFh, openErr = os.Create(listResultFile)
		if openErr != nil {
			retErr = openErr
			logs.Error("Failed to open list result file `%s`", listResultFile)
			return
		}
		defer listResultFh.Close()
	}

	if startDate.IsZero() && endDate.IsZero() {
		listAll = true
	}

	bWriter := bufio.NewWriter(listResultFh)

	bm := storage.NewBucketManager(mac, nil)

	for {
		lastMarker, err := listBucket2(bm, bucket, prefix, delimiter, marker, bWriter, listAll, func(putTime time.Time) bool {
			switch {
			case startDate.IsZero() && endDate.IsZero():
				return true
			case !startDate.IsZero() && endDate.IsZero() && putTime.After(startDate):
				return true
			case !endDate.IsZero() && startDate.IsZero() && putTime.Before(endDate):
				return true
			case putTime.After(startDate) && putTime.Before(endDate):
				return true
			default:
				return false
			}
		})

		if err != nil {
			marker = lastMarker
			fmt.Fprintf(os.Stderr, "ListBucket: %v\n", err)
			fmt.Fprintf(os.Stderr, "marker: %v\n", marker)
		} else {
			if lastMarker == "" {
				break
			} else {
				marker = lastMarker
			}
		}
	}
	return
}

func listBucket2(bm *storage.BucketManager, bucket, prefix, delimiter, marker string, out *bufio.Writer,
	listAll bool, filter func(time.Time) bool) (string, error) {
	var lastMarker string

	entries, err := bm.ListBucket(bucket, prefix, delimiter, marker)

	for listItem := range entries {
		if listItem.Marker != lastMarker {
			lastMarker = listItem.Marker
		}

		putTime := time.Unix(listItem.Item.PutTime/1e7, 0)
		if listAll || filter(putTime) {
			lineData := fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%s\r\n",
				listItem.Item.Key, listItem.Item.Fsize, listItem.Item.Hash,
				listItem.Item.PutTime, listItem.Item.MimeType, listItem.Item.Type, listItem.Item.EndUser)
			_, wErr := out.WriteString(lineData)
			if wErr != nil {
				return lastMarker, wErr
			}
		}
	}
	fErr := out.Flush()
	if fErr != nil {
		return lastMarker, fErr
	}
	return lastMarker, err
}
