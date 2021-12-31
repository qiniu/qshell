package bucket

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"strings"
	"time"
)

type ListApiInfo struct {
	Bucket    string
	Prefix    string
	Marker    string
	Delimiter string
	StartTime time.Time // list item 的 put time 区间的开始时间 【闭区间】
	EndTime   time.Time // list item 的 put time 区间的终止时间 【闭区间】
	Suffixes  []string  // list item 必须包含前缀
	MaxRetry  int       // -1: 无限重试
}

type ListItem storage.ListItem

// List list 某个 bucket 所有的文件
func List(info *ListApiInfo) (<-chan ListItem, error) {
	objects := make(chan ListItem)

	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	go listBucketToChan(workspace.GetContext(), bucketManager, info, objects)
	return objects, nil
}

func listBucketToChan(ctx context.Context, manager *storage.BucketManager, info *ListApiInfo, objects chan<- ListItem) {

	shouldCheckPutTime := !info.StartTime.IsZero() || !info.StartTime.IsZero()
	shouldCheckSuffixes := len(info.Suffixes) > 0
	complete := false
	for retryCount := 0; !complete && (info.MaxRetry < 0 || retryCount <= info.MaxRetry); retryCount++ {
		entries, err := manager.ListBucketContext(ctx, info.Bucket, info.Prefix, info.Delimiter, info.Marker)
		if entries == nil && err == nil {
			// no data
			if info.Marker == "" {
				complete = true
				break
			} else {
				log.Error("meet empty body when list not completed")
				continue
			}
		}

		if err != nil {
			log.ErrorF("marker: %s", info.Marker)
			log.ErrorF("listbucket Error: %v", err)
			time.Sleep(1)
			continue
		}

		for listItem := range entries {
			if listItem.Marker != info.Marker {
				info.Marker = listItem.Marker
			}

			if listItem.Item.IsEmpty() {
				continue
			}

			if shouldCheckPutTime {
				putTime := time.Unix(listItem.Item.PutTime/1e7, 0)
				if !filterByPutTime(putTime, info.StartTime, info.EndTime) {
					continue
				}
			}

			if shouldCheckSuffixes && !filterBySuffixes(listItem.Item.Key, info.Suffixes) {
				continue
			}

			objects <- ListItem(listItem.Item)
		}
		complete = true
	}

	close(objects)
}

func filterByPutTime(putTime, startDate, endDate time.Time) bool {
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