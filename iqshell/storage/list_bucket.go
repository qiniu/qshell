package storage

import (
	"bufio"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/v2/iqshell/utils"
	"os"
	"os/signal"
	"strings"
	"time"
)

func filterByPuttime(putTime, startDate, endDate time.Time) bool {
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
}

func filterBySuffixes(key string, suffixes []string) bool {
	hasSuffix := false
	if len(suffixes) == 0 {
		hasSuffix = true
	}
	for _, s := range suffixes {
		if strings.HasSuffix(key, s) {
			hasSuffix = true
			break
		}
	}
	if hasSuffix {
		return true
	} else {
		return false
	}
}

func errorWarning(marker string, err error) {
	fmt.Fprintf(os.Stderr, "marker: %s\n", marker)
	fmt.Fprintf(os.Stderr, "listbucket Error: %v\n", err)
}

/*
*@param bucket
*@param prefix
*@param marker
*@param listResultFile
*@return listError
 */
func (m *BucketManager) ListFiles(bucket, prefix, marker, listResultFile string) (retErr error) {
	return m.ListBucket2(bucket, prefix, marker, listResultFile, "", time.Time{}, time.Time{}, nil, 20, false, false)
}

func (m *BucketManager) ListBucket2(bucket, prefix, marker, listResultFile, delimiter string, startDate, endDate time.Time, suffixes []string, maxRetry int, appendMode bool, readable bool) (retErr error) {
	lastMarker := marker

	defer func(lastMarker string) {
		if lastMarker != "" {
			fmt.Fprintf(os.Stderr, "Marker: %s\n", lastMarker)
		}
	}(lastMarker)

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		// 捕捉Ctrl-C, 退出下面列举的循环
		<-sigChan
		cancel()
		maxRetry = 0

		fmt.Printf("\nMarker: %s\n", lastMarker)
		os.Exit(1)
	}()

	var listResultFh *os.File

	if listResultFile == "" {
		listResultFh = os.Stdout
	} else {
		var openErr error
		var mode int

		if appendMode {
			mode = os.O_APPEND | os.O_RDWR
		} else {
			mode = os.O_CREATE | os.O_RDWR | os.O_TRUNC
		}
		listResultFh, openErr = os.OpenFile(listResultFile, mode, 0666)
		if openErr != nil {
			retErr = openErr
			logs.Error("Failed to open list result file `%s`", listResultFile)
			return
		}
		defer listResultFh.Close()
	}

	bWriter := bufio.NewWriter(listResultFh)

	notfilterTime := startDate.IsZero() && endDate.IsZero()
	notfilterSuffix := len(suffixes) == 0

	var c int
	for {
		if maxRetry >= 0 && c >= maxRetry {
			break
		}
		entries, lErr := m.ListBucketContext(ctx, bucket, prefix, delimiter, marker)

		if entries == nil && lErr == nil {
			// no data
			if lastMarker == "" {
				break
			} else {
				fmt.Fprintf(os.Stderr, "meet empty body when list not completed\n")
				continue
			}
		}
		if lErr != nil {
			retErr = lErr
			errorWarning(lastMarker, retErr)
			if maxRetry > 0 {
				c++
			}
			time.Sleep(1)
			continue
		}
		var fsizeValue interface{}

		for listItem := range entries {
			if listItem.Marker != lastMarker {
				lastMarker = listItem.Marker
			}
			if listItem.Item.IsEmpty() {
				continue
			}
			if readable {
				fsizeValue = utils.BytesToReadable(listItem.Item.Fsize)
			} else {
				fsizeValue = listItem.Item.Fsize
			}
			if notfilterSuffix && notfilterTime {
				lineData := fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d\t%s\r\n",
					listItem.Item.Key, fsizeValue, listItem.Item.Hash,
					listItem.Item.PutTime, listItem.Item.MimeType, listItem.Item.Type, listItem.Item.EndUser)
				_, wErr := bWriter.WriteString(lineData)
				if wErr != nil {
					retErr = wErr
					errorWarning(lastMarker, retErr)
				}

			} else {
				var hasSuffix = true
				var putTimeValid = true

				if !notfilterTime { // filter by putTime
					putTime := time.Unix(listItem.Item.PutTime/1e7, 0)
					putTimeValid = filterByPuttime(putTime, startDate, endDate)
				}
				if !notfilterSuffix {
					key := listItem.Item.Key
					hasSuffix = filterBySuffixes(key, suffixes)
				}

				if hasSuffix && putTimeValid {
					lineData := fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d\t%s\r\n",
						listItem.Item.Key, fsizeValue, listItem.Item.Hash,
						listItem.Item.PutTime, listItem.Item.MimeType, listItem.Item.Type, listItem.Item.EndUser)
					_, wErr := bWriter.WriteString(lineData)
					if wErr != nil {
						retErr = wErr
						errorWarning(lastMarker, retErr)
					}
				}
			}
		}
		fErr := bWriter.Flush()
		if fErr != nil {
			retErr = fErr
			errorWarning(lastMarker, retErr)
			if maxRetry > 0 {
				c++
			}
		}
		if lastMarker == "" {
			break
		} else {
			marker = lastMarker
		}
	}

	return
}
