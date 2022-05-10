package iqshell

import (
	"github.com/qiniu/go-sdk/v7/storage"
)

func Prefop(persistentId string) (ret storage.PrefopRet, err error) {
	mac, err := GetMac()
	if err != nil {
		return
	}
	opManager := storage.NewOperationManager(mac, nil)
	ret, err = opManager.Prefop(persistentId)
	return
}

func Pfop(bucket, key, fops, pipeline, notifyURL string, notifyForce bool) (persistentId string, err error) {

	mac, err := GetMac()
	if err != nil {
		return
	}
	opManager := storage.NewOperationManager(mac, nil)
	persistentId, err = opManager.Pfop(bucket, key, fops, pipeline, notifyURL, notifyForce)
	return
}
