//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestAsyncFetch(t *testing.T) {
	TestBatchCopy(t)

	fetchKeys := test.Keys
	fetchKeys = append(fetchKeys, "hello10.json")
	content := ""
	for _, key := range fetchKeys {
		content += "http://" + test.BucketDomain + "/" + key + "\t" + "0" + "\t" + "fetch_" + key + "\n"
	}
	path, err := test.CreateFileWithContent("async_fetch.txt", content)
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

	test.RunCmdWithError("abfetch", test.Bucket,
		"-i", path,
		"-s", successLogPath,
		"-e", failLogPath,
		"-g", "1",
		"--file-type", "1",
		"-c", "2")
	defer func() {
		test.RemoveFile(failLogPath)
		test.RemoveFile(successLogPath)
	}()
	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("success log can't empty")
	}

	if !test.IsFileHasContent(failLogPath) {
		t.Fatal("fail log can't empty")
	}
}

func TestAsyncFetchNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("abfetch")
	if !strings.Contains(err, "Bucket can't empty") {
		t.Fail()
	}
}

func TestAsyncFetchDocument(t *testing.T) {
	test.TestDocument("abfetch", t)
}

func TestACheck(t *testing.T) {
	id := "eyJ6b25lIjoibmEwIiwicXVldWUiOiJTSVNZUEhVUy1KT0JTLVYzIiwicGFydF9pZCI6OSwib2Zmc2V0Ijo1NTEzMTU3fQ=="
	result, err := test.RunCmdWithError("acheck", test.Bucket, id, "-D")
	if len(err) > 0 && !strings.Contains(err, "incorrect zone") && len(result) == 0 {
		t.Fail()
	}
}

func TestACheckNoId(t *testing.T) {
	_, err := test.RunCmdWithError("acheck", test.Bucket)
	if !strings.Contains(err, "Id can't empty") {
		t.Fail()
	}
}

func TestACheckNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("acheck")
	if !strings.Contains(err, "Bucket can't empty") {
		t.Fail()
	}
}

func TestACheckDocument(t *testing.T) {
	test.TestDocument("acheck", t)
}
