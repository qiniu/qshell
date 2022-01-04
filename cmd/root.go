package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/spf13/cobra"
	"os"
)

var cfg iqshell.Config

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

// rootCmd cobra root cmd, all other commands is children or subchildren of this root cmd
var rootCmd = &cobra.Command{
	Use:                    "qshell",
	Short:                  "Qiniu commandline tool for managing your bucket and CDN",
	Version:                data.Version,
	BashCompletionFunction: bash_completion_func,
}

func init() {
	cfg = iqshell.Config{
		DebugEnable:    false,
		DDebugEnable:   false,
		ConfigFilePath: "",
		Local:          false,
	}
	rootCmd.PersistentFlags().BoolVarP(&cfg.DebugEnable, "debug", "d", false, "debug mode")
	// ddebug 开启 client debug
	rootCmd.PersistentFlags().BoolVarP(&cfg.DDebugEnable, "ddebug", "D", false, "deep debug mode")
	rootCmd.PersistentFlags().StringVarP(&cfg.ConfigFilePath, "config", "C", "", "config file (default is $HOME/.qshell.json)")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Local, "local", "L", false, "use current directory as config file path")
}

func loadConfig() {
	err := iqshell.Load(cfg)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
