package bucket

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"os"
	"time"
)

func GetBucketManager() (manager *storage.BucketManager, err error) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = errors.New("GetBucketManager: get current account error:" + gErr.Error())
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cfg := workspace.GetConfig()
	r := (&cfg).GetRegion()
	if len(cfg.Hosts.GetOneUc()) > 0 {
		storage.SetUcHost(cfg.Hosts.GetOneUc(), cfg.IsUseHttps())
	}
	manager = storage.NewBucketManager(mac, &storage.Config{
		UseHTTPS:      cfg.IsUseHttps(),
		Region:        r,
		Zone:          r,
		CentralRsHost: cfg.Hosts.GetOneRs(),
	})
	return
}

func CheckExists(bucket, key string) (exists bool, err error) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return false, err
	}

	entry, sErr := bucketManager.Stat(bucket, key)
	if sErr != nil {
		if v, ok := sErr.(*storage.ErrorInfo); !ok {
			err = fmt.Errorf("Check file exists error, %s", sErr.Error())
			return
		} else {
			if v.Code != 612 {
				err = fmt.Errorf("Check file exists error, %s", v.Err)
				return
			} else {
				exists = false
				return
			}
		}
	}
	if entry.Hash != "" {
		exists = true
	}
	return
}

func ListBucketToFile(bucket, prefix, marker, listResultFile, delimiter string, startDate, endDate time.Time, suffixes []string, maxRetry int, appendMode bool, readable bool) (retErr error) {
	lastMarker := marker

	defer func(lastMarker string) {
		if lastMarker != "" {
			fmt.Fprintf(os.Stderr, "Marker: %s\n", lastMarker)
		}
	}(lastMarker)

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
			log.ErrorF("Failed to open list result file `%s`", listResultFile)
			return
		}
		defer listResultFh.Close()
	}

	bWriter := bufio.NewWriter(listResultFh)

	notfilterTime := startDate.IsZero() && endDate.IsZero()
	notfilterSuffix := len(suffixes) == 0

	bucketManager, err := GetBucketManager()
	if err != nil {
		return err
	}

	var c int
	for {
		if maxRetry >= 0 && c >= maxRetry {
			break
		}

		entries, lErr := bucketManager.ListBucketContext(workspace.GetContext(), bucket, prefix, delimiter, marker)
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
					putTimeValid = filterByPutTime(putTime, startDate, endDate)
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

func errorWarning(marker string, err error) {
	log.ErrorF("marker: %s", marker)
	log.ErrorF("listbucket Error: %v", err)
}
