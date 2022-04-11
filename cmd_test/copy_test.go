////go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	copyFile(t, test.Key, "qshell_copy.json")
}

func copyFile(t *testing.T, srcKey, destKey string) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, srcKey, test.Bucket, "-k", destKey, "-w")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestCopyNoExistSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.BucketNotExist, test.Key, test.Bucket, "-k", "qshell_copy.json", "-w")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestCopyNoExistDestBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.BucketNotExist, "-k", "qshell_copy.json", "-w")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestCopyNoExistSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.KeyNotExist, test.Bucket, "-k", "qshell_copy.json", "-w")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestCopyNoSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("copy")
	if !strings.Contains(errs, "SourceBucket can't empty") {
		t.Fail()
	}
}

func TestCopyNoSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket)
	if !strings.Contains(errs, "SourceKey can't empty") {
		t.Fail()
	}
}

func TestCopyNoDestBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.KeyNotExist)
	if !strings.Contains(errs, "DestBucket can't empty") {
		t.Fail()
	}
}

func TestCopyDocument(t *testing.T) {
	test.TestDocument("copy", t)
}

func TestBatchCopy(t *testing.T) {
	batchConfig := ""
	keys := test.OriginKeys
	for i, key := range keys {
		batchConfig += key + "\t" + test.Keys[i] + "\t" + "\n"
	}
	batchConfig += "\n"
	batchConfig += "hello10.json" + "\t" + "hello10_test.json" + "\t" + "\n"
	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_copy_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_copy_fail.txt")

	path, err := test.CreateFileWithContent("batch_copy.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchcopy", test.Bucket, test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4",
		"-w",
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

func TestBatchCopyDocument(t *testing.T) {
	test.TestDocument("batchcopy", t)
}
