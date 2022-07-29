package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage"
)

type MirrorUpdateInfo storage.PrefetchApiInfo

func (info *MirrorUpdateInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func MirrorUpdate(cfg *iqshell.Config, info MirrorUpdateInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	err := storage.Prefetch(storage.PrefetchApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Mirror update Failed, [%s:%s], Error: %v", info.Bucket, info.Key, err)
	} else {
		log.InfoF("Mirror update Success, [%s:%s]", info.Bucket, info.Key)
	}
}
