package rs

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"time"
)

type ListApiInfo struct {
	Bucket    string
	Prefix    string
	Marker    string
	Delimiter string
	MaxRetry  int // -1: 无限重试
}

type ListItem storage.ListItem

// List list 某个 bucket 所有的文件
func List(ctx context.Context, info *ListApiInfo) (<-chan ListItem, error) {
	objects := make(chan ListItem)

	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	go listObjectOfBucketToChan(ctx, bucketManager, info, objects)

	return objects, nil
}

func listObjectOfBucketToChan(ctx context.Context, manager *storage.BucketManager, info *ListApiInfo, objects chan<- ListItem) {
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

			objects <- ListItem(listItem.Item)
		}
		complete = true
	}

	close(objects)
}
