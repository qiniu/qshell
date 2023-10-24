//go:build integration

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

func TestMatch(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.Bucket, test.Key, objectPath)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestMatchNoExistBucket(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.BucketNotExist, test.Key, objectPath)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMatchNoExistSrcKey(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, objectPath)
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestMatchNoExistLocalFile(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, "/user/desktop/a.txt")
	if !strings.Contains(errs, "error:open /user/desktop/a.txt") {
		t.Fail()
	}
}

func TestMatchWithEmptyBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("match", "", test.Key, "/user/desktop/a.txt")
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestMatchNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, "", "/user/desktop/a.txt")
	if !strings.Contains(errs, "Key can't be empty") {
		t.Fail()
	}
}

func TestMatchNoLocalFile(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, "")
	if !strings.Contains(errs, "LocalFile can't be empty") {
		t.Fail()
	}
}

func TestMatchDocument(t *testing.T) {
	test.TestDocument("match", t)
}

func TestBatchMatch(t *testing.T) {
	TestBatchCopy(t)

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	objectPath := filepath.Join(resultDir, test.Key)
	_, _ = test.RunCmdWithError("get", test.Bucket, test.Key, "--domain", utils.Endpoint(false, test.BucketDomain),
		"-o", objectPath, "")
	defer test.RemoveFile(objectPath)

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_match.txt", test.KeysString)
	if err != nil {
		t.Fatal("create batch match config file error:", err)
	}

	test.RunCmdWithError("batchmatch", test.Bucket, resultDir,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4")
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

func TestBatchMatchWithRecord(t *testing.T) {
	TestBatchCopy(t)

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	objectPath := filepath.Join(resultDir, test.Key)
	_, _ = test.RunCmdWithError("get", test.Bucket, test.Key, "--domain", test.BucketDomain,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	keys := test.KeysString + "\nhello10_test.json"
	path, err := test.CreateFileWithContent("batch_match.txt", keys)
	if err != nil {
		t.Fatal("create batch match config file error:", err)
	}

	test.RunCmdWithError("batchmatch", test.Bucket, resultDir,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"-d")

	result, _ := test.RunCmdWithError("batchmatch", test.Bucket, resultDir,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"-d")
	if !strings.Contains(result, "because have done and") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: should skip the work had done")
	}
	if strings.Contains(result, "work redo") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: shouldn't redo because not set --record-redo-while-error")
	}

	result, _ = test.RunCmdWithError("batchmatch", test.Bucket, resultDir,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"--record-redo-while-error",
		"-d")
	if !strings.Contains(result, "because have done and") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: should skip the work had done")
	}
	if !strings.Contains(result, "work redo") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: shouldn redo because set --record-redo-while-error")
	}
}

func TestBatchMatchDocument(t *testing.T) {
	test.TestDocument("batchmatch", t)
}
