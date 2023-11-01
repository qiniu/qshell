//go:build unit

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestBatchSign(t *testing.T) {
	path, err := test.CreateFileWithContent("batch_sign.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create batch sign config file error:", err)
	}

	resultDir, err := test.ResultPath()
	if err != nil {
		t.Fatal("get result dir error:", err)
	}

	resultLogPath := filepath.Join(resultDir, "batch_result.txt")

	result, errs := test.RunCmdWithError("batchsign",
		"-i", path,
		"--outfile", resultLogPath,
	)
	if len(errs) > 0 {
		t.Fail()
	}

	if len(result) == 0 {
		t.Fail()
	}

	defer func() {
		test.RemoveFile(resultLogPath)
	}()

	if !test.IsFileHasContent(resultLogPath) {
		t.Fatal("batch result: output  to file error: file empty")
	}
}

func TestBatchSignDocument(t *testing.T) {
	test.TestDocument("batchsign", t)
}
