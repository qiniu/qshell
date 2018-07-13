package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/conf"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"github.com/tonycai653/iqshell/qshell"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

type HostConfig struct {
	UpHost  string `json:"up"`
	ApiHost string `json:"api"`
	IoHost  string `json:"io"`
	RsHost  string `json:"rs"`
	RsfHost string `json:"rsf"`
}

var (
	DebugFlag   bool
	VersionFlag bool
	HostFile    string
	AccountFile string
)

var RootCmd = &cobra.Command{
	Use:     "qshell",
	Short:   "Qiniu commandline tool for managing your bucket and CDN",
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//set log level
		if DebugFlag {
			logs.SetLevel(logs.LevelDebug)
		}

		//check account file mode
		if AccountFile != "" {
			accountName := strings.TrimSuffix(filepath.Base(AccountFile), filepath.Ext(AccountFile))
			qshell.QAccountName = accountName
			qshell.QAccountFile = AccountFile
		}

		//read host file
		if HostFile != "" {
			hostFp, openErr := os.Open(HostFile)
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
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
	RootCmd.PersistentFlags().StringVarP(&HostFile, "hostfile", "f", "", "host file")
	RootCmd.PersistentFlags().StringVarP(&AccountFile, "account", "a", "", "account file path")
}

func initConfig() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	rpc.UserAgent = UserAgent()

	//parse command
	logs.SetLevel(logs.LevelInformational)
	logs.SetLogger(logs.AdapterConsole)

	curUser, gErr := user.Current()
	if gErr != nil {
		fmt.Println("Error: get current user error,", gErr)
		os.Exit(qshell.STATUS_HALT)
	}
	qshell.QShellRootPath = curUser.HomeDir
}
