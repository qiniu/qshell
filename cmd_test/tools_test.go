//go:build unit

package cmd

import (
	"fmt"
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
	result := test.RunCmd(t, "rpcencode", rpcEncodeString)
	if !strings.Contains(result, rpcDecodeString) {
		t.Fail()
	}
	return
}

func TestRPCEncodeMoreData(t *testing.T) {
	result := test.RunCmd(t, "rpcencode", rpcEncodeString, rpcEncodeString)
	if !strings.Contains(result, rpcDecodeString) {
		t.Fail()
	}
	return
}

func TestRPCEncodeNoData(t *testing.T) {
	_, err := test.RunCmdWithError("rpcencode")
	if !strings.Contains(err, "Data can't empty") {
		t.Fail()
	}
	return
}

func TestRPCEncodeDocument(t *testing.T) {
	test.TestDocument("rpcencode", t)
}

func TestRPCDecode(t *testing.T) {
	result := test.RunCmd(t, "rpcdecode", rpcDecodeString)
	if !strings.Contains(result, rpcEncodeString) {
		t.Fail()
	}
	return
}

func TestRPCDecodeMoreData(t *testing.T) {
	result, _ := test.RunCmdWithError("rpcdecode", rpcDecodeString, rpcDecodeString)
	if !strings.Contains(result, rpcEncodeString) {
		t.Fail()
	}
	return
}

func TestRPCDecodeDocument(t *testing.T) {
	test.TestDocument("rpcdecode", t)
}

const (
	base64EncodeString = "https://qiniu.com/rpc?a=1&b=1"
	base64DecodeString = "aHR0cHM6Ly9xaW5pdS5jb20vcnBjP2E9MSZiPTE="
)

func TestBase64Encode(t *testing.T) {
	result := test.RunCmd(t, "b64encode", base64EncodeString)
	if !strings.Contains(result, base64DecodeString) {
		t.Fail()
	}
	return
}

func TestBase64EncodeNoData(t *testing.T) {
	_, err := test.RunCmdWithError("b64encode")
	if !strings.Contains(err, "Data can't empty") {
		t.Fail()
	}
	return
}

func TestB64EncodeDocument(t *testing.T) {
	test.TestDocument("b64encode", t)
}

func TestBase64Decode(t *testing.T) {
	result := test.RunCmd(t, "b64decode", base64DecodeString)
	if !strings.Contains(result, base64EncodeString) {
		t.Fail()
	}
	return
}

func TestBase64DecodeNoData(t *testing.T) {
	_, err := test.RunCmdWithError("b64decode")
	if !strings.Contains(err, "Data can't empty") {
		t.Fail()
	}
	return
}

func TestB64DecodeDocument(t *testing.T) {
	test.TestDocument("b64decode", t)
}

func TestD2ts(t *testing.T) {
	duration := 0
	currentTime := time.Now()
	timestampString := fmt.Sprintf("%d", currentTime.Unix())
	result := test.RunCmd(t, "d2ts", strconv.Itoa(duration))
	if !strings.Contains(result, timestampString) {
		t.Fail()
	}
	return
}

func TestD2tsNoData(t *testing.T) {
	_, err := test.RunCmdWithError("d2ts")
	if !strings.Contains(err, "args can't empty") {
		t.Fail()
	}
	return
}

func TestD2TsDocument(t *testing.T) {
	test.TestDocument("d2ts", t)
}

const (
	timestamp       = 1641527120
	timestampOfDate = "2022-01-07 11:45:20"
)

func TestTs2D(t *testing.T) {
	result := test.RunCmd(t, "ts2d", strconv.Itoa(timestamp))
	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

func TestTs2DNoData(t *testing.T) {
	_, err := test.RunCmdWithError("ts2d")
	if !strings.Contains(err, "args can't empty") {
		t.Fail()
	}
	return
}

func TestTS2dDocument(t *testing.T) {
	test.TestDocument("ts2d", t)
}

func TestTms2d(t *testing.T) {
	result := test.RunCmd(t, "tms2d", strconv.Itoa(timestamp*1000))
	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

func TestTms2dNoData(t *testing.T) {
	_, err := test.RunCmdWithError("tms2d")
	if !strings.Contains(err, "args can't empty") {
		t.Fail()
	}
	return
}

func TestTMs2dDocument(t *testing.T) {
	test.TestDocument("tms2d", t)
}

func TestTns2d(t *testing.T) {
	tns := timestamp * 1000 * 1000 * 10
	result := test.RunCmd(t, "tns2d", strconv.Itoa(tns))
	if !strings.Contains(result, timestampOfDate) {
		t.Fail()
	}
	return
}

func TestTns2dNoData(t *testing.T) {
	_, err := test.RunCmdWithError("tns2d")
	if !strings.Contains(err, "args can't empty") {
		t.Fail()
	}
	return
}

func TestTNs2dDocument(t *testing.T) {
	test.TestDocument("tns2d", t)
}

const (
	urlEncodeString = "https://qiniu.com/rpc?a=1&b=1"
	urlDecodeString = "https:%2F%2Fqiniu.com%2Frpc%3Fa=1&b=1"
)

func TestUrlEncode(t *testing.T) {
	result := test.RunCmd(t, "urlencode", urlEncodeString)
	if !strings.Contains(result, urlDecodeString) {
		t.Fail()
	}
	return
}

func TestUrlEncodeNoData(t *testing.T) {
	_, err := test.RunCmdWithError("urlencode")
	if !strings.Contains(err, "Data can't empty") {
		t.Fail()
	}
	return
}

func TestUrlEncodeDocument(t *testing.T) {
	test.TestDocument("urlencode", t)
}

func TestUrlDecode(t *testing.T) {
	result := test.RunCmd(t, "urldecode", urlDecodeString)
	if !strings.Contains(result, urlEncodeString) {
		t.Fail()
	}
	return
}

func TestUrlDecodeNoData(t *testing.T) {
	_, err := test.RunCmdWithError("urldecode")
	if !strings.Contains(err, "Data can't empty") {
		t.Fail()
	}
	return
}

func TestUrlDecodeDocument(t *testing.T) {
	test.TestDocument("urldecode", t)
}

func TestReqid(t *testing.T) {
	result := test.RunCmd(t, "reqid", "62kAAIYB06brhtsT")
	if !strings.Contains(result, "2015-05-06/12-14") {
		t.Fail()
	}
	return
}

func TestReqidNoData(t *testing.T) {
	_, err := test.RunCmdWithError("reqid")
	if !strings.Contains(err, "ReqId can't empty") {
		t.Fail()
	}
	return
}

func TestReqIdDocument(t *testing.T) {
	test.TestDocument("reqid", t)
}


func TestIP(t *testing.T) {
	result := test.RunCmd(t, "ip", "180.154.236.170")
	if !strings.Contains(result, "180.154.236.170") {
		t.Fail()
	}
	return
}

func TestIPDocument(t *testing.T) {
	test.TestDocument("ip", t)
}