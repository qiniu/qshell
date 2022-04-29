package object

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func PreFopStatus(persistentId string) (storage.PrefopRet, *data.CodeError) {
	if len(persistentId) == 0 {
		return storage.PrefopRet{}, alert.CannotEmptyError("persistent id", "")
	}

	mac, err := account.GetMac()
	if err != nil {
		return storage.PrefopRet{}, err
	}

	cfg := workspace.GetStorageConfig()
	opManager := storage.NewOperationManager(mac, cfg)
	ret, e := opManager.Prefop(persistentId)
	return ret, data.ConvertError(e)
}

type PreFopApiInfo struct {
	Bucket      string
	Key         string
	Fops        string
	Pipeline    string
	NotifyURL   string
	NotifyForce bool
}

func PreFop(info PreFopApiInfo) (string, *data.CodeError) {
	if len(info.Bucket) == 0 {
		return "", alert.CannotEmptyError("bucket", "")
	}

	if len(info.Key) == 0 {
		return "", alert.CannotEmptyError("key", "")
	}

	if len(info.Fops) == 0 {
		return "", alert.CannotEmptyError("fops", "")
	}

	mac, err := account.GetMac()
	if err != nil {
		return "", err
	}
	cfg := workspace.GetStorageConfig()
	opManager := storage.NewOperationManager(mac, cfg)
	persistentId, e := opManager.Pfop(info.Bucket, info.Key, info.Fops, info.Pipeline, info.NotifyURL, info.NotifyForce)
	return persistentId, data.ConvertError(e)
}
