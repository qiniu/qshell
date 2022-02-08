package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestAsyncFetch(t *testing.T) {
	path, err := test.CreateFileWithContent("async_fetch.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result, errs := test.RunCmdWithError("abfetch", test.Bucket,
		"-i", path,
		"-g", "1",
		"-c", "2")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	result, errs = test.RunCmdWithError("prefop", result)
	if len(errs) > 0 {
		t.Fail()
	}
}
