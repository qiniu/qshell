package operations

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type CleanInfo struct {
}

func (info *CleanInfo) Check() *data.CodeError {
	return nil
}

func Clean(cfg *iqshell.Config, info CleanInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	err := account.CleanUser()
	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
	}
}

type RemoveInfo struct {
	Name string
}

func (info *RemoveInfo) Check() *data.CodeError {
	if len(info.Name) == 0 {
		return alert.CannotEmptyError("user name", "")
	}
	return nil
}

func Remove(cfg *iqshell.Config, info RemoveInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	err := account.RmUser(info.Name)
	if err != nil {
		log.Error(err)
		data.SetCmdStatusError()
		return
	}
}
