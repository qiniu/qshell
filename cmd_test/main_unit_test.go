//go:build unit

package cmd

import (
	"os"
	"testing"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

// TestMain 在单元测试进程中启用 test mode。
//
// cmd.Execute 在非 test mode 下遇到非成功状态会直接 os.Exit，
// 会让同进程内后续用例无法执行；test mode 下改为返回，便于断言错误路径。
func TestMain(m *testing.M) {
	data.SetTestMode()
	os.Exit(m.Run())
}
