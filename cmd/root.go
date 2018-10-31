package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"runtime"
)

var (
	DebugFlag   bool
	VersionFlag bool
	cfgFile     string
)

var RootCmd = &cobra.Command{
	Use:     "qshell",
	Short:   "Qiniu commandline tool for managing your bucket and CDN",
	Version: version,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "config file (default is $HOME/.qshell.json)")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	storage.UserAgent = UserAgent()

	//parse command
	logs.SetLevel(logs.LevelInformational)
	logs.SetLogger(logs.AdapterConsole)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		curUser, gErr := user.Current()
		if gErr != nil {
			fmt.Fprintf(os.Stderr, "get current user: %v\n", gErr)
			os.Exit(1)
		}
		viper.AddConfigPath(curUser.HomeDir)
		viper.SetConfigName(".qshell.json")
	}
	viper.ReadInConfig()
}
