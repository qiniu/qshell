package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestPFop(t *testing.T) {
	result, errs := test.RunCmdWithError("pfop", "qshell-na0", "test_mv.mp4", "avthumb/mp4")
	if len(errs) > 0{
		t.Fail()
	}
	result = strings.ReplaceAll(result, "\n", "")
	result, errs = test.RunCmdWithError("prefop", result)
	if len(errs) > 0{
		t.Fail()
	}
}
