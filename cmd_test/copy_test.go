package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestCopy(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", "qshell_copy.json", "-w")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchCopy(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", "qshell_copy.json", "-w")
	if len(errs) > 0 {
		t.Fail()
	}
}
