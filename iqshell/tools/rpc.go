package tools

import (
	"bufio"
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type RpcInfo struct {
	params []string
}

func RpcDecode(info RpcInfo) {
	if len(info.params) > 0 {
		for _, param := range info.params {
			decodedStr, _ := utils.Decode(param)
			log.Alert(decodedStr)
		}
	} else {
		bScanner := bufio.NewScanner(os.Stdin)
		for bScanner.Scan() {
			toDecode := bScanner.Text()
			decodedStr, _ := utils.Decode(string(toDecode))
			log.Alert(decodedStr)
		}
	}
}

func RpcEncode(info RpcInfo) {
	for _, param := range info.params {
		encodedStr := utils.Encode(param)
		log.Alert(encodedStr)
	}
}
