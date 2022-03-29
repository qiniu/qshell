package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

var (
	restoreKey = "restore_key.json"
)

func TestRestoreArchive(t *testing.T) {
	copyFile(t, test.OriginKeys[0], restoreKey)
	changeType(t, restoreKey, "2")
	_, errs := test.RunCmdWithError("restorear", test.Bucket, restoreKey, "1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestRestoreArchiveNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("restorear", test.BucketNotExist, restoreKey, "1")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestRestoreArchiveNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("restorear", test.Bucket, test.KeyNotExist, "1")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestRestoreArchiveNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("restorear")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestRestoreArchiveNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("restorear", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestRestoreArchiveNoFreezeAfterDays(t *testing.T) {
	_, errs := test.RunCmdWithError("restorear", test.Bucket, test.Key)
	if !strings.Contains(errs, "FreezeAfterDays can't empty") {
		t.Fail()
	}
}

func TestRestoreArchiveDocument(t *testing.T) {
	test.TestDocument("restorear", t)
}

func TestBatchRestoreArchive(t *testing.T) {
	copyFile(t, test.OriginKeys[0], restoreKey)
	changeType(t, restoreKey, "2")

	batchConfig := ""
	keys := []string{restoreKey}
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_restorear_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_restorear_fail.txt")

	path, err := test.CreateFileWithContent("batch_restorear.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchrestorear", test.Bucket, "1",
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4",
		"-y")
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

func TestBatchRestoreArchiveDocument(t *testing.T) {
	test.TestDocument("batchrestorear", t)
}
