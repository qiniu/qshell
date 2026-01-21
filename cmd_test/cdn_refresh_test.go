//go:build integration

package cmd

import (
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func TestCdnRefreshFile(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result, errString := test.RunCmdWithError("cdnrefresh", "-i", path, "--qps", "1", "--size", "2", "-d")
	if !strings.Contains(result, "CDN refresh Code: 200, FlowInfo: success") &&
		!strings.Contains(errString, "count limit error") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnRefreshDirs(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errString := test.RunCmdWithError("cdnrefresh", "--dirs", "-i", path, "--qps", "1", "--size", "2")
	if len(errString) > 0 && !strings.Contains(errString, "count limit error") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnRefreshDocument(t *testing.T) {
	test.TestDocument("cdnrefresh", t)
}
