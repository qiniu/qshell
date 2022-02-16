package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
)

type PreFopStatusInfo struct {
	Id string
}

func (info *PreFopStatusInfo) Check() error {
	if len(info.Id) == 0 {
		return alert.CannotEmptyError("PersistentID", "")
	}
	return nil
}

func PreFopStatus(info PreFopStatusInfo) {
	ret, err := object.PreFopStatus(info.Id)
	if err != nil {
		log.ErrorF("pre fog status error:%v", err)
		return
	}

	log.Alert(ret.String())
}

type PreFopInfo object.PreFopApiInfo

func (info *PreFopInfo) Check() error {
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

func PreFop(info PreFopInfo) {
	persistentId, err := object.PreFop(object.PreFopApiInfo(info))
	if err != nil {
		log.ErrorF("pre fog error:%v", err)
		return
	}
	log.Alert(persistentId)
}
