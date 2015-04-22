package main

import (
	"cli"
	"fmt"
	"github.com/qiniu/log"
	"os"
)

var debugMode = false

var supportedCmds = map[string]cli.CliFunc{
	"account":       cli.Account,
	"dircache":      cli.DirCache,
	"listbucket":    cli.ListBucket,
	"alilistbucket": cli.AliListBucket,
	"prefop":        cli.Prefop,
	"stat":          cli.Stat,
	"delete":        cli.Delete,
	"move":          cli.Move,
	"copy":          cli.Copy,
	"chgm":          cli.Chgm,
	"fetch":         cli.Fetch,
	"prefetch":      cli.Prefetch,
	"batchdelete":   cli.BatchDelete,
	"batchchgm":     cli.BatchChgm,
	"batchrename":   cli.BatchRename,
	"batchmove":     cli.BatchMove,
	"checkqrsync":   cli.CheckQrsync,
	"fput":          cli.FormPut,
	"qupload":       cli.QiniuUpload,
	"qdownload":     cli.QiniuDownload,
	"rput":          cli.ResumablePut,
	"b64encode":     cli.Base64Encode,
	"b64decode":     cli.Base64Decode,
	"urlencode":     cli.Urlencode,
	"urldecode":     cli.Urldecode,
	"ts2d":          cli.Timestamp2Date,
	"tns2d":         cli.TimestampNano2Date,
	"tms2d":         cli.TimestampMilli2Date,
	"d2ts":          cli.Date2Timestamp,
	"ip":            cli.IpQuery,
	"qetag":         cli.Qetag,
	"help":          cli.Help,
	"unzip":         cli.Unzip,
	"privateurl":    cli.PrivateUrl,
	"saveas":        cli.Saveas,
	"reqid":         cli.ReqId,
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
		for cmdName, cliFunc := range supportedCmds {
			if cmdName == cmd {
				cliFunc(cmd, params...)
				hit = true
				break
			}
		}
		if !hit {
			fmt.Println(fmt.Sprintf("Unknow cmd `%s'", cmd))
		}
	} else {
		fmt.Println("Use help or help [cmd1 [cmd2 [cmd3 ...]]] to see supported commands.")
	}
}
