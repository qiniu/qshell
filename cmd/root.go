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

const (
	bash_completion_func = `__qshell_parse_get()
{
    local qshell_output out
    if qshell_output=$(qshell user ls --name 2>/dev/null); then
        out=($(echo "${qshell_output}"))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__qshell_get_resource()
{
    __qshell_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        qshell_user_cu)
            __qshell_get_resource
            return
            ;;
        *)
            ;;
    esac
}
`
)

var RootCmd = &cobra.Command{
	Use:                    "qshell",
	Short:                  "Qiniu commandline tool for managing your bucket and CDN",
	Version:                version,
	BashCompletionFunction: bash_completion_func,
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
		viper.SetConfigName(".qshell")
	}
	if rErr := viper.ReadInConfig(); rErr != nil {
		if _, ok := rErr.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "read config file: %v\n", rErr)
		}
	}
}
