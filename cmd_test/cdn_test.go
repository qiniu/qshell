package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestCdnPrefetch(t *testing.T) {

	path, err := test.CreateFileWithContent("cdn_prefetch.txt", `
https://qshell-na0.qiniupkg.com/hello1.json
https://qshell-na0.qiniupkg.com/hello2.json
https://qshell-na0.qiniupkg.com/hello3.json
https://qshell-na0.qiniupkg.com/hello4.json
https://qshell-na0.qiniupkg.com/hello5.json
https://qshell-na0.qiniupkg.com/hello6.json
https://qshell-na0.qiniupkg.com/hello7.json
`)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	result := test.RunCmd(t, "cdnprefetch", "-i", path, "--qps", "1", "--size", "2")
	if !strings.Contains(result, "CDN prefetch Code: 200, Info: success") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnRefreshFile(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", `
https://qshell-na0.qiniupkg.com/hello1.json
https://qshell-na0.qiniupkg.com/hello2.json
https://qshell-na0.qiniupkg.com/hello3.json
https://qshell-na0.qiniupkg.com/hello4.json
https://qshell-na0.qiniupkg.com/hello5.json
https://qshell-na0.qiniupkg.com/hello6.json
https://qshell-na0.qiniupkg.com/hello7.json
`)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("cdnrefresh", "-i", path, "--qps", "1", "--size", "2")
	if !strings.Contains(errs, "CDN refresh Code: 200, Info: success") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}


func TestCdnRefreshDirs(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", `
https://qshell-na0.qiniupkg.com/hello1.json
https://qshell-na0.qiniupkg.com/hello2.json
https://qshell-na0.qiniupkg.com/hello3.json
https://qshell-na0.qiniupkg.com/hello4.json
https://qshell-na0.qiniupkg.com/hello5.json
https://qshell-na0.qiniupkg.com/hello6.json
https://qshell-na0.qiniupkg.com/hello7.json
`)
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