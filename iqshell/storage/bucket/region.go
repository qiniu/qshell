package bucket

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
)

func Region(cfg *storage.Config, bucket string) (*storage.Region, error) {
	if cfg.Region != nil {
		return cfg.Region, nil
	}

	acc, err := workspace.GetAccount()
	if err != nil {
		return nil, errors.New("get account error:" + err.Error())
	}

	if len(acc.AccessKey) == 0 {
		return nil, errors.New("can't get access key")
	}

	r, err := storage.GetRegion(acc.AccessKey, bucket)
	if err != nil {
		return nil, errors.New("get region error:" + err.Error())
	}

	return r, nil
}
