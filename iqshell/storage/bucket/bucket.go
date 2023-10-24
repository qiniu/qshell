package bucket

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qiniu/qshell/v2/iqshell/common/client"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func GetBucketManager() (manager *storage.BucketManager, err *data.CodeError) {
	acc, gErr := workspace.GetAccount()
	if gErr != nil {
		err = data.NewEmptyError().AppendDescF("GetBucketManager: get current account error:%v", gErr)
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cfg := workspace.GetStorageConfig()
	c := client.DefaultStorageClient()
	manager = storage.NewBucketManagerEx(mac, cfg, &c)
	return
}

type GetBucketApiInfo struct {
	Bucket string
}

type BucketInfo storage.BucketInfo

func GetBucketInfo(info GetBucketApiInfo) (*BucketInfo, *data.CodeError) {
	bucketManager, err := GetBucketManager()
	if err != nil {
		return nil, err
	}

	if bucketInfo, gErr := bucketManager.GetBucketInfo(info.Bucket); gErr != nil {
		return nil, data.ConvertError(gErr)
	} else {
		return (*BucketInfo)(&bucketInfo), nil
	}
}
