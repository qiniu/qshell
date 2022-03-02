package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestQUpload(t *testing.T) {
	fileSizeList := []int{1, 32, 64, 256, 512, 1024, 2 * 1024, 4 * 1024, 5 * 1024, 8 * 1024, 10 * 1024}
	fileSizeList = []int{1, 2}
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
	"log_stdout": "true",
	"overwrite": "true",
	"check_exists": "true",
	"check_size": "true",
	"work_count": 4
}`, test.Bucket, fileDir)
	cfgFile, err := test.CreateFileWithContent("upload_cfg.json", cfgContent)
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

func TestQUploadDocument(t *testing.T) {
	test.TestDocument("qupload", t)
}

func TestQUpload2WithSrcDir(t *testing.T) {
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
		"--src-dir", fileDir,
		"--rescan-local")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestQUpload2WithFileList(t *testing.T) {
	fileSizeList := []int{1, 32, 64, 256, 512, 1024, 2 * 1024, 4 * 1024, 5 * 1024, 8 * 1024, 10 * 1024}
	for _, size := range fileSizeList {
		test.CreateTempFile(size)
	}

	fileDir, err := test.TempPath()
	if err != nil {
		t.Fatal("create upload temp file error:", err)
	}

	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	fileListPath := filepath.Join(resultPath, "qupload2_file_list.txt")
	_, errs := test.RunCmdWithError("dircache", fileDir,
		"-o", fileListPath)
	if len(errs) > 0 {
		t.Fatal("upload2 dircache error:", err)
	}

	result, errs := test.RunCmdWithError("qupload2",
		"--bucket", test.Bucket,
		"--src-dir", fileDir,
		"--file-list", fileListPath, "-d")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestQUpload2NoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("qupload2")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fatal(errs)
	}
}

func TestQUpload2NoSrc(t *testing.T) {
	_, errs := test.RunCmdWithError("qupload2",
		"--bucket", test.Bucket)
	if !strings.Contains(errs, "SrcDir can't empty") {
		t.Fatal(errs)
	}
}

func TestQUpload2NotExistSrcDir(t *testing.T) {
	_, errs := test.RunCmdWithError("qupload2",
		"--bucket", test.Bucket,
		"--src-dir", "/Demo")
	if !strings.Contains(errs, "no such file or directory") {
		t.Fatal(errs)
	}
}

func TestQUpload2Document(t *testing.T) {
	test.TestDocument("qupload2", t)
}
