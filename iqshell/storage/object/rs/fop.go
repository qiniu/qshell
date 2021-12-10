package rs

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
)

func PreFopStatus(persistentId string) (storage.PrefopRet, error) {
	if len(persistentId) == 0 {
		return storage.PrefopRet{}, errors.New(alert.CannotEmpty("persistent id", ""))
	}

	mac, err := account.GetMac()
	if err != nil {
		return storage.PrefopRet{}, err
	}

	opManager := storage.NewOperationManager(mac, nil)
	ret, err := opManager.Prefop(persistentId)
	return ret, err
}

type PreFopApiInfo struct {
	Bucket      string
	Key         string
	Fops        string
	Pipeline    string
	NotifyURL   string
	NotifyForce bool
}

func PreFop(info PreFopApiInfo) (persistentId string, err error) {
	if len(info.Bucket) == 0 {
		return "", errors.New(alert.CannotEmpty("bucket", ""))
	}

	if len(info.Key) == 0 {
		return "", errors.New(alert.CannotEmpty("key", ""))
	}

	if len(info.Fops) == 0 {
		return "", errors.New(alert.CannotEmpty("fops", ""))
	}

	mac, err := account.GetMac()
	if err != nil {
		return
	}
	opManager := storage.NewOperationManager(mac, nil)
	persistentId, err = opManager.Pfop(info.Bucket, info.Key, info.Fops, info.Pipeline, info.NotifyURL, info.NotifyForce)
	return
}
