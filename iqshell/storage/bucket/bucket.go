package bucket

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func GetBucketManager() (manager *storage.BucketManager, err error) {
	acc, gErr := account.GetAccount()
	if gErr != nil {
		err = errors.New("GetBucketManager: get current account error:" + gErr.Error())
		return
	}

	mac := qbox.NewMac(acc.AccessKey, acc.SecretKey)
	cfg := workspace.GetStorageConfig()
	manager = storage.NewBucketManager(mac, cfg)
	return
}
