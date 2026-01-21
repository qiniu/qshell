//go:build integration

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
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
	if !strings.Contains(err, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestBucketDocument(t *testing.T) {
	test.TestDocument("bucket", t)
}
