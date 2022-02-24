package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	TestBatchCopy(t)

	result, errs := test.RunCmdWithError("stat", test.Bucket, test.Key)
	if len(errs) > 0 {
		t.Fail()
	}

	if !strings.Contains(result, "FileHash") {
		t.Fail()
	}
}

func TestBatchStatus(t *testing.T) {
	TestBatchCopy(t)

	path, err := test.CreateFileWithContent("batch_delete.txt", test.KeysString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchstat", test.Bucket, "-i", path)
	if len(errs) > 0 {
		t.Fail()
	}
}
