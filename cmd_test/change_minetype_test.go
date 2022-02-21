package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/jpeg")
	if len(errs) > 0 {
		t.Fail()
	}

	_, errs = test.RunCmdWithError("chgm", test.Bucket, test.Key, "image/png")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestMimeTypeNoExistBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.BucketNotExist, test.Key, "image/jpeg")
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMimeTypeNoExistKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.KeyNotExist, "image/jpeg")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestMimeTypeNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestMimeTypeNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestMimeTypeNoMimeType(t *testing.T) {
	_, errs := test.RunCmdWithError("chgm", test.Bucket, test.Key)
	if !strings.Contains(errs, "MimeType can't empty") {
		t.Fail()
	}
}

func TestMimeTypeDocument(t *testing.T) {
	test.TestDocument("chgm", t)
}

func TestBatchChangeMimeType(t *testing.T) {
	batchConfig := ""
	for _, key := range test.Keys {
		batchConfig += key + "\t" + "image/jpeg" + "\n"
	}

	path, err := test.CreateFileWithContent("batch_chgm.txt", batchConfig)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("batchchgm", test.Bucket, "-i", path, "-y")
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestBatchMimeTypeDocument(t *testing.T) {
	test.TestDocument("batchchgm", t)
}