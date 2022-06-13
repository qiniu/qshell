//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func aTestCdnPrefetch(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_prefetch.txt", test.BucketObjectDomainsString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result := test.RunCmd(t, "cdnprefetch", "-i", path, "--qps", "1", "--size", "2")
	if !strings.Contains(result, "CDN prefetch Code: 200, FlowInfo: success") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnPrefetchDocument(t *testing.T) {
	test.TestDocument("cdnprefetch", t)
}
