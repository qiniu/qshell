//go:build unit

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestAliBucketListNoDataCenter(t *testing.T) {
	_, errs := test.RunCmdWithError("alilistbucket")
	if !strings.Contains(errs, "DataCenter can't be empty") {
		t.Fatal("empty DataCenter check error")
	}
	return
}

func TestAliBucketListNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("alilistbucket", "DataCenter")
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fatal("empty Bucket check error")
	}
	return
}

func TestAliBucketListNoAccessKeyId(t *testing.T) {
	_, errs := test.RunCmdWithError("alilistbucket", "DataCenter", "Bucket")
	if !strings.Contains(errs, "AccessKeyId can't be empty") {
		t.Fatal("empty AccessKeyId check error")
	}

	return
}

func TestAliBucketListNoAccessKeySecret(t *testing.T) {
	_, errs := test.RunCmdWithError("alilistbucket", "DataCenter", "Bucket", "AccessKeyId")
	if !strings.Contains(errs, "AccessKeySecret can't be empty") {
		t.Fatal("empty AccessKeySecret check error")
	}
	return
}

func TestAliBucketListDocument(t *testing.T) {
	test.TestDocument("alilistbucket", t)
}
