package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	TestBatchCopy(t)

	result, errs := test.RunCmdWithError("stat", test.Bucket, test.Key)
	if len(errs) > 0 {
		t.Fail()
	}

	if !strings.Contains(result, "FileHash") {
		t.Fail()
	}
}

func TestBatchStatus(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "status_" + key + "\t" + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath :=  filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_status.txt", batchConfig)
	if err != nil {
		t.Fatal("create config file error:", err)
	}

	test.RunCmdWithError("batchstat", test.Bucket,
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

func TestBatchStatucDocument(t *testing.T) {
	test.TestDocument("batchstat", t)
}
