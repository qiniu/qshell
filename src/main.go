package main

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"qshell"
)

var debugMode = false

const (
	CMD_ACCOUNT    = "account"
	CMD_DIRCACHE   = "dircache"
	CMD_LISTBUCKET = "listbucket"
	CMD_PREFOP     = "prefop"
	CMD_STAT       = "stat"
)

var accountS = qshell.Account{}
var dircacheS = qshell.DirCache{}
var listbucketS = qshell.ListBucket{}
var rsfopS = qshell.RSFop{}

func help(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(qshell.CmdHelpList())
	} else {
		for _, cmd := range cmds {
			fmt.Println(qshell.CmdHelp(cmd))
		}
	}
}

func main() {
	args := os.Args
	argc := len(args)
	log.SetOutputLevel(log.Linfo)
	if argc > 1 {
		cmd := ""
		params := []string{}
		option := args[1]
		if option == "-d" {
			if argc > 2 {
				cmd = args[2]
				if argc > 3 {
					params = args[3:]
				}
			}
			log.SetOutputLevel(log.Ldebug)
		} else {
			cmd = args[1]
			if argc > 2 {
				params = args[2:]
			}
		}

		switch cmd {
		case CMD_ACCOUNT:
			account(params...)
		case CMD_DIRCACHE:
			dircache(params...)
		case CMD_LISTBUCKET:
			listbucket(params...)
		case CMD_PREFOP:
			prefop(params...)
		case CMD_STAT:
			stat(params...)
		default:
			help()
		}
	} else {
		help()
	}
}

func account(params ...string) {
	if len(params) == 0 {
		accountS.Get()
		fmt.Println(accountS.String())
	} else if len(params) == 2 {
		accessKey := params[0]
		secretKey := params[1]
		accountS.Set(accessKey, secretKey)
	} else {
		help(CMD_ACCOUNT)
	}
}

func dircache(params ...string) {
	if len(params) == 2 {
		cacheRootPath := params[0]
		cacheResultFile := params[1]
		dircacheS.Cache(cacheRootPath, cacheResultFile)
	} else {
		help(CMD_DIRCACHE)
	}
}

func listbucket(params ...string) {
	if len(params) == 3 {
		bucket := params[0]
		prefix := params[1]
		listResultFile := params[2]
		//get ak,sk
		accountS.Get()
		if accountS.AccessKey != "" && accountS.SecretKey != "" {
			listbucketS.Account = accountS
			listbucketS.List(bucket, prefix, listResultFile)
		} else {
			log.Error("No AccessKey and SecretKey set error!")
		}
	} else {
		help(CMD_LISTBUCKET)
	}
}

func prefop(params ...string) {
	if len(params) == 1 {
		persistentId := params[0]
		accountS.Get()
		fopRet := qshell.FopRet{}
		err := rsfopS.Prefop(persistentId, &fopRet)
		if err != nil {
			log.Error(fmt.Sprintf("Can not get fop status for `%s'", persistentId), err)
		} else {
			fmt.Println(fopRet.String())
		}
	} else {
		help(CMD_PREFOP)
	}
}

func stat(params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		entry, err := client.Stat(nil, bucket, key)
		if err != nil {
			log.Error("Stat error,", err)
		} else {
			qshell.PrintStat(bucket, key, entry)
		}
	} else {
		help(CMD_STAT)
	}
}
