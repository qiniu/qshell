//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/jpeg")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/png")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestMimeTypeNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.BucketNotExist, test.Key, "image/jpeg")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMimeTypeNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.KeyNotExist, "image/jpeg")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestMimeTypeNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestMimeTypeNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestMimeTypeNoMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.Key)
	if !strings.Contains(errs, "MimeType can't empty") {
		t.Fail()
	}
}

func TestMimeTypeDocument(t *testing.T) {
	test.TestDocument("chgm", t)
}

// 批量操作
func TestBatchChangeMimeType(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "image/jpeg" + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_chgm_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_chgm_fail.txt")

	path, err := test.CreateFileWithContent("batch_chgm.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchchgm", test.Bucket,
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

func TestBatchChangeMimeTypeRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "image/jpeg" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_chgm.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchchgm", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y")

	result, _ := test.RunCmdWithError("batchchgm", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if strings.Contains(result, "work redo") {
		t.Fatal("batch result: shouldn't redo because not set --record-redo-while-error")
	}

	result, _ = test.RunCmdWithError("batchchgm", test.Bucket, test.Bucket,
		"-i", path,
		"--enable-record",
		"--record-redo-while-error",
		"--worker", "4",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if !strings.Contains(result, "work redo") {
		t.Fatal("batch result: shouldn redo because set --record-redo-while-error")
	}
}

func TestBatchMimeTypeDocument(t *testing.T) {
	test.TestDocument("batchchgm", t)
}
