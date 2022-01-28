package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestFormUpload(t *testing.T) {
	path, err := test.CreateTempFile(1 * 1024)
	if err != nil {
		t.Fatal("create form upload file error:", err)
	}

	result, errs := test.RunCmdWithError("fput", test.Bucket, "qshell_fput_1M", path, "-d")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestResumeUpload(t *testing.T) {
	path, err := test.CreateTempFile(5 * 1024)
	if err != nil {
		t.Fatal("create resume upload file error:", err)
	}

	result, errs := test.RunCmdWithError("rput", test.Bucket, "qshell_rput_5M", path)
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestQUpload(t *testing.T) {
	fileSizeList := []int{1, 32, 64, 256, 512, 1024, 2 * 1024, 4 * 1024, 5 * 1024, 8 * 1024, 10 * 1024}
	for _, size := range fileSizeList {
		test.CreateTempFile(size)
	}

	fileDir, err := test.TempPath()
	if err != nil {
		t.Fatal("create upload temp file error:", err)
	}

	cfgContent := fmt.Sprintf(`{
	"bucket": "%s",
	"src_dir": "%s",
	"overwrite": "true",
	"work_count": 4
}`, test.Bucket, fileDir)
	cfgFile, err  := test.CreateFileWithContent("upload_cfg.json", cfgContent)
	defer test.RemoveFile(cfgFile)

	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result, errs := test.RunCmdWithError("qupload", cfgFile)
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestQUpload2(t *testing.T) {
	fileSizeList := []int{1, 32, 64, 256, 512, 1024, 2 * 1024, 4 * 1024, 5 * 1024, 8 * 1024, 10 * 1024}
	for _, size := range fileSizeList {
		test.CreateTempFile(size)
	}

	fileDir, err := test.TempPath()
	if err != nil {
		t.Fatal("create upload temp file error:", err)
	}

	result, errs := test.RunCmdWithError("qupload2",
		"--bucket", test.Bucket,
		"--src-dir", fileDir)
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestSync(t *testing.T) {

	url := "https://qshell-na0.qiniupkg.com/10240K.tmp"
	result, errs := test.RunCmdWithError("sync", url, test.Bucket, "-d")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}
