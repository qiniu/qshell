package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
)

type PreFopStatusInfo struct {
	Id     string
	Bucket string // 用于查询 region，私有云必须，公有云可选
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

	ret, err := object.PreFopStatus(object.PreFopStatusApiInfo{
		Id:     info.Id,
		Bucket: info.Bucket,
	})
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("prefop status error:%v", err)
		return
	}

	log.Alert(ret.String())
}

type PfopInfo object.PfopApiInfo

func (info *PfopInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}

	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}

	if len(info.Fops) == 0 && len(info.WorkflowTemplateID) == 0 {
		return alert.CannotEmptyError("Fops", "")
	}
	return nil
}

func Pfop(cfg *iqshell.Config, info PfopInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	persistentId, err := object.Pfop(object.PfopApiInfo(info))
	if err != nil {
		data.SetCmdStatusError()
		log.ErrorF("pfop error:%v", err)
		return
	}
	log.Alert(persistentId)
}
