package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	result, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", "qshell_copy.json", "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	if !strings.Contains(result, "CDN refresh Code: 200, Info: success") {
		t.Fail()
	}
}
