//go:build unit

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

func TestBatchSign(t *testing.T) {
	path, err := test.CreateFileWithContent("batch_sign.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create batch sign config file error:", err)
	}

	result, errs := test.RunCmdWithError("batchsign", "-i", path)
	if len(errs) > 0 {
		t.Fail()
	}

	if len(result) == 0 {
		t.Fail()
	}
}

func TestBatchSignDocument(t *testing.T) {
	test.TestDocument("batchsign", t)
}
