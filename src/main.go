package main

import (
	"cli"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"os/user"
	"path/filepath"
	"qiniu/api.v6/conf"
	"qiniu/rpc"
	"qshell"
	"runtime"
	"strings"
)

var supportedCmds = map[string]cli.CliFunc{
	"account":       cli.Account,
	"dircache":      cli.DirCache,
	"listbucket":    cli.ListBucket,
	"listbucket2":   cli.ListBucket2,
	"alilistbucket": cli.AliListBucket,
	"prefop":        cli.Prefop,
	"stat":          cli.Stat,
	"delete":        cli.Delete,
	"move":          cli.Move,
	"copy":          cli.Copy,
	"chgm":          cli.Chgm,
	"chtype":        cli.Chtype,
	"expire":        cli.DeleteAfterDays,
	"sync":          cli.Sync,
	"fetch":         cli.Fetch,
	"prefetch":      cli.Prefetch,
	"batchstat":     cli.BatchStat,
	"batchdelete":   cli.BatchDelete,
	"batchchgm":     cli.BatchChgm,
	"batchchtype":   cli.BatchChtype,
	"batchexpire":   cli.BatchDeleteAfterDays,
	"batchrename":   cli.BatchRename,
	"batchcopy":     cli.BatchCopy,
	"batchmove":     cli.BatchMove,
	"batchsign":     cli.BatchSign,
	"fput":          cli.FormPut,
	"rput":          cli.ResumablePut,
	"get":           cli.GetFileFromBucket,
	"qupload":       cli.QiniuUpload,
	"qupload2":      cli.QiniuUpload2,
	"qdownload":     cli.QiniuDownload,
	"b64encode":     cli.Base64Encode,
	"b64decode":     cli.Base64Decode,
	"urlencode":     cli.Urlencode,
	"urldecode":     cli.Urldecode,
	"rpcencode":     cli.RpcEncode,
	"rpcdecode":     cli.RpcDecode,
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

type HostConfig struct {
	UpHost  string `json:"up"`
	ApiHost string `json:"api"`
	IoHost  string `json:"io"`
	RsHost  string `json:"rs"`
	RsfHost string `json:"rsf"`
}

func main() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	rpc.UserAgent = cli.UserAgent()

	//parse command
	logs.SetLevel(logs.LevelInformational)
	logs.SetLogger(logs.AdapterConsole)

	if len(os.Args) <= 1 {
		fmt.Println("Use help or help [cmd1 [cmd2 [cmd3 ...]]] to see supported commands.")
		os.Exit(qshell.STATUS_HALT)
	}

	//global options
	var debugMode bool
	var helpMode bool
	var versionMode bool
	var multiUserMode bool
	var hostFile string

	//account file
	var accountFile string

	flag.BoolVar(&debugMode, "d", false, "debug mode")
	flag.BoolVar(&multiUserMode, "m", false, "multi user mode")
	flag.BoolVar(&helpMode, "h", false, "show help")
	flag.BoolVar(&versionMode, "v", false, "show version")
	flag.StringVar(&hostFile, "f", "", "host file")
	flag.StringVar(&accountFile, "c", "", "account file path")

	flag.Parse()

	if helpMode {
		cli.Help("help")
		return
	}

	if versionMode {
		cli.Version()
		return
	}

	//set log level
	if debugMode {
		logs.SetLevel(logs.LevelDebug)
	}

	//set qshell root path
	if multiUserMode {
		logs.Debug("Entering multiple user mode")
		pwd, gErr := os.Getwd()
		if gErr != nil {
			fmt.Println("Error: get current work dir error,", gErr)
			os.Exit(qshell.STATUS_HALT)
		}
		qshell.QShellRootPath = pwd
	} else {
		logs.Debug("Entering single user mode")
		curUser, gErr := user.Current()
		if gErr != nil {
			fmt.Println("Error: get current user error,", gErr)
			os.Exit(qshell.STATUS_HALT)
		}

		qshell.QShellRootPath = curUser.HomeDir
		//check account file mode
		if accountFile != "" {
			accountName := strings.TrimSuffix(filepath.Base(accountFile), filepath.Ext(accountFile))
			qshell.QAccountName = accountName
			qshell.QAccountFile = accountFile
		}
	}

	//read host file
	if hostFile != "" {
		hostFp, openErr := os.Open(hostFile)
		if openErr != nil {
			fmt.Println("Error: open specified host file error,", openErr)
			os.Exit(qshell.STATUS_HALT)
			return
		}

		var hostCfg HostConfig
		decoder := json.NewDecoder(hostFp)
		decodeErr := decoder.Decode(&hostCfg)
		if decodeErr != nil {
			fmt.Println("Error: read host file error,", decodeErr)
			os.Exit(qshell.STATUS_HALT)
			return
		}

		conf.UP_HOST = hostCfg.UpHost
		conf.RS_HOST = hostCfg.RsHost
		conf.RSF_HOST = hostCfg.RsfHost
		conf.IO_HOST = hostCfg.IoHost
		conf.API_HOST = hostCfg.ApiHost

		//bucket domains
		qshell.BUCKET_RS_HOST = hostCfg.RsHost
		qshell.BUCKET_API_HOST = hostCfg.ApiHost

		cli.IsHostFileSpecified = true
	}

	//set cmd and params
	args := flag.Args()
	if len(args) >= 1 {
		cmd := args[0]
		params := args[1:]

		if cliFunc, ok := supportedCmds[cmd]; ok {
			cliFunc(cmd, params...)
		} else {
			fmt.Printf("Error: unknown cmd `%s`\n", cmd)
			os.Exit(qshell.STATUS_HALT)
		}
	} else {
		fmt.Println("Error: sub cmd is required")
		os.Exit(qshell.STATUS_HALT)
	}
}
