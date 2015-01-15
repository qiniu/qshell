package cli

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/log"
	"strconv"
)

func Base64Encode(cmd string, params ...string) {
	if len(params) == 2 {
		urlSafe, err := strconv.ParseBool(params[0])
		if err != nil {
			log.Error("Invalid bool value or <UrlSafe>, must true or false")
			return
		}
		dataToEncode := params[1]
		dataEncoded := ""
		if urlSafe {
			dataEncoded = base64.URLEncoding.EncodeToString([]byte(dataToEncode))
		} else {
			dataEncoded = base64.StdEncoding.EncodeToString([]byte(dataToEncode))
		}
		fmt.Println(dataEncoded)
	} else {
		CmdHelp(cmd)
	}
}
func Base64Decode(cmd string, params ...string) {
	if len(params) == 2 {
		urlSafe, err := strconv.ParseBool(params[0])
		if err != nil {
			log.Error("Invalid bool value or <UrlSafe>, must true or false")
			return
		}
		dataToDecode := params[1]
		var dataDecoded []byte
		if urlSafe {
			dataDecoded, err = base64.URLEncoding.DecodeString(dataToDecode)
			if err != nil {
				log.Error("Failed to decode `", dataToDecode, "' in url safe mode.")
				return
			}
		} else {
			dataDecoded, err = base64.StdEncoding.DecodeString(dataToDecode)
			if err != nil {
				log.Error("Failed to decode `", dataToDecode, "' in standard mode.")
				return
			}
		}
		fmt.Println(string(dataDecoded))
	} else {
		CmdHelp(cmd)
	}
}
