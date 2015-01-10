package main

import (
	"cli"
	"github.com/qiniu/log"
	"os"
)

var debugMode = false

var supportCmds = map[string]cli.CliFunc{
	"account":    cli.Account,
	"dircache":   cli.DirCache,
	"listbucket": cli.ListBucket,
	"prefop":     cli.Prefop,
	"stat":       cli.Stat,
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
		hit := false
		for cmdName, cliFunc := range supportCmds {
			if cmdName == cmd {
				cliFunc(cmd, params...)
				hit = true
				break
			}
		}
		if !hit {
			cli.Help()
		}
	} else {
		cli.Help()
	}
}
