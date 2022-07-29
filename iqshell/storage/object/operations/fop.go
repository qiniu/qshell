package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
)

type PreFopStatusInfo struct {
	Id string
}

func (info *PreFopStatusInfo) Check() *data.CodeError {
	if len(info.Id) == 0 {
		return alert.CannotEmptyError("PersistentID", "")
	}
	return nil
}

func PreFopStatus(cfg *iqshell.Config, info PreFopStatusInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	ret, err := object.PreFopStatus(info.Id)
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("pre fog status error:%v", err)
		return
	}

	log.Alert(ret.String())
}

type PreFopInfo object.PreFopApiInfo

func (info *PreFopInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}

	if len(info.Fops) == 0 {
		return alert.CannotEmptyError("Fops", "")
	}
	return nil
}

func PreFop(cfg *iqshell.Config, info PreFopInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	persistentId, err := object.PreFop(object.PreFopApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("pre fog error:%v", err)
		return
	}
	log.Alert(persistentId)
}
