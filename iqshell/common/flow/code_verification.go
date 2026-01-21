package flow

import (
	"fmt"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

func UserCodeVerification() (success bool) {
	code := utils.CreateRandString(6)
	log.Warning(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", code))

	confirm := ""
	_, err := fmt.Scanln(&confirm)
	if err != nil {
		_, _ = fmt.Fprintf(data.Stdout(), "scan error:%v\n", err)
		return false
	}

	if code != confirm {
		_, _ = fmt.Fprintln(data.Stdout(), "Task quit!")
		return false
	}

	return true
}
