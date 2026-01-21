package list

import (
	"context"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

	dir := strings.Join(rets.CommonPrefixes, info.Delimiter)
	for _, item := range rets.Items {
		if handler(rets.Marker, dir, Item(item)) {
			break
		}
	}
	return hasMore, nil
}
