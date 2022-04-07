//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestRename(t *testing.T) {
	TestBatchCopy(t)

	TestUserIntegration(t)
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
	if !strings.Contains(errs, "SourceBucket can't empty") {
		t.Fail()
	}
}

func TestRenameNoSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.Bucket)
	if !strings.Contains(errs, "SourceKey can't empty") {
		t.Fail()
	}
}

func TestRenameNoDestKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rename", test.Bucket, test.Key)
	if !strings.Contains(errs, "DestKey can't empty") {
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

func TestBatchRenameDocument(t *testing.T) {
	test.TestDocument("batchrename", t)
}
