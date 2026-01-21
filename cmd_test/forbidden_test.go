//go:build integration

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestForbidden(t *testing.T) {
	_, errs := test.RunCmdWithError("forbidden", test.Bucket, test.Key)
	if len(errs) > 0 && !strings.Contains(errs, "already in normal stat") {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("forbidden", test.Bucket, test.Key, "-r")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestForbiddenNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("forbidden", test.BucketNotExist, test.Key, "0")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestForbiddenNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("forbidden", test.Bucket, test.KeyNotExist, "0")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestForbiddenNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("forbidden")
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestForbiddenNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("forbidden", test.Bucket)
	if !strings.Contains(errs, "Key can't be empty") {
		t.Fail()
	}
}

func TestForbiddenDocument(t *testing.T) {
	test.TestDocument("forbidden", t)
}

func TestBatchForbidden(t *testing.T) {
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

	successLogPath := filepath.Join(resultDir, "batch_forbidden_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_forbidden_fail.txt")

	path, err := test.CreateFileWithContent("batch_forbidden.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchforbidden", test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4",
		"-y")
	defer func() {
		// back
		test.RunCmdWithError("batchforbidden", test.Bucket, "-i", path, "-y", "-r")

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

func TestBatchForbiddenRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_forbidden.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchforbidden", test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-d",
		"-y")
	defer func() {
		// back
		test.RunCmdWithError("batchforbidden", test.Bucket, "-i", path, "-y", "-r")
	}()

	result, _ := test.RunCmdWithError("batchforbidden", test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: should skip success work")
	}
	if strings.Contains(result, "work redo") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: shouldn't redo because not set --record-redo-while-error")
	}

	result, _ = test.RunCmdWithError("batchforbidden", test.Bucket,
		"-i", path,
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

func TestBatchForbiddenDocument(t *testing.T) {
	test.TestDocument("batchforbidden", t)
}
