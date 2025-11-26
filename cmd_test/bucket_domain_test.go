//go:build integration

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestBucketDomain(t *testing.T) {
	result, errs := test.RunCmdWithError("domains", test.Bucket)
	if len(errs) > 0 {
		t.Fatal("error:", errs)
	}

	if !strings.Contains(result, test.BucketDomain) {
		t.Fatal("no expected domain:", result)
	}

	return
}

func TestBucketDomainNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("domains")
	if !strings.Contains(err, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestBucketDomainDocument(t *testing.T) {
	test.TestDocument("domains", t)
}
