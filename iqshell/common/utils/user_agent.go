package utils

import (
	"fmt"
	"runtime"

	"github.com/qiniu/qshell/v2/iqshell/common/version"
)

// 生成客户端代理名称
func UserAgent() string {
	return fmt.Sprintf("QiniuQShell/%s (%s; %s; %s)", version.Version(), runtime.GOOS, runtime.GOARCH, runtime.Version())
}
