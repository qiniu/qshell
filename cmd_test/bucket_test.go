package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestBucket(t *testing.T) {

	ret, errs := test.RunCmdWithError("bucket", test.Bucket)
	if len(errs) > 0 {
		t.Fatal("get bucket info error:" + errs)
	}

	if !strings.Contains(ret, test.Bucket) {
		t.Fatal("UnExcepted bucket info:" + ret)
	}

	return
}

func TestBucketNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("bucket")
	if !strings.Contains(err, "Bucket can't empty") {
		t.Fail()
	}
}

func TestBucketDocument(t *testing.T) {
	test.TestDocument("bucket", t)
}
