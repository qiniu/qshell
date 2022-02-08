package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/jpeg")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/png")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchChangeMimeType(t *testing.T) {
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