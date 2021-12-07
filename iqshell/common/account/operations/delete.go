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
		os.Exit(data.STATUS_ERROR)
	}
}

type RemoveInfo struct {
	Name string
}
func Remove(info RemoveInfo) {
	if len(info.Name) == 0 {
		log.Error(alert.CannotEmpty("user name", ""))
		return
	}

	err := account.RmUser(info.Name)
	if err != nil {
		log.Error(err)
		os.Exit(data.STATUS_ERROR)
	}
}
