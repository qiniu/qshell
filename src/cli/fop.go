package cli

import (
	"fmt"
	"os"
	"qshell/qiniu/rpc"
	"qshell/qshell"
)

func Prefop(cmd string, params ...string) {
	if len(params) == 1 {
		persistentId := params[0]
		fopRet := qshell.FopRet{}
		err := qshell.Prefop(persistentId, &fopRet)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Prefop error,", v.Code, v.Err)
			} else {
				fmt.Println("Prefop error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		} else {
			fmt.Println(fopRet.String())
		}
	} else {
		CmdHelp(cmd)
	}
}
