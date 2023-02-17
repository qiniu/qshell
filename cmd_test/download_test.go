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
	test.RemoveRootPath()

	keys := test.KeysString + "\nhello_10.json"
	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", keys)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}
	defer test.RemoveFile(keysFilePath)

	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			KeyFile:    keysFilePath,
			Bucket:     test.Bucket,
			Prefix:     "hell",
			Suffixes:   ".json",
			IoHost:     "",
			Public:     true,
			CheckHash:  true,
			Referer:    "",
			CdnDomain:  "",
			RecordRoot: "",
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
	test.RemoveRootPath()

	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			KeyFile:         "",
			SavePathHandler: "{{pathJoin .DestDir (replace \"hello\" \"lala\" .Key)}}",
			Bucket:          test.Bucket,
			Prefix:          "hello3,hello5,hello7",
			Suffixes:        "",
			IoHost:          test.BucketDomain,
			Public:          true,
			CheckSize:       true,
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

	test.RunCmdWithError("qdownload", "-c", "4", path, "-d")
	if test.FileCountInDir(cfg.DestDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(cfg.LogFile.Value()) {
		t.Fatal("log file should has content")
	}
}

func TestDownloadNoBucket(t *testing.T) {
	test.RemoveRootPath()

	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			KeyFile:    "",
			Bucket:     "",
			Prefix:     "hello3,hello5,hello7",
			Suffixes:   "",
			IoHost:     test.BucketDomain,
			Public:     true,
			CheckHash:  true,
			Referer:    "",
			CdnDomain:  "",
			RecordRoot: "",
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
	test.RemoveRootPath()

	cfg := &DownloadCfg{
		DownloadCfg: operations.DownloadCfg{
			Bucket:     test.Bucket,
			Prefix:     "hello3,hello5,hello7",
			Suffixes:   "",
			Public:     true,
			CheckHash:  true,
			Referer:    "",
			CdnDomain:  "",
			RecordRoot: "",
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

func TestDownload2NoBucket(t *testing.T) {
	_, errs := test.RunCmdWithError("qdownload2", "-c", "4")
	if !strings.Contains(errs, "bucket can't empty") {
		t.Fail()
	}
	return
}

func TestDownload2AllFilesFromBucket(t *testing.T) {
	test.RemoveRootPath()

	rootPath, err := test.RootPath()
	if err != nil {
		t.Fatal("get root path error:", err)
	}

	destDir := filepath.Join(rootPath, "download2")
	logPath := filepath.Join(rootPath, "download2_log")
	defer func() {
		test.RemoveFile(destDir)
		test.RemoveFile(logPath)
	}()

	test.RunCmdWithError("qdownload2",
		"--bucket", test.Bucket,
		"--dest-dir", destDir,
		"--suffixes", ".json",
		"--public",
		"--log-file", logPath,
		"--log-level", "info",
		"-c", "4")
	if test.FileCountInDir(destDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(logPath) {
		t.Fatal("log file should has content")
	}

	logContent := test.FileContent(logPath)
	if strings.Contains(logContent, "[D]") {
		t.Fatal("shouldn't has debug log")
	}

	return
}

func TestDownload2WithKeyFile(t *testing.T) {
	test.RemoveRootPath()

	keys := test.KeysString + "\nhello_10.json"
	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", keys)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	rootPath, err := test.RootPath()
	if err != nil {
		t.Fatal("get root path error:", err)
	}

	destDir := filepath.Join(rootPath, "download2")
	logPath := filepath.Join(rootPath, "download2_log.private")
	successLogPath := filepath.Join(rootPath, "download2_success.txt")
	failLogPath := filepath.Join(rootPath, "download2_fail.txt")
	defer func() {
		test.RemoveFile(keysFilePath)
		test.RemoveFile(destDir)
		test.RemoveFile(logPath)
	}()

	test.RunCmdWithError("qdownload2",
		"--bucket", test.Bucket,
		"--dest-dir", destDir,
		"--key-file", keysFilePath,
		"--log-file", logPath,
		"--log-level", "debug",
		"-s", successLogPath,
		"-f", failLogPath,
		"-c", "4",
		"-d")
	if test.FileCountInDir(destDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(logPath) {
		t.Fatal("log file should has content")
	}

	if !test.IsFileHasContent(successLogPath) {
		t.Fatal("success log file should has content")
	}

	if !test.IsFileHasContent(failLogPath) {
		t.Fatal("fail log file should has content")
	}

	logContent := test.FileContent(logPath)
	if !strings.Contains(logContent, "?e=") {
		t.Fatal("download url should private")
	}

	if !strings.Contains(logContent, "work consumer 3 start") {
		t.Fatal("download should have consumer 3")
	}

	if strings.Contains(logContent, "work consumer 4 start") {
		t.Fatal("download shouldn't have consumer 4")
	}
	return
}

func TestDownload2PublicWithKeyFile(t *testing.T) {
	test.RemoveRootPath()

	keys := test.KeysString + "\nhello_10.json"
	keysFilePath, err := test.CreateFileWithContent("download_keys.txt", keys)
	if err != nil {
		t.Fatal("create cdn config file error:", err)
	}

	rootPath, err := test.RootPath()
	if err != nil {
		t.Fatal("get root path error:", err)
	}

	destDir := filepath.Join(rootPath, "download2")
	logPath := filepath.Join(rootPath, "download2_log.public")
	defer func() {
		test.RemoveFile(keysFilePath)
		test.RemoveFile(destDir)
		test.RemoveFile(logPath)
	}()

	test.RunCmdWithError("qdownload2",
		"--bucket", test.Bucket,
		"--dest-dir", destDir,
		"--key-file", keysFilePath,
		"--log-file", logPath,
		"--log-level", "debug",
		"--public",
		"-c", "4",
		"-d")
	if test.FileCountInDir(destDir) < 2 {
		t.Fail()
	}

	if !test.IsFileHasContent(logPath) {
		t.Fatal("log file should has content")
	}

	logContent := test.FileContent(logPath)
	if strings.Contains(logContent, "?e=") {
		t.Fatal("download url should public")
	}

	return
}

func TestDownload2Document(t *testing.T) {
	test.TestDocument("qdownload2", t)
}
