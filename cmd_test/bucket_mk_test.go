//go:build integration

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestMkBucket(t *testing.T) {

	_, errs := test.RunCmdWithError("mkbucket", test.Bucket, "--region", "z1", "--private")
	if len(errs) == 0 {
		t.Fatal("should return bucket exists")
	}

	if !strings.Contains(errs, "the bucket already exists") {
		t.Fatal("expected error:bucket exists, but:" + errs)
	}

	return
}

func TestMkBucketNotExistRegion(t *testing.T) {

	_, errs := test.RunCmdWithError("mkbucket", test.Bucket, "--region", "z10", "--private")
	if len(errs) == 0 {
		t.Fatal("should return bucket exists")
	}

	if !strings.Contains(errs, "invalid region parameter") {
		t.Fatal("expected error:invalid region parameter, but:" + errs)
	}

	return
}

func TestMkBucketNoBucket(t *testing.T) {
	_, err := test.RunCmdWithError("mkbucket")
	if !strings.Contains(err, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestMkBucketDocument(t *testing.T) {
	test.TestDocument("mkbucket", t)
}
