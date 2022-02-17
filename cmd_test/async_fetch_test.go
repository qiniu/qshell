package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestAsyncFetch(t *testing.T) {
	path, err := test.CreateFileWithContent("async_fetch.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create path error:", err)
	}

	successLogPath, err := test.CreateFileWithContent("async_fetch_success_log.txt", "")
	if err != nil {
		t.Fatal("create successLogPath error:", err)
	}

	failLogPath, err := test.CreateFileWithContent("async_fetch_fail_log.txt", "")
	if err != nil {
		t.Fatal("create failLogPath error:", err)
	}

	_, errs := test.RunCmdWithError("abfetch", test.Bucket,
		"-i", path,
		"-s", successLogPath,
		"-e", failLogPath,
		"-g", "1",
		"-c", "2")
	if len(errs) > 0 {
		t.Fail()
	}

	if !test.IsFileHasContent(successLogPath) {
		t.Fail()
	}
	test.RemoveFile(successLogPath)

	if !test.IsFileHasContent(failLogPath) {
		t.Fail()
	}
	test.RemoveFile(failLogPath)
}

func TestAsyncFetchNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("abfetch")
	if !strings.Contains(err, "bucket can't empty") {
		t.Fail()
	}
}

func TestAsyncFetchDocument(t *testing.T) {
	result, _ := test.RunCmdWithError("abfetch", test.Bucket, test.DocumentOption)
	if strings.HasPrefix(result, "# 简介\n`abfetch`") {
		t.Fail()
	}
}
