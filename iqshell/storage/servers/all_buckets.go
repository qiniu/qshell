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
	Id                  *data.String `json:"id"`
	Tbl                 *data.String `json:"tbl"`
	CTime               *data.Int    `json:"ctime"` // 返回为 0
	FileNum             *data.Int64  `json:"file_num"`
	StorageSize         *data.Int64  `json:"storage_size"`
	Region              *data.String `json:"region"`
	Global              *data.Bool   `json:"global"`
	Perm                *data.Int    `json:"perm"`
	ShareUsers          []*UserInfo  `json:"share_users"`
	Versioning          *data.Bool   `json:"versioning"`
	AllowNullKey        *data.Bool   `json:"allow_nullkey"`
	EncryptionEnabled   *data.Bool   `json:"encryption_enabled"`
	NotAllowAccessByTbl *data.Bool   `json:"not_allow_access_by_tbl"`
}

func (i *BucketInfo) DescriptionString() string {
	return fmt.Sprintf("%s", i.Tbl.Value())
}

func (i *BucketInfo) DetailDescriptionString() string {
	sizeString := utils.FormatFileSize(i.StorageSize.Value())
	return fmt.Sprintf("%s\t%s\t%d\t%d(%s)", i.Tbl.Value(), i.Region.Value(), i.FileNum.Value(), i.StorageSize.Value(), sizeString)
}

type ListApiInfo struct {
	Shared bool
	Region string
}

// AllBuckets List list 所有 bucket
func AllBuckets(info ListApiInfo) (buckets []BucketInfo, err *data.CodeError) {
	return allBuckets(workspace.GetConfig(), info)
}

func allBuckets(cfg *config.Config, info ListApiInfo) ([]BucketInfo, *data.CodeError) {
	// https://github.com/qbox/product/blob/master/kodo/bucket/tblmgr.md#v3buckets%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7%E7%AC%A6%E5%90%88%E6%9D%A1%E4%BB%B6%E7%9A%84%E7%A9%BA%E9%97%B4%E4%BF%A1%E6%81%AF%E5%8C%85%E6%8B%AC%E7%A9%BA%E9%97%B4%E6%96%87%E4%BB%B6%E4%BF%A1%E6%81%AF
	bucketManager, err := bucket.GetBucketManager()
	if err != nil {
		return nil, err
	}

	ucHost := cfg.Hosts.GetOneUc()
	reqURL := fmt.Sprintf("%s/v3/buckets?shared=%v", utils.Endpoint(cfg.UseHttps.Value(), ucHost), info.Shared)
	if len(info.Region) > 0 {
		reqURL = fmt.Sprintf("%s&region=%s", reqURL, info.Region)
	}
	var buckets []BucketInfo
	e := bucketManager.Client.CredentialedCall(workspace.GetContext(), bucketManager.Mac, auth.TokenQiniu, &buckets, "POST", reqURL, nil)
	return buckets, data.ConvertError(e)
}
