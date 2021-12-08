package utils

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"runtime"
)

// 生成客户端代理名称
func UserAgent() string {
	return fmt.Sprintf("QShell/%s (%s; %s; %s)", data.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
