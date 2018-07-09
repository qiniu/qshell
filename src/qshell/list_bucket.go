package qshell

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"net/url"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	"qiniu/api.v6/rsf"
	"qiniu/rpc"
	"strconv"
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
	defer bWriter.Flush()

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
			lineData := fmt.Sprintln(fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%d\t%s",
				entry.Key, entry.Fsize, entry.Hash, entry.PutTime, entry.MimeType, entry.FileType,
				entry.FileStatus, entry.EndUser))
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

type listBucketRetV2 struct {
	Item   rsf.ListItem `json:"item"`
	Marker string       `json:"marker"`
	Dir    string       `json:"dir"`
}

// ListFilesV2 改进版本的 ListFiles 以解决 ListFiles 的超时问题
// https://github.com/qbox/product/blob/master/kodo/rsf.md#v2list-%E5%88%97%E5%87%BA%E5%86%85%E5%AE%B9
// 在这个方法中，需要注意的是即使函数返回的 err 不为 nil，entries 也有可能有值，另外 nextMarker 也有可能不为空，所以
// 正确的逻辑是检查 hasNext 是否为 true，如果有表示还可以继续使用 nextMarker 来进行list，另外保存下 entries 的记录
// 在这些逻辑处理完毕之后，检查下是否 err 不为 nil，如果不为 nil，应该打印一个 Warnning 的日志表示 list 曾经出现过错误
// 另外如果你希望全量列举空间的话，limit 参数设置为 0 即可。
func ListBucketV2(mac *digest.Mac, bucket, prefix, marker, listResultFile string) (nextMarker string, err error) {
	var listResultFh *os.File
	if listResultFile == "stdout" {
		listResultFh = os.Stdout
	} else {
		var openErr error
		//if marker not empty, continue the list
		if marker != "" {
			listResultFh, openErr = os.OpenFile(listResultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if openErr != nil {
				err = openErr
				return
			}
		} else {
			listResultFh, openErr = os.Create(listResultFile)
			if openErr != nil {
				err = openErr
				return
			}
		}
	}
	defer listResultFh.Close()
	bWriter := bufio.NewWriter(listResultFh)
	defer bWriter.Flush()

	//init request
	reqURL := fmt.Sprintf("%s%s", conf.RSF_HOST, makeListURLV2(bucket, prefix, "", marker, 0))
	req, newErr := http.NewRequest("POST", reqURL, nil)
	if newErr != nil {
		err = newErr
		return
	}

	accessToken, signErr := mac.SignRequest(req, false)
	if signErr != nil {
		err = signErr
		return
	}
	req.Header.Add("Authorization", "QBox "+accessToken)
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = respErr
		return
	}

	bScanner := bufio.NewScanner(resp.Body)

	for bScanner.Scan() {
		eachLine := bScanner.Text()
		//try to parse into ListItem
		var listItem listBucketRetV2
		mErr := json.Unmarshal([]byte(eachLine), &listItem)
		if mErr != nil {
			//return
			err = mErr
			return
		}

		nextMarker = listItem.Marker
		entry := listItem.Item
		//write entries
		lineData := fmt.Sprintln(fmt.Sprintf("%s\t%d\t%s\t%d\t%s\t%d\t%d\t%s",
			entry.Key, entry.Fsize, entry.Hash, entry.PutTime, entry.MimeType, entry.FileType,
			entry.FileStatus, entry.EndUser))
		_, wErr := bWriter.WriteString(lineData)
		if wErr != nil {
			err = wErr
			return
		}
	}
	return
}

func makeListURLV2(bucket, prefix, delimiter, marker string, limit int) string {
	query := make(url.Values)
	query.Add("bucket", bucket)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if delimiter != "" {
		query.Add("delimiter", delimiter)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}
	return fmt.Sprintf("/v2/list?%s", query.Encode())
}
