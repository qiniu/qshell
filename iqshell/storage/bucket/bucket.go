package bucket

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func GetBucketManager() (manager *storage.BucketManager, err *data.CodeError) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("GetBucketManager: get current account error:%v", gErr)
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cfg := workspace.GetStorageConfig()
	manager = storage.NewBucketManager(mac, cfg)
	return
}

type GetBucketApiInfo struct {
	Bucket string
}

type BucketInfo struct {
	Id                     string   `json:"id"`
	Bucket                 string   `json:"tbl"`
	Region                 string   `json:"region"`
	Global                 bool     `json:"global"`
	Line                   bool     `json:"line"`
	CreateTime             int64    `json:"ctime"`
	Versioning             bool     `json:"versioning"`
	Private                int      `json:"private"`
	Product                string   `json:"product"`
	SysTags                []string `json:"systags"`
	MultiRegionEnabled     bool     `json:"multiregion_enabled"`
	MultiRegionEverEnabled bool     `json:"multiregion_ever_enabled"`
}

func GetBucketInfo(info GetBucketApiInfo) (bucketInfo *BucketInfo, err *data.CodeError) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	cfg := workspace.GetConfig()
	ucHost := cfg.Hosts.GetOneUc()
	if len(ucHost) == 0 {
		return nil, data.NewEmptyError().AppendDesc("can't get uc host")
	}

	bucketInfo = &BucketInfo{}
	url := utils.Endpoint(cfg.IsUseHttps(), ucHost)
	reqURL := fmt.Sprintf("%s/bucket/%s", url, info.Bucket)
	e := bucketManager.Client.CredentialedCall(context.Background(), bucketManager.Mac, auth.TokenQiniu, bucketInfo, "GET", reqURL, nil)
	return bucketInfo, data.ConvertError(e)
}
