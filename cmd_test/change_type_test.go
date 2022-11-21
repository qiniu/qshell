//go:build integration

package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestChangeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, test.Key, "0")
	if len(errs) > 0 && !strings.Contains(errs, "already in normal stat") {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chtype", test.Bucket, test.Key, "1")
	if len(errs) > 0 {
		t.Fail()
	}
}

func changeType(t *testing.T, key string, ty string) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, key, ty)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestChangeTypeNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.BucketNotExist, test.Key, "0")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestChangeTypeNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, test.KeyNotExist, "0")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestChangeTypeNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestChangeTypeNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestChangeTypeNoMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, test.Key)
	if !strings.Contains(errs, "Type can't empty") {
		t.Fail()
	}
}

func TestChangeTypeDocument(t *testing.T) {
	test.TestDocument("chtype", t)
}

func TestBatchChangeType(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "0" + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_chtype_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_chtype_fail.txt")

	path, err := test.CreateFileWithContent("batch_chtype.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("batchchtype", test.Bucket,
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

	//back
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "1" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_chtype.txt", batchConfig)
	if err != nil {
		t.Fatal("create chtype config file error:", err)
	}

	test.RunCmdWithError("batchchtype", test.Bucket, "-i", path, "-y")
}

func TestBatchChangeTypeRecord(t *testing.T) {
	TestBatchCopy(t)

	batchConfig := ""
	keys := test.Keys
	keys = append(keys, "hello10.json")
	for _, key := range keys {
		batchConfig += key + "\t" + "0" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_chtype.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch move config file error:", err)
	}

	test.RunCmdWithError("batchchtype", test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-d",
		"-y")

	result, _ := test.RunCmdWithError("batchchtype", test.Bucket,
		"-i", path,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and") {
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

	result, _ = test.RunCmdWithError("batchchtype", test.Bucket,
		"-i", path,
		"--enable-record",
		"--record-redo-while-error",
		"--worker", "4",
		"--min-worker", "10",
		"--worker-count-increase-period", "50",
		"-y",
		"-d")
	if !strings.Contains(result, "because have done and success") {
		t.Fatal("batch result: should skip success work")
	}
	if !strings.Contains(result, "work redo") {
		fmt.Println("=========================== result start ===========================")
		fmt.Println(result)
		fmt.Println("=========================== result   end ===========================")
		t.Fatal("batch result: should redo because set --record-redo-while-error")
	}

	//back
	batchConfig = ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "1" + "\n"
	}

	path, err = test.CreateFileWithContent("batch_chtype.txt", batchConfig)
	if err != nil {
		t.Fatal("create chtype config file error:", err)
	}

	test.RunCmdWithError("batchchtype", test.Bucket, "-i", path, "-y")
}

func TestBatchChangeTypeDocument(t *testing.T) {
	test.TestDocument("batchchtype", t)
}
