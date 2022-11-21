package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestQShellDocument01(t *testing.T) {
	prefix := "# 简介\n`qshell`"
	result, _ := test.RunCmdWithError()
	if !strings.HasPrefix(result, prefix) {
		t.Fatal("document test fail for cmd: qshell")
	}
}

func TestQShellDocument02(t *testing.T) {
	prefix := "# 简介\n`qshell`"
	result, _ := test.RunCmdWithError("--doc")
	if !strings.HasPrefix(result, prefix) {
		t.Fatal("document test fail for cmd: qshell")
	}
}
