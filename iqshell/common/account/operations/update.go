package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

// ChangeInfo 切换账户
type ChangeInfo struct {
	Name string
}

func (info *ChangeInfo) Check() *data.CodeError {
	return nil
}

func Change(cfg *iqshell.Config, info ChangeInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	name, err := account.ChUser(info.Name)
	if err != nil {
		log.ErrorF("user change to %s failed, error:%v", name, err)
		return
	} else {
		log.InfoF("user change to %s success", name)
	}
}
