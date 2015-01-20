package cli

import (
	"github.com/qiniu/log"
	"qshell"
	"strconv"
)

func QiniuDownload(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		var threadCount int64 = 5
		var downConfig string
		var err error
		if len(params) == 1 {
			downConfig = params[0]
		} else {
			threadCount, err = strconv.ParseInt(params[0], 10, 64)
			if err != nil {
				log.Error("Invalid value for <ThreadCount>", params[0])
				return
			}
			downConfig = params[1]
		}
		qshell.QiniuDownload(int(threadCount), downConfig)
	} else {
		CmdHelp(cmd)
	}
}
