//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

// ------------------- v1 -----------------------
func TestResumeV1Upload(t *testing.T) {
	test.RunCmdWithError("delete", test.Bucket, "qshell_rput_5M")

	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
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

	path, err = test.CreateTempFile(5*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	// not overwrite
	result, errs = test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
		"--mimetype", "image/png",
		"--storage", "1",
		"--worker", "4")
	if !strings.Contains(errs, "upload error:file exists") {
		t.Fatal(result)
	}

	// overwrite
	result, errs = test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
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

func TestResumeV1UploadWithUploadHost(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_uploadHost_5M", path,
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

	path, err = test.CreateTempFile(5*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}
}

// ------------------- v2 -----------------------
func TestResumeV2Upload(t *testing.T) {
	test.RunCmdWithError("delete", test.Bucket, "qshell_rput_5M")

	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
		"--mimetype", "image/jpg",
		"--storage", "0",
		"--resumable-api-v2")
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

	path, err = test.CreateTempFile(5*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	// not overwrite
	_, errs = test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
		"--mimetype", "image/png",
		"--storage", "1",
		"--resumable-api-v2")
	if !strings.Contains(errs, "upload error:file exists") {
		t.Fatal(errs)
	}

	// overwrite
	result, errs = test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path,
		"--mimetype", "image/png",
		"--storage", "1",
		"--overwrite",
		"--resumable-api-v2")

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "MimeType: image/png") {
		t.Fatal(result)
	}
}

func TestResumeV2UploadWithUploadHost(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_v2_uploadHost_5M", path,
		"--mimetype", "image/jpg",
		"--storage", "0",
		"--up-host", "up-na0.qiniup.com",
		"--overwrite",
		"--resumable-api-v2")
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

	path, err = test.CreateTempFile(5*1024 + 1)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}
}

func TestResumeUploadWithWrongUploadHost(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	_, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_v1_uploadHost_5M", path,
		"--mimetype", "image/jpg",
		"--storage", "0",
		"--up-host", "up-mock.qiniup.com",
		"--overwrite")
	if !strings.Contains(errs, "dial tcp: lookup up-mock.qiniup.com: no such host") &&
		!strings.Contains(errs, "Upload file error") {
		t.Fail()
	}
}

func TestResumeUploadNoExistBucket(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	_, errs := test.RunCmdWithError("rput", test.BucketNotExist, "qshell_rput_5M", path)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestResumeUploadNoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("rput")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestResumeUploadNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("rput", test.Bucket)
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestResumeUploadNoLocalFilePath(t *testing.T) {
	_, errs := test.RunCmdWithError("rput", test.Bucket, test.Key)
	if !strings.Contains(errs, "LocalFile can't empty") {
		t.Fail()
	}
}

func TestResumeUploadDocument(t *testing.T) {
	test.TestDocument("rput", t)
}
