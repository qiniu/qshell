package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestBucketDomain(t *testing.T) {
	result, errs := test.RunCmdWithError("domains", test.Bucket)
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.BucketDomain) {
		t.Fatal("no expected domain:%", test.BucketDomain)
	}

	return
}

func TestBucketDomainDocument(t *testing.T) {
	result, errs := test.RunCmdWithError("domains", test.Bucket)
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.BucketDomain) {
		t.Fatal("no expected domain:%", test.BucketDomain)
	}

	return
}

func TestBucketList(t *testing.T) {
	result, errs := test.RunCmdWithError("listbucket", test.Bucket, "--prefix", "hello")
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.Key) {
		t.Fatal("no expected key:% but not exist", test.BucketDomain)
	}

	return
}

func TestBucketListToFile(t *testing.T) {
	rootPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get root path error:", err)
		return
	}
	file := filepath.Join(rootPath, test.Bucket + "_listbucket.txt")
	_, errs := test.RunCmdWithError("listbucket", test.Bucket, "--prefix", "hello", "-o", file)
	defer test.RemoveFile(file)

	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !test.IsFileHasContent(file) {
		t.Fatal("list bucket to file error: file empty")
	}

	return
}

func TestBucketListNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("listbucket")
	if !strings.Contains(err, "Bucket can't empty") {
		t.Fail()
	}
}
