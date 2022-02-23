package cmd

import (
	"encoding/json"
	"errors"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadWithKeyFile(t *testing.T) {
	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", test.KeysString)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(keysFilePath)

	cfg := &config.Download{
		ThreadCount: 4,
		KeyFile:     keysFilePath,
		Bucket:      test.Bucket,
		Prefix:      "hell",
		Suffixes:    ".json",
		IoHost:      "",
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(path)

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir)
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadFromBucket(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: 4,
		KeyFile:     "",
		Bucket:      test.Bucket,
		Prefix:      "hello3,hello5,hello7",
		Suffixes:    "",
		IoHost:      test.BucketDomain,
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func(){
		test.RemoveFile(path)
		test.RemoveFile(cfg.LogFile)
	}()

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir)
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadNoBucket(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: 4,
		KeyFile:     "",
		Bucket:      "",
		Prefix:      "hello3,hello5,hello7",
		Suffixes:    "",
		IoHost:      test.BucketDomain,
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func(){
		test.RemoveFile(path)
		test.RemoveFile(cfg.LogFile)
	}()

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path)
	if !strings.Contains(errs, "bucket can't empty") {
		t.Fail()
	}
	return
}

func TestDownloadNoDomain(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: 4,
		KeyFile:     "/user",
		Bucket:      "",
		Prefix:      "hello3,hello5,hello7",
		Suffixes:    "",
		Public:      true,
		CheckHash:   true,
		Referer:     "",
		CdnDomain:   "",
		UseHttps:    true,
		BatchNum:    0,
		RecordRoot:  "",
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func(){
		test.RemoveFile(path)
		test.RemoveFile(cfg.LogFile)
	}()

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path)
	if !strings.Contains(errs, "bucket / io_host / cdn_domain one them should has value") {
		t.Fail()
	}
	return
}

func TestDocumentDocument(t *testing.T) {
	test.TestDocument("qdownload", t)
}

func createDownloadConfigFile(cfg *config.Download) (cfgPath string, err error) {
	if len(cfg.DestDir) == 0 {
		rootPath, err := test.RootPath()
		if err != nil {
			return "", errors.New("get root path error:" + err.Error())
		}
		cfg.DestDir = filepath.Join(rootPath, "download")
	} else if cfg.DestDir == "empty" {
		cfg.DestDir = ""
	}
	cfg.LogSetting = &config.LogSetting{
		LogLevel:  config.DebugKey,
		LogFile:   filepath.Join(cfg.DestDir, "log.txt"),
		LogRotate: 7,
		LogStdout: "true",
	}

	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return "", errors.New("json marshal error:" + err.Error())
	}

	cfgPath, err = test.CreateFileWithContent("download_config.txt", string(data))
	if err != nil {
		err = errors.New("create cdn config file error:" + err.Error())
	}
	return cfgPath, err
}