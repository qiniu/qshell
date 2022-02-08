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

	if len(result) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("delete", test.Bucket, key)
	if len(errs) > 0 {
		t.Fail()
	}
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
