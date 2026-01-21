//go:build integration

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestRename(t *testing.T) {
	TestBatchCopy(t)

	key := "qshell_rename.json"
	_, errs := test.RunCmdWithError("rename", test.Bucket, test.Key, key, "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	// back
	_, errs = test.RunCmdWithError("rename", test.Bucket, key, test.Key, "-w")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestRenameNoExistSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.BucketNotExist, test.Key, "qshell_rename.json", "-w")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestRenameNoExistSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.Bucket, test.KeyNotExist, "qshell_rename.json", "-w")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestRenameNoSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("rename")
	if !strings.Contains(errs, "SourceBucket can't be empty") {
		t.Fail()
	}
}

func TestRenameNoSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.Bucket)
	if !strings.Contains(errs, "SourceKey can't be empty") {
		t.Fail()
	}
}

func TestRenameNoDestKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.Bucket, test.Key)
	if !strings.Contains(errs, "DestKey can't be empty") {
		t.Fail()
	}
}

func TestRenameDocument(t *testing.T) {
	test.TestDocument("rename", t)
}

func TestBatchRename(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "rename_" + key + "\t" + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_rename.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch rename config file error:", err)
	}

	test.RunCmdWithError("batchrename", test.Bucket,
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

func TestBatchRenameWithRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "move_" + key + "\t" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_rename.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchrename", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-w")

	result, _ := test.RunCmdWithError("batchrename", test.Bucket, test.Bucket,
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

	result, _ = test.RunCmdWithError("batchrename", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--record-redo-while-error",
		"--worker", "4",
		"--min-worker", "10",
		"--worker-count-increase-period", "50",
		"-y",
		"-w",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: should skip success work")
	}
	if !strings.Contains(result, "work redo") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: shouldn redo because set --record-redo-while-error")
	}
}

func TestBatchRenameDocument(t *testing.T) {
	test.TestDocument("batchrename", t)
}
