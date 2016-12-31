package cli

import (
	"fmt"
	"qshell"
)

var ForceMode bool
var ListMarker string

type CliFunc func(cmd string, params ...string)

var accountS = qshell.Account{}

func Account(cmd string, params ...string) {
	if len(params) == 0 {
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		fmt.Println(accountS.String())
	} else if len(params) == 2 || len(params) == 3 {
		accessKey := params[0]
		secretKey := params[1]
		sErr := accountS.Set(accessKey, secretKey)
		if sErr != nil {
			fmt.Println(sErr)
		}
	} else {
		CmdHelp(cmd)
	}
}
