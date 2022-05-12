//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestMkBucket(t *testing.T) {

	_, errs := test.RunCmdWithError("mkbucket", test.Bucket, "--region", "z1", "--private")
	if len(errs) == 0 {
		t.Fatal("should return bucket exists")
	}

	if !strings.Contains(errs, "error:bucket exists") {
		t.Fatal("expected error:bucket exists, but:" + errs)
	}

	return
}

func TestMkBucketNotExistRegion(t *testing.T) {

	_, errs := test.RunCmdWithError("mkbucket", test.Bucket, "--region", "z10", "--private")
	if len(errs) == 0 {
		t.Fatal("should return bucket exists")
	}

	if !strings.Contains(errs, "error:invalid region parameter") {
		t.Fatal("expected error:error:invalid region parameter, but:" + errs)
	}

	return
}

func TestMkBucketNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("mkbucket")
	if !strings.Contains(err, "Bucket can't empty") {
		t.Fail()
	}
}

func TestMkBucketDocument(t *testing.T) {
	test.TestDocument("mkbucket", t)
}