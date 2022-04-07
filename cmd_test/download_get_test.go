//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	TestCopy(t)

	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	path := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", path)
	defer test.RemoveFile(path)

	if len(errs) > 0 {
		t.Fail()
	}
	if !test.IsFileHasContent(path) {
		t.Fatal("get file content can't empty")
	}
}

func TestGetWithDomain(t *testing.T) {
	TestCopy(t)

	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	path := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"--domain", test.BucketDomain,
		"-o", path)
	defer test.RemoveFile(path)

	if len(errs) > 0 {
		t.Fail()
	}
	if !test.IsFileHasContent(path) {
		t.Fatal("get file content can't empty")
	}
}

func TestGetNoExistDomain(t *testing.T) {
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"--domain", "qiniu.mock.com")
	if !strings.Contains(errs, "lookup qiniu.mock.com: no such host") {
		t.Fail()
	}
}

func TestGetNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("get", test.BucketNotExist, test.Key)
	if !strings.Contains(errs, "query region error, no such bucket") {
		t.Fail()
	}
}

func TestGetNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("get", test.Bucket, test.KeyNotExist)
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestGetNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("get")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestGetNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("get", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestGetDocument(t *testing.T) {
	test.TestDocument("get", t)
}
