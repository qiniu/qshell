//go:build integration

package cmd

import (
	"encoding/json"
	"errors"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download/operations"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadWithKeyFile(t *testing.T) {
	keys := test.KeysString + "\nhello_10.json"
	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", keys)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(keysFilePath)

	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
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
			RecordRoot:  "",
		},
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(path)

	test.RunCmdWithError("qdownload", "-c", "4", path, "-d")
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir)
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadFromBucket(t *testing.T) {
	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			ThreadCount:     4,
			KeyFile:         "",
			SavePathHandler: "{{pathJoin .DestDir (replace \"hello\" \"lala\" .Key)}}",
			Bucket:          test.Bucket,
			Prefix:          "hello3,hello5,hello7",
			Suffixes:        "",
			IoHost:          test.BucketDomain,
			Public:          true,
			CheckHash:       true,
			Referer:         "",
			CdnDomain:       "",
			RecordRoot:      "",
		},
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir)
		test.RemoveFile(cfg.LogFile.Value())
	}()

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}
}

func TestDownloadNoBucket(t *testing.T) {
	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
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
			RecordRoot:  "",
		},
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir)
		test.RemoveFile(cfg.LogFile.Value())
	}()

	_, errs := test.RunCmdWithError("qdownload", "-c", "4", path)
	if !strings.Contains(errs, "bucket can't empty") {
		t.Fail()
	}
	return
}

func TestDownloadNoDomain(t *testing.T) {
	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			ThreadCount: 4,
			Bucket:      test.Bucket,
			Prefix:      "hello3,hello5,hello7",
			Suffixes:    "",
			Public:      true,
			CheckHash:   true,
			Referer:     "",
			CdnDomain:   "",
			RecordRoot:  "",
		},
	}
	path, err := createDownloadConfigFile(cfg)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer func() {
		test.RemoveFile(path)
		test.RemoveFile(cfg.DestDir)
		test.RemoveFile(cfg.LogFile.Value())
	}()

	test.RunCmdWithError("qdownload", "-c", "4", path)
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}

	err = test.RemoveFile(cfg.DestDir)
	if err != nil {
		t.Log("remove file error:", err)
	}
	return
}

func TestDownloadDocument(t *testing.T) {
	test.TestDocument("qdownload", t)
}

func createDownloadConfigFile(cfg *DownloadCfg) (cfgPath string, err error) {
	if len(cfg.DestDir) == 0 {
		rootPath, err := test.RootPath()
		if err != nil {
			return "", data.NewEmptyError().AppendDesc("get root path error:" + err.Error())
		}
		cfg.DestDir = filepath.Join(rootPath, "download")
	}
	cfg.LogLevel = data.NewString(config.DebugKey)
	cfg.LogFile = data.NewString(filepath.Join(cfg.DestDir, "log.txt"))
	cfg.LogRotate = data.NewInt(7)
	cfg.LogStdout = data.NewBool(true)

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

type DownloadCfg struct {
	operations.DownloadCfg
	LogLevel  *data.String `json:"log_level"`
	LogFile   *data.String `json:"log_file"`
	LogRotate *data.Int    `json:"log_rotate"`
	LogStdout *data.Bool   `json:"log_stdout"`
}
