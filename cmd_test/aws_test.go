//go:build unit

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestAwsBucketListNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("awslist")
	if !strings.Contains(errs, "AWS bucket can't be empty") {
		t.Fatal("empty Bucket check error")
	}

	return
}

func TestAwsBucketListNoRegion(t *testing.T) {
	_, errs := test.RunCmdWithError("awslist", "bucket")
	if !strings.Contains(errs, "AWS region can't be empty") {
		t.Fatal("empty Region check error")
	}

	return
}

func TestAwsBucketListNoId(t *testing.T) {
	_, errs := test.RunCmdWithError("awslist", "bucket", "region")
	if !strings.Contains(errs, "AWS ID and SecretKey can't be empty") {
		t.Fatal("empty AWS ID check error")
	}

	return
}

func TestAwsBucketListNoSecret(t *testing.T) {
	_, errs := test.RunCmdWithError("awslist", "bucket", "region",
		"--aws-id", "id")
	if !strings.Contains(errs, "AWS ID and SecretKey can't be empty") {
		t.Fatal("empty Bucket check error")
	}

	return
}

func TestAwsBucketListDocument(t *testing.T) {
	test.TestDocument("awslist", t)
}

func TestAwsFetchNoAwsBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("awsfetch")
	if !strings.Contains(errs, "AWS bucket can't be empty") {
		t.Fatal("empty AWS Bucket check error")
	}

	return
}

func TestAwsFetchNoAwsRegion(t *testing.T) {
	_, errs := test.RunCmdWithError("awsfetch", "bucket")
	if !strings.Contains(errs, "AWS region can't be empty") {
		t.Fatal("empty Region check error")
	}

	return
}

func TestAwsFetchNoQiniuBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("awsfetch", "bucket", "region")
	if !strings.Contains(errs, "Qiniu bucket can't be empty") {
		t.Fatal("empty Qiniu Bucket check error")
	}

	return
}

func TestAwsFetchNoId(t *testing.T) {
	_, errs := test.RunCmdWithError("awsfetch", "bucket", "region", "qiniu_bucket")
	if !strings.Contains(errs, "AWS ID and SecretKey can't be empty") {
		t.Fatal("empty AWS ID check error")
	}

	return
}

func TestAwsFetchNoSecret(t *testing.T) {
	_, errs := test.RunCmdWithError("awsfetch", "bucket", "region", "qiniu_bucket",
		"--aws-id", "id")
	if !strings.Contains(errs, "AWS ID and SecretKey can't be empty") {
		t.Fatal("empty Bucket check error")
	}

	return
}

func TestAwsFetchDocument(t *testing.T) {
	test.TestDocument("awsfetch", t)
}
