//go:build integration

package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestDelete(t *testing.T) {
	deleteKey := "qshell_delete.json"
	copyFile(t, test.Key, deleteKey)
	deleteFile(t, deleteKey)
}

func deleteFile(t *testing.T, deleteKey string) {
	_, errs := test.RunCmdWithError("delete", test.Bucket, deleteKey)
	if len(errs) > 0 && !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestDeleteNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("delete", test.BucketNotExist, test.Key)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestDeleteNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("delete", test.Bucket, test.KeyNotExist)
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestDeleteNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("delete")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestDeleteNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("delete", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestDeleteDocument(t *testing.T) {
	test.TestDocument("delete", t)
}

func TestBatchDelete(t *testing.T) {
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

	successLogPath := filepath.Join(resultDir, "batch_delete_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_delete_fail.txt")

	path, err := test.CreateFileWithContent("batch_delete.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchdelete", test.Bucket,
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

func TestBatchDeleteRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_delete.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchdelete", test.Bucket,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"-y")

	result, _ := test.RunCmdWithError("batchdelete", test.Bucket,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if strings.Contains(result, "work redo") {
		t.Fatal("batch result: shouldn't redo because not set --record-redo-while-error")
	}

	result, _ = test.RunCmdWithError("batchdelete", test.Bucket,
		"-i", path,
		"--worker", "4",
		"--enable-record",
		"--record-redo-while-error",
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

func TestBatchDeleteDocument(t *testing.T) {
	test.TestDocument("batchdelete", t)
}

func TestDeleteAfter(t *testing.T) {
	TestBatchCopy(t)

	deleteKey := "qshell_delete_after.json"
	_, errs := test.RunCmdWithError("copy", test.Bucket, test.Key, test.Bucket, "-k", deleteKey, "-w")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("expire", test.Bucket, deleteKey, "1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchDeleteAfter(t *testing.T) {
	TestBatchCopy(t)

	// copy
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "delete_after_" + key + "\n"
	}

	path, err := test.CreateFileWithContent("batch_delete_after_copy.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchcopy", test.Bucket, test.Bucket, "-i", path, "-w", "-y")
	if len(errs) > 0 {
		t.Fail()
	}

	// delete
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += "delete_after_" + key + "\t" + "1" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_delete_after.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch expire after config file error:", err)
	}

	_, errs = test.RunCmdWithError("batchexpire", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}
