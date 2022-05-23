//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
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
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "move_" + key + "\t" + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_move.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchmove", test.Bucket, test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4",
		"-y",
		"-w")
	defer func() {
		test.RemoveFile(successLogPath)
		test.RemoveFile(failLogPath)
	}()

	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("batch result: success log to file error: file empty")
	}

	if !test.IsFileHasContent(failLogPath) {
		t.Fatal("batch result: fail log  to file error: file empty")
	}
}

func TestBatchMoveRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "move_" + key + "\t" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_move.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchmove", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-w")

	result, _ := test.RunCmdWithError("batchmove", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-w",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if strings.Contains(result, "work redo") {
		t.Fatal("batch result: shouldn't redo because not set --record-redo-while-error")
	}

	result, _ = test.RunCmdWithError("batchmove", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--record-redo-while-error",
		"--worker", "4",
		"-y",
		"-w",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if !strings.Contains(result, "work redo") {
		t.Fatal("batch result: shouldn redo because set --record-redo-while-error")
	}
}

func TestBatchMoveDocument(t *testing.T) {
	test.TestDocument("batchmove", t)
}
