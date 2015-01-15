package cli

import (
	"fmt"
	"qshell"
)

type CliFunc func(cmd string, params ...string)

var accountS = qshell.Account{}
var dircacheS = qshell.DirCache{}
var listbucketS = qshell.ListBucket{}
var rsfopS = qshell.RSFop{}

func Account(cmd string, params ...string) {
	if len(params) == 0 {
		accountS.Get()
		fmt.Println(accountS.String())
	} else if len(params) == 2 {
		accessKey := params[0]
		secretKey := params[1]
		accountS.Set(accessKey, secretKey)
	} else {
		CmdHelp(cmd)
	}
}
