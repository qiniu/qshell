//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestFormUpload(t *testing.T) {
	test.RunCmdWithError("delete", test.Bucket, "qshell_fput_1M")

	path, err := test.CreateTempFile(1 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("fput", test.Bucket, "qshell_fput_1M", path,
		"--mimetype", "image/jpg",
		"--storage", "0")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "MimeType: image/jpg") {
		t.Fatal(result)
	}

	path, err = test.CreateTempFile(1*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	// not overwrite
	_, errs = test.RunCmdWithError("fput", test.Bucket, "qshell_fput_1M", path,
		"--mimetype", "image/jpg",
		"--storage", "1")
	if !strings.Contains(errs, "upload error:file exists") {
		t.Fatal(errs)
	}

	// overwrite
	result, errs = test.RunCmdWithError("fput", test.Bucket, "qshell_fput_1M", path,
		"--mimetype", "image/png",
		"--storage", "1",
		"--overwrite")

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "MimeType: image/png") {
		t.Fatal(result)
	}
}

func TestFormUploadWithUploadHost(t *testing.T) {
	path, err := test.CreateTempFile(1 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("fput", test.Bucket, "qshell_fput_uploadHost_1M", path,
		"--mimetype", "image/jpg",
		"--storage", "0",
		"--up-host", "up-na0.qiniup.com",
		"--overwrite")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "MimeType: image/jpg") {
		t.Fatal(result)
	}

	path, err = test.CreateTempFile(1*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}
}

func TestFormUploadWithWrongUploadHost(t *testing.T) {
	path, err := test.CreateTempFile(1 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	_, errs := test.RunCmdWithError("fput", test.Bucket, "qshell_fput_uploadHost_1M", path,
		"--mimetype", "image/jpg",
		"--storage", "0",
		"--up-host", "up-mock.qiniup.com",
		"--overwrite")
	if !strings.Contains(errs, "dial tcp: lookup up-mock.qiniup.com: no such host") {
		t.Fail()
	}
}

func TestFormUploadNoExistBucket(t *testing.T) {
	path, err := test.CreateTempFile(1 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	_, errs := test.RunCmdWithError("fput", test.BucketNotExist, "qshell_fput_1M", path)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestFormUploadNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("fput")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestFormUploadNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("fput", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestFormUploadNoLocalFilePath(t *testing.T) {
	_, errs := test.RunCmdWithError("fput", test.Bucket, test.Key)
	if !strings.Contains(errs, "LocalFile can't empty") {
		t.Fail()
	}
}

func TestFormUploadDocument(t *testing.T) {
	test.TestDocument("fput", t)
}
