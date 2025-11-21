package servers

import (
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type UserInfo struct {
	UserId *data.Int `json:"uid"`
	Perm   *data.Int `json:"perm"`
}

type BucketQuota struct {
}

type BucketInfo struct {
	Name        *data.String `json:"name"`
	Tbl         *data.String `json:"tbl"`
	FileNum     *data.Int64  `json:"file_num"`
	StorageSize *data.Int64  `json:"storage_size"`
	Region      *data.String `json:"region"`

	//CTime       *data.String `json:"ctime"`
	//Global              *data.Bool   `json:"global"`
	//Perm                *data.Int    `json:"perm"`
	//ShareUsers          []*UserInfo  `json:"share_users"`
	//Versioning          *data.Bool   `json:"versioning"`
	//AllowNullKey        *data.Bool   `json:"allow_nullkey"`
	//EncryptionEnabled   *data.Bool   `json:"encryption_enabled"`
	//NotAllowAccessByTbl *data.Bool   `json:"not_allow_access_by_tbl"`
}

func (i *BucketInfo) BucketName() string {
	if i.Name != nil {
		return i.Name.Value()
	}

	if i.Tbl != nil {
		return i.Tbl.Value()
	}

	return ""
}
func (i *BucketInfo) DescriptionString() string {
	return fmt.Sprintf("%s", i.BucketName())
}

func (i *BucketInfo) DetailDescriptionString() string {
	sizeString := utils.FormatFileSize(i.StorageSize.Value())
	return fmt.Sprintf("%-20s\t%-10d\t%-10s\t%s", i.Region.Value(), i.FileNum.Value(), sizeString, i.BucketName())
}

func BucketInfoDetailDescriptionStringFormat() string {
	return fmt.Sprintf("%-20s\t%-10s\t%-10s\t%s", "Region", "FileNum", "StorageSize", "Bucket")
}

type BucketsResponse struct {
	NextMarker  string       `json:"next_marker"`
	IsTruncated bool         `json:"is_truncated"`
	Buckets     []BucketInfo `json:"buckets"`
}

type ListApiInfo struct {
	Region string
	Marker string
	Limit  int
	Detail bool
}

type BucketHandler func(bucket *BucketInfo, err *data.CodeError)

// AllBuckets List 所有 bucket
func AllBuckets(info ListApiInfo, handler BucketHandler) {
	// 分页获取没有更详细的信息，所以不能使用分页
	if info.Detail {
		allBuckets(workspace.GetConfig(), info, handler)
		return
	}

	// 对于只需要 bucket 名的情况，走分页获取
	allBucketsByPage(workspace.GetConfig(), info, handler)
}

// allBuckets 一次获取获取所有 Bucket
func allBuckets(cfg *config.Config, info ListApiInfo, handler BucketHandler) {
	// https://github.com/qbox/product/blob/master/kodo/bucket/tblmgr.md#v3buckets%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7%E7%AC%A6%E5%90%88%E6%9D%A1%E4%BB%B6%E7%9A%84%E7%A9%BA%E9%97%B4%E4%BF%A1%E6%81%AF%E5%8C%85%E6%8B%AC%E7%A9%BA%E9%97%B4%E6%96%87%E4%BB%B6%E4%BF%A1%E6%81%AF
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		handler(nil, err)
		return
	}

	ucHost := cfg.Hosts.GetOneUc()
	reqURL := fmt.Sprintf("%s/v3/buckets?shared=rd", utils.Endpoint(cfg.UseHttps.Value(), ucHost))
	if len(info.Region) > 0 {
		reqURL = fmt.Sprintf("%s&region=%s", reqURL, info.Region)
	}
	var buckets []BucketInfo
	rErr := bucketManager.Client.CredentialedCall(workspace.GetContext(), bucketManager.Mac, auth.TokenQiniu, &buckets, "POST", reqURL, nil)
	if rErr != nil {
		handler(nil, data.ConvertError(rErr))
		return
	}

	for _, b := range buckets {
		handler(&b, nil)
	}
}

// allBucketsByPage 分页获取 Bucket
func allBucketsByPage(cfg *config.Config, info ListApiInfo, handler BucketHandler) {
	if handler == nil {
		return
	}

	marker := info.Marker
	for {
		info.Marker = marker
		resp, err := allBucketsOnePage(cfg, info)
		if err != nil {
			handler(nil, err)
			break
		}

		for _, b := range resp.Buckets {
			handler(&b, nil)
		}

		if !resp.IsTruncated || len(resp.NextMarker) == 0 {
			break
		}

		marker = resp.NextMarker
	}
}

func allBucketsOnePage(cfg *config.Config, info ListApiInfo) (*BucketsResponse, *data.CodeError) {
	// 支持分页：https://github.com/qbox/product/blob/master/kodo/bucket/tblmgr.md#bucketsapiversionv4-%E6%94%AF%E6%8C%81%E5%88%86%E9%A1%B5%E8%BF%94%E5%9B%9E%E8%A1%A8%E4%BF%A1%E6%81%AF
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	if info.Limit <= 0 {
		info.Limit = 50
	}

	ucHost := cfg.Hosts.GetOneUc()
	reqURL := fmt.Sprintf("%s/buckets?apiVersion=v4&limit=%d", utils.Endpoint(cfg.UseHttps.Value(), ucHost), info.Limit)
	if len(info.Region) > 0 {
		reqURL = fmt.Sprintf("%s&region=%s", reqURL, info.Region)
	}
	if len(info.Marker) > 0 {
		reqURL = fmt.Sprintf("%s&marker=%s", reqURL, info.Marker)
	}

	var resp BucketsResponse
	e := bucketManager.Client.CredentialedCall(workspace.GetContext(), bucketManager.Mac, auth.TokenQiniu, &resp, "GET", reqURL, nil)
	if e != nil {
		return nil, data.ConvertError(e)
	}
	return &resp, nil
}
