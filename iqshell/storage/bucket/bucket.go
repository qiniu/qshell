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
	cfg := workspace.GetConfig()
	r := (&cfg).GetRegion()
	if len(cfg.Hosts.GetOneUc()) > 0 {
		storage.SetUcHost(cfg.Hosts.GetOneUc(), cfg.IsUseHttps())
	}
	manager = storage.NewBucketManager(mac, &storage.Config{
		UseHTTPS:      cfg.IsUseHttps(),
		Region:        r,
		Zone:          r,
		CentralRsHost: cfg.Hosts.GetOneRs(),
	})
	return
}
