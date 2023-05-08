//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	key := "fetch_key.json"
	result, errs := test.RunCmdWithError("fetch", test.BucketObjectDomain, test.Bucket, "-k", key)
	if len(errs) > 0 {
		t.Fail()
	}

	if !strings.Contains(result, key) {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("delete", test.Bucket, key)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestFetchNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("fetch", test.BucketObjectDomain, test.BucketNotExist, "-k", test.Key)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestFetchNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("fetch", test.BucketObjectDomain)
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestFetchNoKey(t *testing.T) {
	result, _ := test.RunCmdWithError("fetch", test.BucketObjectDomain, test.Bucket)
	if !strings.Contains(result, "Key:FvySxBAiQRAd1iSF4XrC4SrDrhff") {
		t.Fail()
	}
}

func TestFetchDocument(t *testing.T) {
	test.TestDocument("fetch", t)
}

func TestBatchFetch(t *testing.T) {
	batchConfig := ""
	domains := test.BucketObjectDomains
	domains = append(domains, "https://qshell-na0.qiniupkg.com/hello10.json")
	for _, domain := range domains {
		name := "batch_fetch_" + filepath.Base(domain)
		batchConfig += domain + "\t" + name + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_copy_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_copy_fail.txt")

	path, err := test.CreateFileWithContent("batch_fetch.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch fetch config file error:", err)
	}

	test.RunCmdWithError("batchfetch", test.Bucket,
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

	result, _ := test.RunCmdWithError("batchfetch", test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4",
		"-y")
	if !strings.Contains(result, "Skip") {
		t.Fatal("batch result: redo should skip")
	}
}

func TestBatchFetchWithRecord(t *testing.T) {
	batchConfig := ""
	domains := []string{test.BucketObjectDomain, "https://qshell-na0.qiniupkg.com/hello10.json"}
	for _, domain := range domains {
		name := "batch_fetch_" + filepath.Base(domain)
		batchConfig += domain + "\t" + name + "\n"
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	successLogPath := filepath.Join(resultDir, "batch_copy_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_copy_fail.txt")

	path, err := test.CreateFileWithContent("batch_fetch.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch fetch config file error:", err)
	}

	test.RunCmdWithError("batchfetch", test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--enable-record",
		"--worker", "1",
		"-y")
	defer func() {
		test.RemoveFile(successLogPath)
		test.RemoveFile(failLogPath)
	}()

	result, _ := test.RunCmdWithError("batchfetch", test.Bucket,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--enable-record",
		"--worker", "4",
		"-y",
		"-d")
	if !strings.Contains(result, "Skip") {
		t.Fatal("batch result: redo should skip")
	}
}

func TestBatchFetchDocument(t *testing.T) {
	test.TestDocument("batchfetch", t)
}
