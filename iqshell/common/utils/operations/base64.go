package operations

import (
	"encoding/base64"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type Base64Info struct {
	Data    string
	UrlSafe bool
}

func (info *Base64Info) Check() *data.CodeError {
	if len(info.Data) == 0 {
		return alert.CannotEmptyError("Data", "")
	}
	return nil
}

// Base64Encode base64编码数据
func Base64Encode(cfg *iqshell.Config, info Base64Info) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	dataEncoded := ""
	if info.UrlSafe {
		log.DebugF("Url safe")
		dataEncoded = base64.URLEncoding.EncodeToString([]byte(info.Data))
	} else {
		log.DebugF("No url safe")
		dataEncoded = base64.StdEncoding.EncodeToString([]byte(info.Data))
	}
	log.Alert(dataEncoded)
}

// Base64Decode 解码base64编码的数据
func Base64Decode(cfg *iqshell.Config, info Base64Info) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if info.UrlSafe {
		dataDecoded, err := base64.URLEncoding.DecodeString(info.Data)
		if err != nil {
			data.SetCmdStatusError()
			log.Error("Failed to decode `", info.Data, "' in url safe mode.")
			return
		}
		log.Alert(string(dataDecoded))
	} else {
		dataDecoded, err := base64.StdEncoding.DecodeString(info.Data)
		if err != nil {
			data.SetCmdStatusError()
			log.Error("Failed to decode `", info.Data, "' in standard mode.")
			return
		}
		log.Alert(string(dataDecoded))
	}
}
