package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestFromUpload(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result, errs := test.RunCmdWithError("rput",test.Bucket, "qshell_rput_5M", path)
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}
