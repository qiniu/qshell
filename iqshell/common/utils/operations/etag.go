package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type EtagInfo struct {
	FilePath string
}

func (info *EtagInfo) Check() error {
	if len(info.FilePath) == 0 {
		return alert.CannotEmptyError("LocalFilePath", "")
	}
	return nil
}

// CreateEtag 计算文件的hash值，使用七牛的etag算法
func CreateEtag(cfg *iqshell.Config, info EtagInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if len(info.FilePath) == 0 {
		log.Error(alert.CannotEmpty("file path", ""))
		return
	}

	etag, err := utils.GetEtag(info.FilePath)
	if err != nil {
		log.Error(err)
		return
	}
	log.Alert(etag)
}
