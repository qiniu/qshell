package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestChangeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, test.Key, "2")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chtype", test.Bucket, test.Key, "1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchChangeType(t *testing.T) {
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "2" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_chtype.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchchtype", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
