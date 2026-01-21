package operations

import (
	"bufio"
	"os"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type RpcInfo struct {
	Params []string
}

func (info *RpcInfo) Check() *data.CodeError {
	return nil
}

func RpcDecode(cfg *iqshell.Config, info RpcInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if len(info.Params) > 0 {
		for _, param := range info.Params {
			decodedStr, _ := utils.Decode(param)
			log.Alert(decodedStr)
		}
	} else {
		bScanner := bufio.NewScanner(os.Stdin)
		for bScanner.Scan() {
			toDecode := bScanner.Text()
			decodedStr, _ := utils.Decode(toDecode)
			log.Alert(decodedStr)
		}
	}
}

func RpcEncode(cfg *iqshell.Config, info RpcInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	if len(info.Params) == 0 {
		data.SetCmdStatusError()
		log.Error(alert.CannotEmpty("Data", ""))
		return
	}

	for _, param := range info.Params {
		encodedStr := utils.Encode(param)
		log.Alert(encodedStr)
	}
}
