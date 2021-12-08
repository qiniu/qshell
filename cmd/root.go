package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/spf13/cobra"
	"os"
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
	Version:                data.Version,
	BashCompletionFunction: bash_completion_func,
}

var initFuncs []func()

func OnInitialize(f ...func()) {
	initFuncs = append(initFuncs, f...)
}

func init() {
	cobra.OnInitialize(func() {
		initConfig()
		for _, f := range initFuncs {
			f()
		}
	})

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	// ddebug 开启 client debug
	RootCmd.PersistentFlags().BoolVarP(&DeepDebugInfo, "ddebug", "D", false, "deep debug mode")
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "config file (default is $HOME/.qshell.json)")
	RootCmd.PersistentFlags().BoolVarP(&local, "local", "L", false, "use current directory as config file path")
}

func initConfig() {
	workspacePath := ""
	if local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		workspacePath = dir
	}

	err := iqshell.Load(iqshell.Config{
		DebugEnable:    DebugFlag,
		DDebugEnable:   DeepDebugInfo,
		ConfigFilePath: cfgFile,
		WorkspacePath:  workspacePath,
	})
	_, _ = fmt.Fprintf(os.Stderr, "load error: %v\n", err)
}
