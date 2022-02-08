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

func TestDeleteAfter(t *testing.T) {
	deleteKey := "qshell_delete_after.json"
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", deleteKey, "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("expire", test.Bucket, deleteKey, "1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchDeleteAfter(t *testing.T) {
	// copy
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "delete_after_" + key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_delete_after_copy.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchcopy", test.Bucket, test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}

	// delete
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += "delete_after_" + key + "\t" + "1" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_delete_after.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch expire after config file error:", err)
	}

	_, errs = test.RunCmdWithError("batchexpire", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
