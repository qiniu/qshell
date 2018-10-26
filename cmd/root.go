package cmd

import (
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"runtime"
)

var (
	DebugFlag   bool
	VersionFlag bool
)

var RootCmd = &cobra.Command{
	Use:     "qshell",
	Short:   "Qiniu commandline tool for managing your bucket and CDN",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//set log level
		if DebugFlag {
			logs.SetLevel(logs.LevelDebug)
		}
		err := qshell.ReadConfigFile()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
}

func initConfig() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	storage.UserAgent = UserAgent()

	//parse command
	logs.SetLevel(logs.LevelInformational)
	logs.SetLogger(logs.AdapterConsole)
}
