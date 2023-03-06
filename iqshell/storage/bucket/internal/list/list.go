package list

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strings"
)

type ApiVersion string

const (
	ApiVersionV1 = "v1"
	ApiVersionV2 = "v2"
)

type ApiInfo struct {
	Manager    *storage.BucketManager
	ApiVersion ApiVersion
	Bucket     string
	Prefix     string
	Delimiter  string
	Marker     string
	V1Limit    int
}

type Item storage.ListItem

func (l *Item) IsNull() bool {
	if l == nil {
		return true
	}

	return len(l.Key) == 0 && len(l.Hash) == 0 &&
		len(l.MimeType) == 0 && len(l.EndUser) == 0 &&
		l.PutTime == 0 && l.Type == 0 && l.Fsize == 0
}

func (l *Item) PutTimeString() string {
	if l.PutTime < 1 {
		return ""
	}
	return fmt.Sprintf("%d", l.PutTime)
}

func (l *Item) FileSizeString() string {
	if l.Fsize < 1 {
		return ""
	}
	return fmt.Sprintf("%d", l.Fsize)
}

func (l *Item) FileTypeString() string {
	if l.Type < 1 {
		return ""
	}
	return fmt.Sprintf("%d", l.Type)
}

type Handler func(marker string, dir string, item Item) (stop bool)

func ListBucket(ctx context.Context, info ApiInfo, handler Handler) (hasMore bool, err *data.CodeError) {
	if info.Manager == nil {
		return true, alert.CannotEmptyError("bucket manager", "")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if ctx.Err() != nil {
		return false, data.ConvertError(ctx.Err())
	}

	if info.ApiVersion == ApiVersionV1 {
		log.DebugF("list by api v1, marker:%s", info.Marker)
		return listBucketByV1(ctx, info, handler)
	} else {
		log.DebugF("list by api v2, marker:%s", info.Marker)
		return listBucketByV2(ctx, info, handler)
	}
}

func listBucketByV1(ctx context.Context, info ApiInfo, handler Handler) (hasMore bool, err *data.CodeError) {
	rets, commonPrefixes, marker, hasMore, e := info.Manager.ListFiles(info.Bucket, info.Prefix, info.Delimiter, info.Marker, info.V1Limit)
	if e == nil && rets == nil {
		return hasMore, data.NewError(0, "v1 meet empty body when list not completed")
	}

	if e != nil {
		return hasMore, data.ConvertError(e)
	}

	dir := strings.Join(commonPrefixes, info.Delimiter)
	for _, item := range rets {
		if handler(marker, dir, Item(item)) {
			break
		}
	}
	return hasMore, nil
}

func listBucketByV2(ctx context.Context, info ApiInfo, handler Handler) (hasMore bool, err *data.CodeError) {
	ret, e := info.Manager.ListBucketContext(ctx, info.Bucket, info.Prefix, info.Delimiter, info.Marker)
	if e == nil && ret == nil {
		return true, data.NewError(0, "v2 meet empty body when list not completed")
	}

	if e != nil {
		return true, data.ConvertError(e)
	}

	marker := ""
	for item := range ret {
		marker = item.Marker
		if handler(item.Marker, item.Dir, Item(item.Item)) {
			break
		}
	}
	return len(marker) > 0, nil
}
