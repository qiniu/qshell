package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestCdnRefreshFile(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result, _ := test.RunCmdWithError("cdnrefresh", "-i", path, "--qps", "1", "--size", "2")
	if !strings.Contains(result, "CDN refresh Code: 200, Info: success") {
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

	_, errs := test.RunCmdWithError("cdnrefresh", "--dirs", "-i", path, "--qps", "1", "--size", "2")
	if len(errs) > 0 {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnRefreshDocument(t *testing.T) {
	test.TestDocument("cdnrefresh", t)
}
