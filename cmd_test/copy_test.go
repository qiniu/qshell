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
	batchCopyConfig := ""
	for _, key := range test.Keys {
		batchCopyConfig += key + "\t" + "copy_" + key + "\t" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_copy.txt", batchCopyConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchcopy", test.Bucket, test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
