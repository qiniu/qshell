package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

// ChangeInfo 切换账户
type ChangeInfo struct {
	Name string
}

func (info *ChangeInfo)Check() error {
	return nil
}

func Change(info ChangeInfo) {
	err := account.ChUser(info.Name)
	if err != nil {
		log.ErrorF("user change error:", err)
		os.Exit(data.StatusError)
		return
	}
}
