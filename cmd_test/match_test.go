package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestMatch(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.Bucket, test.Key, objectPath)
	if len(errs) > 0 {
		t.Fail()
	}
}

func TestMatchNoExistBucket(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.BucketNotExist, test.Key, objectPath)
	if !strings.Contains(errs, "no such bucket") {
		t.Fail()
	}
}

func TestMatchNoExistSrcKey(t *testing.T) {
	resultPath, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result path error:", err)
	}
	objectPath := filepath.Join(resultPath, test.Key)
	_, errs := test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	_, errs = test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, objectPath)
	if !strings.Contains(errs, "no such file or directory") {
		t.Fail()
	}
}

func TestMatchNoExistLocalFile(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, "/user/desktop/a.txt")
	if !strings.Contains(errs, "open /user/desktop/a.txt: no such file or directory") {
		t.Fail()
	}
}

func TestMatchWithEmptyBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("match", "", test.Key, "/user/desktop/a.txt")
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestMatchNoKey(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, "", "/user/desktop/a.txt")
	if !strings.Contains(errs, "Key can't empty") {
		t.Fail()
	}
}

func TestMatchNoLocalFile(t *testing.T) {
	_, errs := test.RunCmdWithError("match", test.Bucket, test.KeyNotExist, "")
	if !strings.Contains(errs, "LocalFile can't empty") {
		t.Fail()
	}
}

func TestMatchDocument(t *testing.T) {
	test.TestDocument("match", t)
}

func TestBatchMatch(t *testing.T) {

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	objectPath := filepath.Join(resultDir, test.Key)
	_, _ = test.RunCmdWithError("get", test.Bucket, test.Key,
		"-o", objectPath)
	defer test.RemoveFile(objectPath)

	successLogPath := filepath.Join(resultDir, "batch_success.txt")
	failLogPath := filepath.Join(resultDir, "batch_fail.txt")

	path, err := test.CreateFileWithContent("batch_match.txt", test.KeysString)
	if err != nil {
		t.Fatal("create batch match config file error:", err)
	}

	test.RunCmdWithError("batchmatch", test.Bucket, resultDir,
		"-i", path,
		"--success-list", successLogPath,
		"--failure-list", failLogPath,
		"--worker", "4")
	defer func() {
		test.RemoveFile(successLogPath)
		test.RemoveFile(failLogPath)
	}()

	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("batch result: success log to file error: file empty")
	}

	if !test.IsFileHasContent(failLogPath) {
		t.Fatal("batch result: fail log  to file error: file empty")
	}
}

func TestBatchMatchDocument(t *testing.T) {
	test.TestDocument("batchmatch", t)
}
