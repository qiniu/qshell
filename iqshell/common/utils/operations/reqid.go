package operations

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strconv"
	"time"
)

type ReqIdInfo struct {
	ReqId string
}

func (info *ReqIdInfo) Check() *data.CodeError {
	if len(info.ReqId) == 0 {
		return alert.CannotEmptyError("ReqId", "")
	}
	return nil
}

// DecodeReqId 解析reqid， 打印人工可读的字符串
func DecodeReqId(cfg *iqshell.Config, info ReqIdInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	decodedBytes, err := base64.URLEncoding.DecodeString(info.ReqId)
	if err != nil || len(decodedBytes) < 4 {
		data.SetCmdStatusError()
		log.Error("Invalid reqid", info.ReqId, err)
		return
	}

	newBytes := decodedBytes[4:]
	newBytesLen := len(newBytes)
	newStr := ""
	for i := newBytesLen - 1; i >= 0; i-- {
		newStr += fmt.Sprintf("%02X", newBytes[i])
	}

	unixNano, err := strconv.ParseInt(newStr, 16, 64)
	if err != nil {
		data.SetCmdStatusError()
		log.Error("Invalid reqid", info.ReqId, err)
		return
	}

	dstDate := time.Unix(0, unixNano)
	log.AlertF("%04d-%02d-%02d/%02d-%02d", dstDate.Year(), dstDate.Month(), dstDate.Day(),
		dstDate.Hour(), dstDate.Minute())
}
