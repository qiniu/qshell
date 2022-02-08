package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
)

// 切换账户
type ChangeInfo struct {
	Name string
}

func Change(info ChangeInfo) {
	err := account.ChUser(info.Name)
	if err != nil {
		log.ErrorF("user change error:", err)
		os.Exit(data.StatusError)
		return
	}
}
