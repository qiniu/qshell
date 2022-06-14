//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
	"testing"
)

func TestSyncV1(t *testing.T) {
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	result, errs := test.RunCmdWithError("sync", url, test.Bucket,
		"-k", "1024K.tmp",
		"--overwrite",
		"-d")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "Key: 1024K.tmp") {
		t.Fatal(result)
	}
}

func TestSyncV1WithKey(t *testing.T) {
	key := "sync_v1_1024K.mp4"
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	result, errs := test.RunCmdWithError("sync", url, test.Bucket,
		"-k", key)
	if len(errs) > 0 {
		t.Fail()
	}

	test.RunCmdWithError("delete", test.Bucket, key)

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "Key: "+key) {
		t.Fatal(result)
	}
}

func TestSyncV2(t *testing.T) {
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	result, errs := test.RunCmdWithError("sync", url, test.Bucket,
		"-k", "1024K.tmp",
		"--overwrite")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestSyncV2WithKey(t *testing.T) {
	key := "sync_v2_key.mp4"
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"

	test.RunCmdWithError("delete", test.Bucket, key)

	result, errs := test.RunCmdWithError("sync", url, test.Bucket, "--resumable-api-v2", "--key", key, "-d")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}

	if !strings.Contains(result, "Key: "+key) {
		t.Fatal(result)
	}
}

func TestSyncWithUploadHost(t *testing.T) {
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	result, errs := test.RunCmdWithError("sync", url, test.Bucket,
		"-k", "1024K.tmp",
		"--up-host", "up-na0.qiniup.com",
		"--overwrite")
	if len(errs) > 0 {
		t.Fail()
	}

	result = strings.ReplaceAll(result, "\n", "")
	if !strings.Contains(result, "Upload File success") {
		t.Fatal(result)
	}
}

func TestSyncWithWrongUploadHost(t *testing.T) {
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	_, errs := test.RunCmdWithError("sync", url, test.Bucket,
		"-k", "1024K.tmp",
		"--up-host", "up-mock.qiniup.com")
	if !strings.Contains(errs, "dial tcp: lookup up-mock.qiniup.com: no such host") &&
		!strings.Contains(errs, "Sync file error") {
		t.Fail()
	}
}

func TestSyncNoUrl(t *testing.T) {
	_, errs := test.RunCmdWithError("sync")
	if !strings.Contains(errs, "SrcResUrl can't empty") {
		t.Fail()
	}
}

func TestSyncNoBucket(t *testing.T) {
	url := "https://qshell-na0.qiniupkg.com/1024K.tmp"
	_, errs := test.RunCmdWithError("sync", url)
	if !strings.Contains(errs, "Bucket can't empty") {
		t.Fail()
	}
}

func TestSyncDocument(t *testing.T) {
	test.TestDocument("sync", t)
}
