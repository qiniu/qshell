package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	rpcEncodeString = "https://qiniu.com/rpc?a=1&b=1"
	rpcDecodeString = "https:!!qiniu.com!rpc'3Fa=1&b=1"
)

func TestRPCEncode(t *testing.T) {
	result := ""
	test.NewTestFlow("rpcencode", rpcEncodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, rpcDecodeString) {
		t.Fail()
	}
	return
}

func TestRPCDecode(t *testing.T) {
	result := ""
	test.NewTestFlow("rpcdecode", rpcDecodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, rpcEncodeString) {
		t.Fail()
	}
	return
}

const (
	base64EncodeString = "https://qiniu.com/rpc?a=1&b=1"
	base64DecodeString = "aHR0cHM6Ly9xaW5pdS5jb20vcnBjP2E9MSZiPTE="
)

func TestBase64Encode(t *testing.T) {
	result := ""
	test.NewTestFlow("b64encode", base64EncodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, base64DecodeString) {
		t.Fail()
	}
	return
}

func TestBase64Decode(t *testing.T) {
	result := ""
	test.NewTestFlow("b64decode", base64DecodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, base64EncodeString) {
		t.Fail()
	}
	return
}

func TestD2ts(t *testing.T) {
	duration := 0
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	result := ""
	test.NewTestFlow("d2ts", strconv.Itoa(duration)).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, timeString) {
		t.Fail()
	}
	return
}

const (
	timestamp       = 1641527120
	timestampOfDate = "2022-01-07 11:45:20"
)

func TestTs2d(t *testing.T) {
	result := ""
	test.NewTestFlow("ts2d", strconv.Itoa(timestamp)).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

func TestTms2d(t *testing.T) {
	result := ""
	test.NewTestFlow("tms2d", strconv.Itoa(timestamp*1000)).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

func TestTns2d(t *testing.T) {
	result := ""
	tns := timestamp * 1000 * 1000 * 10
	test.NewTestFlow("tns2d", strconv.Itoa(tns)).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

const (
	urlEncodeString = "https://qiniu.com/rpc?a=1&b=1"
	urlDecodeString = "https:%2F%2Fqiniu.com%2Frpc%3Fa=1&b=1"
)

func TestUrlEncode(t *testing.T) {
	result := ""
	test.NewTestFlow("urlencode", urlEncodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, urlDecodeString) {
		t.Fail()
	}
	return
}

func TestUrlDecode(t *testing.T) {
	result := ""
	test.NewTestFlow("urldecode", urlDecodeString).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, urlEncodeString) {
		t.Fail()
	}
	return
}

func TestReqid(t *testing.T) {
	result := ""
	test.NewTestFlow("reqid", "62kAAIYB06brhtsT").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(test.DefaultTestErrorHandler(t)).Run()

	if !strings.Contains(result, "2015-05-06/12-14") {
		t.Fail()
	}
	return
}
