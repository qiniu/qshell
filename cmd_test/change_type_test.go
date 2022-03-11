package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestChangeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chtype", test.Bucket, test.Key, "0")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chtype", test.Bucket, test.Key, "1")
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
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchchtype", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchChangeTypeDocument(t *testing.T) {
	test.TestDocument("batchchtype", t)
}
