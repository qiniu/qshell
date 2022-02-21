package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
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
	result, errs := test.RunCmdWithError("listbucket", test.Bucket)
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.Key) {
		t.Fatal("no expected key:% but not exist", test.BucketDomain)
	}

	return
}
