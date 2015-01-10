package cli

import (
	"fmt"
	"os"
	"qshell"
)

type CliFunc func(cmd string, params ...string)

var accountS = qshell.Account{}
var dircacheS = qshell.DirCache{}
var listbucketS = qshell.ListBucket{}
var rsfopS = qshell.RSFop{}

func Help(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(CmdHelpList())
	} else {
		for _, cmd := range cmds {
			fmt.Println(CmdHelp(cmd))
		}
	}
}

func Account(cmd string, params ...string) {
	if len(params) == 0 {
		accountS.Get()
		fmt.Println(accountS.String())
	} else if len(params) == 2 {
		accessKey := params[0]
		secretKey := params[1]
		accountS.Set(accessKey, secretKey)
	} else {
		Help(cmd)
	}
}
