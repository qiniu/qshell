package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestBuckets(t *testing.T) {
	result, errs := test.RunCmdWithError("buckets")
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.Bucket) {
		t.Fatal("no expected bucket:%", test.Bucket)
	}
	return
}

func TestBucketsDocument(t *testing.T) {
	test.TestDocument("buckets", t)
}
