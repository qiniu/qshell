package list

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"strings"
)

func listBucketByV1(ctx context.Context, info ApiInfo, handler Handler) (hasMore bool, err *data.CodeError) {
	rets, hasMore, e := info.Manager.ListFilesWithContext(ctx, info.Bucket,
		storage.ListInputOptionsMarker(info.Marker),
		storage.ListInputOptionsPrefix(info.Prefix),
		storage.ListInputOptionsLimit(info.V1Limit),
		storage.ListInputOptionsDelimiter(info.Delimiter))
	if e == nil && rets == nil {
		return hasMore, data.NewError(0, "v1 meet empty body when list not completed")
	}

	if e != nil {
		return hasMore, data.ConvertError(e)
	}

	marker := rets.Marker
	commonPrefixes := rets.CommonPrefixes
	items := rets.Items
	dir := strings.Join(commonPrefixes, info.Delimiter)
	for _, item := range items {
		if handler(marker, dir, Item(item)) {
			break
		}
	}
	return hasMore, nil
}
