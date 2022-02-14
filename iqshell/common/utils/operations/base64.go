package operations

import (
	"encoding/base64"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type Base64Info struct {
	Data    string
	UrlSafe bool
}

func (info *Base64Info) Check() error {
	if len(info.Data) == 0 {
		return alert.CannotEmptyError("Data", "")
	}
	return nil
}

// Base64Encode base64编码数据
func Base64Encode(info Base64Info) {
	dataEncoded := ""
	if info.UrlSafe {
		dataEncoded = base64.URLEncoding.EncodeToString([]byte(info.Data))
	} else {
		dataEncoded = base64.StdEncoding.EncodeToString([]byte(info.Data))
	}
	log.Alert(dataEncoded)
}

// Base64Decode 解码base64编码的数据
func Base64Decode(info Base64Info) {
	if info.UrlSafe {
		dataDecoded, err := base64.URLEncoding.DecodeString(info.Data)
		if err != nil {
			log.Error("Failed to decode `", info.Data, "' in url safe mode.")
			return
		}
		log.Alert(string(dataDecoded))
	} else {
		dataDecoded, err := base64.StdEncoding.DecodeString(info.Data)
		if err != nil {
			log.Error("Failed to decode `", info.Data, "' in standard mode.")
			return
		}
		log.Alert(string(dataDecoded))
	}
}
