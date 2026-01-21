package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

type listBucketV2Ret struct {
	Marker string `json:"marker"`
	Item   Item   `json:"item"`
	Dir    string `json:"dir"`
}

func listBucketByV2(ctx context.Context, info ApiInfo, handler Handler) (hasMore bool, err *data.CodeError) {
	ctx = auth.WithCredentialsType(ctx, info.Manager.Mac, auth.TokenQiniu)
	reqHost, reqErr := info.Manager.RsfReqHost(info.Bucket)
	if reqErr != nil {
		return false, data.ConvertError(reqErr)
	}

	reqURL := fmt.Sprintf("%s%s", reqHost, createListBucketV2Uri(info))
	resp, reqErr := info.Manager.Client.DoRequestWith(ctx, "POST", reqURL, nil, nil, 0)
	if reqErr != nil {
		return false, data.ConvertError(reqErr)
	}
	if resp.StatusCode/100 != 2 {
		return false, data.ConvertError(client.ResponseError(resp))
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var ret listBucketV2Ret
	for {
		dErr := dec.Decode(&ret)
		if dErr != nil {
			if dErr != io.EOF {
				return false, data.NewEmptyError().AppendDescF("decode error: %v", err)
			}
			break
		}

		if workspace.IsCmdInterrupt() {
			return false, data.NewError(0, "user cancel")
		}

		hasMore = len(ret.Marker) > 0
		if handler != nil {
			handler(ret.Marker, ret.Dir, ret.Item)
		}
	}

	return hasMore, nil
}

func createListBucketV2Uri(info ApiInfo) string {
	query := make(url.Values)
	query.Add("bucket", info.Bucket)
	if info.Prefix != "" {
		query.Add("prefix", info.Prefix)
	}
	if info.Delimiter != "" {
		query.Add("delimiter", info.Delimiter)
	}
	if info.Marker != "" {
		query.Add("marker", info.Marker)
	}
	return fmt.Sprintf("/v2/list?%s", query.Encode())
}
