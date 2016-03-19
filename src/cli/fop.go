package cli

import (
	"fmt"
	"qiniu/rpc"
	"qshell"
)

func Prefop(cmd string, params ...string) {
	if len(params) == 1 {
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		persistentId := params[0]
		fopRet := qshell.FopRet{}
		err := rsfopS.Prefop(persistentId, &fopRet)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Prefop error,", v.Code, v.Err)
			} else {
				fmt.Println("Prefop error,", err)
			}
		} else {
			fmt.Println(fopRet.String())
		}
	} else {
		CmdHelp(cmd)
	}
}
