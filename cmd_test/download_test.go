package cmd

import (
	"encoding/json"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadWithKeyFile(t *testing.T) {
	rootPath, err := test.RootPath()
	if err != nil {
		t.Fatal("get root path err:", err)
		return
	}

	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", test.Keys)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(keysFilePath)

	d := config.Download{
		ThreadCount: 4,
		KeyFile:     keysFilePath,
		DestDir:     filepath.Join(rootPath, "download"),
		Bucket:      test.Bucket,
		Prefix:      "",
		Suffixes:    "",
		IoHost:      "",
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
		//LogSetting: &config.LogSetting{
		//	LogLevel:  "",
		//	LogFile:   "",
		//	LogRotate: 0,
		//	LogStdout: false,
		//},
		Tasks: &config.Tasks{},
		Retry: &config.Retry{},
	}
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		t.Fatal("json marshal error:", err)
		return
	}

	path, err := test.CreateFileWithContent("download_config.txt", string(data))
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(path)

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path, "-d")
	if !strings.Contains(errs, "CDN refresh Code: 200, Info: success") {
		t.Fail()
	}

	return
}


func TestDownload(t *testing.T) {
	rootPath, err := test.RootPath()
	if err != nil {
		t.Fatal("get root path err:", err)
		return
	}

	d := config.Download{
		ThreadCount: 4,
		KeyFile:     "",
		DestDir:     filepath.Join(rootPath, "download"),
		Bucket:      test.Bucket,
		Prefix:      "",
		Suffixes:    "",
		IoHost:      "",
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
		//LogSetting: &config.LogSetting{
		//	LogLevel:  "",
		//	LogFile:   "",
		//	LogRotate: 0,
		//	LogStdout: false,
		//},
		Tasks: &config.Tasks{},
		Retry: &config.Retry{},
	}
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		t.Fatal("json marshal error:", err)
		return
	}

	path, err := test.CreateFileWithContent("download_config.txt", string(data))
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(path)

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path, "-d")
	if !strings.Contains(errs, "CDN refresh Code: 200, Info: success") {
		t.Fail()
	}

	return
}