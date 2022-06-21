//go:build integration

package cmd

import (
	"encoding/json"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/operations"
	"path/filepath"
	"strings"
	"testing"
)

func TestQUpload(t *testing.T) {
	deleteFile(t, "test/1K.tmp")
	deleteFile(t, "test/32K.tmp")
	deleteFile(t, "test/64K.tmp")
	deleteFile(t, "test/256K.tmp")
	copyFile(t, test.Key, "test/512K.tmp")
	copyFile(t, test.Key, "test/1024K.tmp")

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

	successLogPath := filepath.Join(resultPath, "qupload2_success.txt")
	failLogPath := filepath.Join(resultPath, "qupload2_fail.txt")
	overwriteLogPath := filepath.Join(resultPath, "qupload2_overwrite.txt")
	logPath := filepath.Join(resultPath, "qupload2_log.txt")
	recordPath := filepath.Join(resultPath, "record")

	cfg := struct {
		operations.UploadConfig
		LogFile    string `json:"log_file"`
		RecordRoot string `json:"record_root"`
	}{
		UploadConfig: operations.UploadConfig{
			UpHost:                 "",
			SrcDir:                 fileDir,
			FileList:               "",
			IgnoreDir:              false,
			SkipFilePrefixes:       "",
			SkipPathPrefixes:       "",
			SkipFixedStrings:       "",
			SkipSuffixes:           "",
			FileEncoding:           "",
			Bucket:                 test.Bucket,
			ResumableAPIV2:         false,
			ResumableAPIV2PartSize: 0,
			PutThreshold:           0,
			KeyPrefix:              "test/",
			Overwrite:              true,
			CheckExists:            false,
			CheckHash:              true,
			CheckSize:              true,
			RescanLocal:            true,
			FileType:               0,
			DeleteOnSuccess:        false,
			DisableResume:          false,
			DisableForm:            false,
			WorkerCount:            4,
			RecordRoot:             "",
			Policy:                 nil,
		},
		LogFile:    logPath,
		RecordRoot: recordPath,
	}
	cfgContent, _ := json.Marshal(cfg)
	cfgFile, err := test.CreateFileWithContent("upload_cfg.json", string(cfgContent))
	defer test.RemoveFile(cfgFile)

	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	test.RunCmdWithError("qupload", cfgFile,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--overwrite-list", overwriteLogPath)
	defer func() {
		test.RemoveFile(successLogPath)
		test.RemoveFile(failLogPath)
		test.RemoveFile(overwriteLogPath)
		test.RemoveFile(logPath)
		test.RemoveFile(recordPath)
	}()

	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("batch result: success log to file error: file empty")
	}

	if test.IsFileHasContent(failLogPath) {
		t.Fatal("batch result: fail log to file error: file should empty")
	}

	if !test.IsFileHasContent(overwriteLogPath) {
		t.Fatal("batch result: overwrite log to file error: file empty")
	}

	if !test.IsFileHasContent(logPath) {
		t.Fatal("batch result: log to file error: file empty")
	}
}

func TestQUploadDocument(t *testing.T) {
	test.TestDocument("qupload", t)
}

func TestQUpload2WithSrcDir(t *testing.T) {
	deleteFile(t, "1K.tmp")
	deleteFile(t, "32K.tmp")
	deleteFile(t, "64K.tmp")
	deleteFile(t, "256K.tmp")

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
		"--rescan-local",
		"--check-exists")
	defer func() {
		test.RemoveFile(fileListPath)
	}()

	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestQUpload2WithFileList(t *testing.T) {
	deleteFile(t, "1K.tmp")
	deleteFile(t, "32K.tmp")
	deleteFile(t, "64K.tmp")
	deleteFile(t, "256K.tmp")
	copyFile(t, test.Key, "512K.tmp")
	copyFile(t, test.Key, "1024K.tmp")
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

	successLogPath := filepath.Join(resultPath, "qupload2_success.txt")
	failLogPath := filepath.Join(resultPath, "qupload2_fail.txt")
	overwriteLogPath := filepath.Join(resultPath, "qupload2_overwrite.txt")
	logPath := filepath.Join(resultPath, "qupload2_log.txt")
	recordPath := filepath.Join(resultPath, "record")

	fileListPath := filepath.Join(resultPath, "qupload2_file_list.txt")
	_, errs := test.RunCmdWithError("dircache", fileDir,
		"-o", fileListPath)
	if len(errs) > 0 {
		t.Fatal("upload2 dircache error:", err)
	}

	if err := test.AppendToFile(fileListPath, `
mock01.jpg	10485760	16455233472998522
mock02.jpg	10485760	16455233472998522
`); err != nil {
		t.Fatal("upload2 upload file list append error:", err)
	}

	test.RunCmdWithError("qupload2",
		"--bucket", test.Bucket,
		"--src-dir", fileDir,
		"--file-list", fileListPath,
		"--overwrite",
		"--check-exists",
		"--check-hash",
		"--check-size",
		"--file-type", "1",
		"--rescan-local", "false",
		"--ignore-dir", "",
		"--key-prefix", "test/",
		"--skip-file-prefixes", "",
		"--skip-fixed-strings", "",
		"--skip-path-prefixes", "",
		"--skip-suffixes", "",
		"--thread-count", "4",
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--overwrite-list", overwriteLogPath,
		"--record-root", recordPath,
		"--log-file", logPath,
		"--log-level", "debug",
		"--log-rotate", "10",
		"--up-host", "",
		"-d")

	defer func() {
		test.RemoveFile(successLogPath)
		test.RemoveFile(failLogPath)
		test.RemoveFile(overwriteLogPath)
		test.RemoveFile(logPath)
		test.RemoveFile(recordPath)
	}()

	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("batch result: success log to file error: file empty")
	}

	if !test.IsFileHasContent(failLogPath) {
		t.Fatal("batch result: fail log  to file error: file empty")
	}

	if !test.IsFileHasContent(overwriteLogPath) {
		t.Fatal("batch result: overwrite log to file error: file empty")
	}

	if !test.IsFileHasContent(logPath) {
		t.Fatal("batch result: log to file error: file empty")
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
	if !strings.Contains(errs, "invalid SrcDir:") {
		t.Fatal(errs)
	}
}

func TestQUpload2Document(t *testing.T) {
	test.TestDocument("qupload2", t)
}
