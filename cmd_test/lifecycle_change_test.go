//go:build integration

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestChangeLifecycle(t *testing.T) {
	TestBatchCopy(t)

	_, errs := test.RunCmdWithError("chlifecycle", test.Bucket, test.Key,
		"--to-ia-after-days", "30",
		"--to-archive-after-days", "60",
		"--to-deep-archive-after-days", "180",
		"--delete-after-days", "365")
	if len(errs) > 0 {
		t.Fail()
	}

	// back
	_, errs = test.RunCmdWithError("chlifecycle", test.Bucket, test.Key,
		"--to-ia-after-days", "-1",
		"--to-archive-after-days", "-1",
		"--to-deep-archive-after-days", "-1",
		"--delete-after-days", "-1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestChangeLifecycleNoExistSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chlifecycle", test.BucketNotExist, test.Key,
		"--to-ia-after-days", "-1",
		"--to-archive-after-days", "-1",
		"--to-deep-archive-after-days", "-1",
		"--delete-after-days", "-1")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestChangeLifecycleNoExistSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chlifecycle", test.Bucket, test.KeyNotExist,
		"--to-ia-after-days", "-1",
		"--to-archive-after-days", "-1",
		"--to-deep-archive-after-days", "-1",
		"--delete-after-days", "-1")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestChangeLifecycleNoSrcBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chlifecycle")
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestChangeLifecycleNoSrcKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chlifecycle", test.Bucket)
	if !strings.Contains(errs, "Key can't be empty") {
		t.Fail()
	}
}

func TestChangeLifecycleNoAction(t *testing.T) {
	_, errs := test.RunCmdWithError("chlifecycle", test.Bucket, test.KeyNotExist)
	if !strings.Contains(errs, "must set at least one value of lifecycle") {
		t.Fail()
	}
}

func TestChangeLifecycleDocument(t *testing.T) {
	test.TestDocument("chlifecycle", t)
}

func TestBatchChangeLifecycle(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_chlifecycle.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch change lifecycle config file error:", err)
	}

	test.RunCmdWithError("batchchlifecycle", test.Bucket,
		"-i", path,
		"--to-ia-after-days", "30",
		"--delete-after-days", "365",
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

func TestBatchChangeLifecycleWithRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_change_lifecycle.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch change lifecycle config file error:", err)
	}

	test.RunCmdWithError("batchchlifecycle", test.Bucket,
		"-i", path,
		"--to-ia-after-days", "30",
		"--to-archive-after-days", "60",
		"--to-deep-archive-after-days", "180",
		"--delete-after-days", "365",
		"--enable-record",
		"--worker", "4",
		"-y")

	result, _ := test.RunCmdWithError("batchchlifecycle", test.Bucket,
		"-i", path,
		"--to-ia-after-days", "30",
		"--to-archive-after-days", "60",
		"--to-deep-archive-after-days", "180",
		"--delete-after-days", "365",
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

	result, _ = test.RunCmdWithError("batchchlifecycle", test.Bucket, test.Bucket,
		"-i", path,
		"--to-ia-after-days", "30",
		"--to-archive-after-days", "60",
		"--to-deep-archive-after-days", "180",
		"--delete-after-days", "365",
		"--enable-record",
		"--record-redo-while-error",
		"--worker", "4",
		"--min-worker", "10",
		"--worker-count-increase-period", "50",
		"-y",
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

func TestBatchChangeLifecycleNoAction(t *testing.T) {
	batchConfig := ""
	keys := test.Keys
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_change_lifecycle.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch change lifecycle config file error:", err)
	}
	defer func() {
		_ = test.RemoveFile(path)
	}()

	_, errs := test.RunCmdWithError("batchchlifecycle", test.Bucket,
		"-i", path,
		"-y")
	if !strings.Contains(errs, "must set at least one value of lifecycle") {
		t.Fail()
	}
}

func TestBatchChangeLifecycleDocument(t *testing.T) {
	test.TestDocument("batchchlifecycle", t)
}
