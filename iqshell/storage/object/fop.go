package object

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

type PreFopStatusApiInfo struct {
	Id     string
	Bucket string // 用于查询 region，私有云必须，公有云可选
}

func PreFopStatus(info PreFopStatusApiInfo) (storage.PrefopRet, *data.CodeError) {
	if len(info.Id) == 0 {
		return storage.PrefopRet{}, alert.CannotEmptyError("persistent id", "")
	}

	opManager, err := getOperationManager(info.Bucket)
	if err != nil {
		return storage.PrefopRet{}, err
	}
	ret, e := opManager.Prefop(info.Id)
	return ret, data.ConvertError(e)
}

type PreFopApiInfo struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	Fops        string `json:"fops"`
	Pipeline    string `json:"pipeline"`
	NotifyURL   string `json:"notify_url"`
	NotifyForce bool   `json:"notify_force"`
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

	opManager, err := getOperationManager(info.Bucket)
	if err != nil {
		return "", err
	}
	persistentId, e := opManager.Pfop(info.Bucket, info.Key, info.Fops, info.Pipeline, info.NotifyURL, info.NotifyForce)
	return persistentId, data.ConvertError(e)
}

func getOperationManager(bucket string) (*storage.OperationManager, *data.CodeError) {
	mac, err := account.GetMac()
	if err != nil {
		return nil, err
	}
	cfg := workspace.GetStorageConfig()
	opManager := storage.NewOperationManager(mac, cfg)
	if len(bucket) == 0 && (opManager.Cfg.Region == nil || opManager.Cfg.Zone == nil || len(opManager.Cfg.Region.ApiHost) != 0) {
		if region, e := storage.GetZone(opManager.Mac.AccessKey, bucket); e != nil {
			return nil, data.ConvertError(e)
		} else {
			opManager.Cfg.Region = region
			opManager.Cfg.Zone = region
			opManager.Cfg.CentralRsHost = region.RsHost
		}
	}
	return opManager, nil
}
