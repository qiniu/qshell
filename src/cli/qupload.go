package cli

import (
	"github.com/qiniu/log"
	"qshell"
	"strconv"
)

func QiniuUpload(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		putThresold := qshell.PUT_THRESHOLD
		var uploadConfigFile string
		var err error
		if len(params) == 2 {
			putThresold, err = strconv.ParseInt(params[0], 10, 64)
			if err != nil {
				log.Error("Invalid <PutThresholdInBytes> value,", params[0])
				return
			}
			uploadConfigFile = params[1]
		} else {
			uploadConfigFile = params[0]
		}
		qshell.QiniuUpload(putThresold, uploadConfigFile)
	} else {
		CmdHelp(cmd)
	}
}
