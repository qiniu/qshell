//go:build integration

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

const fopObjectKey = "test_mv.mp4"
const fopObjectValue = "avthumb/mp4"

func TestFop(t *testing.T) {
	result, errs := test.RunCmdWithError("pfop", test.Bucket, fopObjectKey, fopObjectValue)
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	result, errs = test.RunCmdWithError("prefop", result)
	if len(errs) > 0 && !strings.Contains(errs, "400 Bad Request") {
		t.Fail()
	}
}

func TestFopNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop", test.BucketNotExist, fopObjectKey, fopObjectValue)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestFopNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop", test.Bucket, test.KeyNotExist, fopObjectValue)
	if !strings.Contains(errs, "NotFound") {
		t.Fail()
	}
}

func TestFopNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop")
	if !strings.Contains(errs, "Bucket can't be empty") {
		t.Fail()
	}
}

func TestFopNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop", test.Bucket)
	if !strings.Contains(errs, "Key can't be empty") {
		t.Fail()
	}
}

func TestFopNoFopValue(t *testing.T) {
	_, errs := test.RunCmdWithError("pfop", test.Bucket, fopObjectKey)
	if !strings.Contains(errs, "Fops can't be empty") {
		t.Fail()
	}
}

func TestFopDocument(t *testing.T) {
	test.TestDocument("pfop", t)
}

func TestPreFopNoID(t *testing.T) {
	_, errs := test.RunCmdWithError("prefop")
	if !strings.Contains(errs, "PersistentID can't be empty") {
		t.Fail()
	}
}

func TestPreFopDocument(t *testing.T) {
	test.TestDocument("prefop", t)
}
