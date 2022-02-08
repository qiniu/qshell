package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestMove(t *testing.T) {
	key := "qshell_move.json"
	_, errs := test.RunCmdWithError("move", test.Bucket, test.Key, test.Bucket, "-k", key, "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	// back
	_, errs = test.RunCmdWithError("move", test.Bucket, key, test.Bucket, "-k", test.Key, "-w")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchMove(t *testing.T) {
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "move_" + key + "\t" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_move.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchmove", test.Bucket, test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}

	//back
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += "move_" + key + "\t" + key + "\t" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_move_back.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file after batch move error:", err)
	}

	_, errs = test.RunCmdWithError("batchmove", test.Bucket, test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
