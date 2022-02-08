package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestDelete(t *testing.T) {
	deleteKey := "qshell_delete.json"
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", deleteKey, "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("delete", test.Bucket, deleteKey)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchDelete(t *testing.T) {
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += "copy_" + key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_delete.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchdelete", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
