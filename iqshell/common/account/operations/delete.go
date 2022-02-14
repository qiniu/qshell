package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

func Clean() {
	err := account.CleanUser()
	if err != nil {
		log.Error(err)
		os.Exit(data.StatusError)
	}
}

type RemoveInfo struct {
	Name string
}

func (info *RemoveInfo) Check() error {
	if len(info.Name) == 0 {
		return alert.CannotEmptyError("user name", "")
	}
	return nil
}

func Remove(info RemoveInfo) {
	err := account.RmUser(info.Name)
	if err != nil {
		log.Error(err)
		os.Exit(data.StatusError)
	}
}
