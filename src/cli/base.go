package cli

import (
	"fmt"
	"os"
	"qshell/qshell"
)

type CliFunc func(cmd string, params ...string)

func Account(cmd string, params ...string) {
	if len(params) == 0 {
		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}
		fmt.Println(account.String())
	} else if len(params) == 2 || len(params) == 3 {
		accessKey := params[0]
		secretKey := params[1]
		sErr := qshell.SetAccount(accessKey, secretKey)
		if sErr != nil {
			fmt.Println(sErr)
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}
