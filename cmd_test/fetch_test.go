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
	if !strings.Contains(errs, "Bucket can't empty") {
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
	for _, domain := range test.BucketObjectDomains {
		name := "batch_fetch_" + filepath.Base(domain)
		batchConfig += domain + "\t" + name + "\n"
	}
	path, err := test.CreateFileWithContent("batch_fetch.txt", batchConfig)
	if err != nil {
		t.Fatal("create batch fetch config file error:", err)
	}

	result, errs := test.RunCmdWithError("batchfetch", test.Bucket,
		"-i", path,
		"-c", "2")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	result, errs = test.RunCmdWithError("prefop", result)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchFetchDocument(t *testing.T) {
	test.TestDocument("batchfetch", t)
}
