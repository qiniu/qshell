package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/astaxie/beego/logs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/qiniu/api.v7/v7/client"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 开启命令行的调试模式
	DebugFlag     bool
	DeepDebugInfo bool

	// qshell 版本信息， qshell -v
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

// cobra root cmd, all other commands is children or subchildren of this root cmd
var RootCmd = &cobra.Command{
	Use:                    "qshell",
	Short:                  "Qiniu commandline tool for managing your bucket and CDN",
	Version:                version,
	BashCompletionFunction: bash_completion_func,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&DeepDebugInfo, "ddebug", "D", false, "deep debug mode")
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

	if DeepDebugInfo {
		DebugFlag = true
	}
	//parse command
	if DebugFlag {
		logs.SetLevel(logs.LevelDebug)
		client.TurnOnDebug()
		// master 已合并, v7.5.0 分支没包含次参数 //
		// client.DeepDebugInfo = DeepDebugInfo
	} else {
		logs.SetLevel(logs.LevelInformational)
	}
	logs.SetLogger(logs.AdapterConsole)

	var jsonConfigFile string

	if cfgFile != "" {
		if !strings.HasSuffix(cfgFile, ".json") {
			jsonConfigFile = cfgFile + ".json"
			os.Rename(cfgFile, jsonConfigFile)
		}
		viper.SetConfigFile(jsonConfigFile)
	} else {
		homeDir, hErr := homedir.Dir()
		if hErr != nil {
			fmt.Fprintf(os.Stderr, "get current home directory: %v\n", hErr)
			os.Exit(1)
		}
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".qshell")
	}

	if local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		iqshell.SetRootPath(dir + "/.qshell")
	} else {
		homeDir, hErr := homedir.Dir()
		if hErr != nil {
			fmt.Fprintf(os.Stderr, "get current home directory: %v\n", hErr)
			os.Exit(1)
		}
		iqshell.SetRootPath(homeDir + "/.qshell")
	}
	rootPath := iqshell.RootPath()

	iqshell.SetDefaultAccDBPath(filepath.Join(rootPath, "account.db"))
	iqshell.SetDefaultAccPath(filepath.Join(rootPath, "account.json"))
	iqshell.SetDefaultRsHost(storage.DefaultRsHost)
	iqshell.SetDefaultRsfHost(storage.DefaultRsfHost)
	iqshell.SetDefaultIoHost("iovip.qbox.me")
	iqshell.SetDefaultApiHost(storage.DefaultAPIHost)

	if rErr := viper.ReadInConfig(); rErr != nil {
		if _, ok := rErr.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "read config file: %v\n", rErr)
		}
	}
	os.Rename(jsonConfigFile, cfgFile)
}
