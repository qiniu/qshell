package cmd

import (
	"encoding/json"
	"errors"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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
		ThreadCount: data.NewInt(4),
		KeyFile:     data.NewString(keysFilePath),
		Bucket:      data.NewString(test.Bucket),
		Prefix:      data.NewString("hell"),
		Suffixes:    data.NewString(".json"),
		IoHost:      data.NewString(""),
		Public:      data.NewBool(true),
		CheckHash:   data.NewBool(true),
		Referer:     data.NewString(""),
		CdnDomain:   data.NewString(""),
		UseHttps:    data.NewBool(true),
		BatchNum:    data.NewInt(0),
		RecordRoot:  data.NewString(""),
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(path)

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir.Value()) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir.Value())
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadFromBucket(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: data.NewInt(4),
		KeyFile:     data.NewString(""),
		Bucket:      data.NewString(test.Bucket),
		Prefix:      data.NewString("hello3,hello5,hello7"),
		Suffixes:    data.NewString(""),
		IoHost:      data.NewString(test.BucketDomain),
		Public:      data.NewBool(true),
		CheckHash:   data.NewBool(true),
		Referer:     data.NewString(""),
		CdnDomain:   data.NewString(""),
		UseHttps:    data.NewBool(true),
		BatchNum:    data.NewInt(0),
		RecordRoot:  data.NewString(""),
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir.Value())
		test.RemoveFile(cfg.LogFile.Value())
	}()

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir.Value()) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir.Value())
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadNoBucket(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: data.NewInt(4),
		KeyFile:     data.NewString(""),
		Bucket:      data.NewString(""),
		Prefix:      data.NewString("hello3,hello5,hello7"),
		Suffixes:    data.NewString(""),
		IoHost:      data.NewString(test.BucketDomain),
		Public:      data.NewBool(true),
		CheckHash:   data.NewBool(true),
		Referer:     data.NewString(""),
		CdnDomain:   data.NewString(""),
		UseHttps:    data.NewBool(true),
		BatchNum:    data.NewInt(0),
		RecordRoot:  data.NewString(""),
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir.Value())
		test.RemoveFile(cfg.LogFile.Value())
	}()

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path)
	if !strings.Contains(errs, "bucket can't empty") {
		t.Fail()
	}
	return
}

func TestDownloadNoDomain(t *testing.T) {
	cfg := &config.Download{
		ThreadCount: data.NewInt(4),
		Bucket:      data.NewString(test.Bucket),
		Prefix:      data.NewString("hello3,hello5,hello7"),
		Suffixes:    data.NewString(""),
		Public:      data.NewBool(true),
		CheckHash:   data.NewBool(true),
		Referer:     data.NewString(""),
		CdnDomain:   data.NewString(""),
		UseHttps:    data.NewBool(true),
		BatchNum:    data.NewInt(0),
		RecordRoot:  data.NewString(""),
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir.Value())
		test.RemoveFile(cfg.LogFile.Value())
	}()

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir.Value()) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir.Value())
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDocumentDocument(t *testing.T) {
	test.TestDocument("qdownload", t)
}

func createDownloadConfigFile(cfg *config.Download) (cfgPath string, err error) {
	if data.Empty(cfg.DestDir) {
		rootPath, err := test.RootPath()
		if err != nil {
			return "", errors.New("get root path error:" + err.Error())
		}
		cfg.DestDir = data.NewString(filepath.Join(rootPath, "download"))
	} else if cfg.DestDir.Value() == "empty" {
		cfg.DestDir = data.NewString("")
	}
	cfg.LogSetting = &config.LogSetting{
		LogLevel:  data.NewString(config.DebugKey),
		LogFile:   data.NewString(filepath.Join(cfg.DestDir.Value(), "log.txt")),
		LogRotate: data.NewInt(7),
		LogStdout: data.NewBool(true),
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
