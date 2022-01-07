package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestPFop(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop", "qiniutest", "test.avi", "avthumb/mp4")
	if !strings.Contains(errs, "400014") {
		t.Fail()
	}
}
