package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestBatchRename(t *testing.T) {
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "rename_" + key + "\t" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_rename.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch rename config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchrename", test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}

	//back
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += "rename_" + key + "\t" + key + "\t" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_rename_back.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch rename config file after batch rename error:", err)
	}

	_, errs = test.RunCmdWithError("batchrename", test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
