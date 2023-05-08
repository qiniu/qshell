package list

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
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

func (l *Item) StorageTypeString() string {
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
