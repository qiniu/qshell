package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	DebugFlag   bool
	VersionFlag bool
	cfgFile     string
	local       bool
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
	RootCmd.PersistentFlags().BoolVarP(&local, "local", "L", false, "use current directory as config file path")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("local", RootCmd.PersistentFlags().Lookup("local"))
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

	if local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		viper.Set("path.root_path", dir+"/.qshell")
	} else {
		curUser, gErr := user.Current()
		if gErr != nil {
			fmt.Fprintf(os.Stderr, "Error: get current user error: %v\n", gErr)
			os.Exit(1)
		}
		viper.Set("path.root_path", curUser.HomeDir+"/.qshell")
	}
	rootPath := viper.GetString("path.root_path")

	viper.SetDefault("path.acc_db_path", filepath.Join(rootPath, "account.db"))
	viper.SetDefault("path.acc_path", filepath.Join(rootPath, "account.json"))
	viper.SetDefault("hosts.up_host", "upload.qiniup.com")
	viper.SetDefault("hosts.rs_host", storage.DefaultRsHost)
	viper.SetDefault("hosts.rsf_host", storage.DefaultRsfHost)
	viper.SetDefault("hosts.io_host", "iovip.qbox.me")
	viper.SetDefault("hosts.api_host", storage.DefaultAPIHost)

	if rErr := viper.ReadInConfig(); rErr != nil {
		if _, ok := rErr.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "read config file: %v\n", rErr)
		}
	}
}
