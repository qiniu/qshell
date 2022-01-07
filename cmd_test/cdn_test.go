package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestCdnPrefetch(t *testing.T) {

	path, err := test.CreateFileWithContent("cdn_prefetch.txt", `
http://if-pbl.qiniudn.com/hello1.txt
http://if-pbl.qiniudn.com/hello2.txt
http://if-pbl.qiniudn.com/hello3.txt
http://if-pbl.qiniudn.com/hello4.txt
http://if-pbl.qiniudn.com/hello5.txt
http://if-pbl.qiniudn.com/hello6.txt
http://if-pbl.qiniudn.com/hello7.txt
`)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("cdnprefetch", "-i", path, "--qps", "1", "--size", "1")
	if !strings.Contains(errs, "400014") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}

func TestCdnRefreshFile(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", `
http://if-pbl.qiniudn.com/hello1.txt
http://if-pbl.qiniudn.com/hello2.txt
http://if-pbl.qiniudn.com/hello3.txt
http://if-pbl.qiniudn.com/hello4.txt
http://if-pbl.qiniudn.com/hello5.txt
http://if-pbl.qiniudn.com/hello6.txt
http://if-pbl.qiniudn.com/hello7.txt
`)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("cdnrefresh", "-i", path, "--qps", "1", "--size", "1")
	if !strings.Contains(errs, "400014") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}


func TestCdnRefreshDirs(t *testing.T) {
	path, err := test.CreateFileWithContent("cdn_refresh.txt", `
http://if-pbl.qiniudn.com/hello1.txt
http://if-pbl.qiniudn.com/hello2.txt
http://if-pbl.qiniudn.com/hello3.txt
http://if-pbl.qiniudn.com/hello4.txt
http://if-pbl.qiniudn.com/hello5.txt
http://if-pbl.qiniudn.com/hello6.txt
http://if-pbl.qiniudn.com/hello7.txt
`)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	_, errs := test.RunCmdWithError("cdnrefresh", "--dirs", path, "--qps", "1", "--size", "1")
	if !strings.Contains(errs, "400014") {
		t.Fail()
	}

	test.RemoveFile(path)

	return
}