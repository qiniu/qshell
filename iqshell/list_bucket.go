package iqshell

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
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
	return m.ListBucket2(bucket, prefix, marker, listResultFile, "", time.Time{}, time.Time{}, nil, 5)
}

func (m *BucketManager) ListBucket2(bucket, prefix, marker, listResultFile, delimiter string, startDate, endDate time.Time, suffixes []string, maxRetry int) (retErr error) {

	var listResultFh *os.File

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

	bWriter := bufio.NewWriter(listResultFh)

	lastMarker := marker
	notfilterTime := startDate.IsZero() && endDate.IsZero()
	notfilterSuffix := len(suffixes) == 0

	for i := 0; i < maxRetry; i++ {
		entries, lErr := m.ListBucket(bucket, prefix, delimiter, marker)

		if entries == nil && lErr == nil {
			// no data
			return
		}
		if retErr != nil {
			errorWarning(lastMarker, retErr)
		}

		for listItem := range entries {
			if listItem.Marker != lastMarker {
				lastMarker = listItem.Marker
			}
			if listItem.Item.IsEmpty() {
				continue
			}
			if notfilterSuffix && notfilterTime {
				lineData := fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%s\r\n",
					listItem.Item.Key, listItem.Item.Fsize, listItem.Item.Hash,
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
					lineData := fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%s\r\n",
						listItem.Item.Key, listItem.Item.Fsize, listItem.Item.Hash,
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
		}
		if lastMarker == "" {
			break
		} else {
			marker = lastMarker
		}
	}
	return
}
