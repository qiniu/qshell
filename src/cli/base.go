package cli

import (
	"fmt"
	"qshell"
)

var ForceMode bool

type CliFunc func(cmd string, params ...string)

var accountS = qshell.Account{}
var dircacheS = qshell.DirCache{}
var listbucketS = qshell.ListBucket{}
var rsfopS = qshell.RSFop{}

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
		var zone string
		if len(params) == 3 {
			zone = params[2]
			if !qshell.IsValidZone(zone) {
				fmt.Println(fmt.Sprintf("Invalid zone '%s'", zone))
			}
		}
		sErr := accountS.Set(accessKey, secretKey, zone)
		if sErr != nil {
			fmt.Println(sErr)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Zone(cmd string, params ...string) {
	if len(params) == 0 {
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		fmt.Println("Current zone:", accountS.Zone)
	} else if len(params) == 1 {
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		accessKey := accountS.AccessKey
		secretKey := accountS.SecretKey
		zone := params[0]
		if !qshell.IsValidZone(zone) {
			fmt.Println(fmt.Sprintf("Invalid zone '%s'", zone))
			return
		}

		sErr := accountS.Set(accessKey, secretKey, zone)
		if sErr != nil {
			fmt.Println(sErr)
		}
	} else {
		CmdHelp(cmd)
	}
}
