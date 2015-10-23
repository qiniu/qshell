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
	} else if len(params) == 2 || len(params) == 3 {
		accessKey := params[0]
		secretKey := params[1]
		zone := qshell.ZoneNB
		if len(params) == 3 {
			val := params[2]
			switch val {
			case qshell.ZoneNB, qshell.ZoneBC, qshell.ZoneAWS:
				zone = val
			default:
				fmt.Println(fmt.Sprintf("invalid zone '%s'", zone))
				return
			}
		}
		accountS.Set(accessKey, secretKey, zone)
	} else {
		CmdHelp(cmd)
	}
}
