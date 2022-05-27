package utils

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/version"
	"runtime"
)

// 生成客户端代理名称
func UserAgent() string {
	return fmt.Sprintf("QShell/%s (%s; %s; %s)", version.Version(), runtime.GOOS, runtime.GOARCH, runtime.Version())
}
