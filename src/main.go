package main

import (
	"cli"
	"flag"
	"fmt"
	"os"
	"qiniu/log"
	"qiniu/rpc"
	"runtime"
)

var supportedCmds = map[string]cli.CliFunc{
	"account":       cli.Account,
	"zone":          cli.Zone,
	"dircache":      cli.DirCache,
	"listbucket":    cli.ListBucket,
	"alilistbucket": cli.AliListBucket,
	"prefop":        cli.Prefop,
	"stat":          cli.Stat,
	"delete":        cli.Delete,
	"move":          cli.Move,
	"copy":          cli.Copy,
	"chgm":          cli.Chgm,
	"sync":          cli.Sync,
	"fetch":         cli.Fetch,
	"prefetch":      cli.Prefetch,
	"batchstat":     cli.BatchStat,
	"batchdelete":   cli.BatchDelete,
	"batchchgm":     cli.BatchChgm,
	"batchrename":   cli.BatchRename,
	"batchcopy":     cli.BatchCopy,
	"batchmove":     cli.BatchMove,
	"batchsign":     cli.BatchSign,
	"checkqrsync":   cli.CheckQrsync,
	"fput":          cli.FormPut,
	"qupload":       cli.QiniuUpload,
	"qupload2":      cli.QiniuUpload2,
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
	"m3u8delete":    cli.M3u8Delete,
	"m3u8replace":   cli.M3u8Replace,
	"buckets":       cli.GetBuckets,
	"domains":       cli.GetDomainsOfBucket,
	"cdnrefresh":    cli.CdnRefresh,
	"cdnprefetch":   cli.CdnPrefetch,
}

func main() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	rpc.UserAgent = cli.UserAgent()

	//parse command
	log.SetOutputLevel(log.Linfo)
	log.SetOutput(os.Stdout)

	if len(os.Args) <= 1 {
		fmt.Println("Use help or help [cmd1 [cmd2 [cmd3 ...]]] to see supported commands.")
		return
	}

	//flag
	var debugMode bool
	var forceMode bool
	var helpMode bool
	var versionMode bool

	flag.BoolVar(&debugMode, "d", false, "debug mode")
	flag.BoolVar(&forceMode, "f", false, "force mode")
	flag.BoolVar(&helpMode, "h", false, "show help")
	flag.BoolVar(&versionMode, "v", false, "show version")

	flag.Parse()

	if helpMode {
		cli.Help("help")
		return
	}

	if versionMode {
		cli.Version()
		return
	}

	if debugMode {
		log.SetOutputLevel(log.Ldebug)
	}

	if forceMode {
		cli.ForceMode = true
	}

	args := flag.Args()

	cmd := args[0]
	params := args[1:]

	if cmd == "" {
		fmt.Println("Error: no subcommand specified")
		return
	}

	if cliFunc, ok := supportedCmds[cmd]; ok {
		cliFunc(cmd, params...)
	} else {
		fmt.Println(fmt.Sprintf("Error: unknown cmd `%s`", cmd))
	}

}
