package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
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

func TestMoveNoExistSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("move", test.BucketNotExist, test.Key, test.Bucket, "-k", "qshell_move.json", "-w")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMoveNoExistDestBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("move", test.Bucket, test.Key, test.BucketNotExist, "-k", "qshell_move.json", "-w")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMoveNoExistSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("move", test.Bucket, test.KeyNotExist, test.Bucket, "-k", "qshell_move.json", "-w")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestMoveNoSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("move")
	if !strings.Contains(errs, "SourceBucket can't empty") {
		t.Fail()
	}
}

func TestMoveNoSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("move", test.Bucket)
	if !strings.Contains(errs, "SourceKey can't empty") {
		t.Fail()
	}
}

func TestMoveNoDestBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("move", test.Bucket, test.KeyNotExist)
	if !strings.Contains(errs, "DestBucket can't empty") {
		t.Fail()
	}
}

func TestMoveDocument(t *testing.T) {
	test.TestDocument("move", t)
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

func TestBatchMoveDocument(t *testing.T) {
	test.TestDocument("batchmove", t)
}