package cli

import (
	"fmt"
	"github.com/qiniu/log"
	"qshell"
)

func Prefop(cmd string, params ...string) {
	if len(params) == 1 {
		persistentId := params[0]
		accountS.Get()
		fopRet := qshell.FopRet{}
		err := rsfopS.Prefop(persistentId, &fopRet)
		if err != nil {
			log.Error(fmt.Sprintf("Can not get fop status for `%s',", persistentId), err)
		} else {
			fmt.Println(fopRet.String())
		}
	} else {
		CmdHelp(cmd)
	}
}
