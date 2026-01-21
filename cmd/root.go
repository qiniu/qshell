package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/version"
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

var rootCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "qshell",
		Short:                  "Qiniu commandline tool for managing your bucket and CDN",
		Version:                version.Version(),
		BashCompletionFunction: bash_completion_func,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QShellType
			docs.ShowCmdDocument(docs.QShellType)
		},
	}
	cmd.PersistentFlags().BoolVarP(&cfg.StdoutColorful, "colorful", "", false, "console colorful mode")
	cmd.PersistentFlags().BoolVarP(&cfg.Silence, "silence", "", false, "silence mode, The console only outputs warnings、errors and some important information")
	cmd.PersistentFlags().BoolVarP(&cfg.DebugEnable, "debug", "d", false, "debug mode")
	// ddebug 开启 client debug
	cmd.PersistentFlags().BoolVarP(&cfg.DDebugEnable, "ddebug", "D", false, "deep debug mode")
	cmd.PersistentFlags().StringVarP(&cfg.ConfigFilePath, "config", "C", "", "set config file (default is $HOME/.qshell.json)")
	cmd.PersistentFlags().BoolVarP(&cfg.Local, "local", "L", false, "use current directory qshell workspace (default is $HOME/.qshell)")
	cmd.PersistentFlags().BoolVarP(&cfg.Document, "doc", "", false, "document of command")
	return cmd
}

func Execute() {
	cfg := &iqshell.Config{
		Document:       false,
		DebugEnable:    false,
		DDebugEnable:   false,
		ConfigFilePath: "",
		Local:          false,
		CmdCfg: config.Config{
			Log: &config.LogSetting{},
		},
	}

	rootCmd := rootCmdBuilder(cfg)
	load(rootCmd, cfg)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		data.SetCmdStatusError()
	}

	if !data.IsTestMode() && data.GetCmdStatus() != data.StatusOK {
		os.Exit(data.GetCmdStatus())
	}
}

type Loader func(superCmd *cobra.Command, cfg *iqshell.Config)

var loaders = make([]Loader, 0, 20)

func registerLoader(l Loader) {
	if l != nil {
		loaders = append(loaders, l)
	}
}

func load(superCmd *cobra.Command, cfg *iqshell.Config) {
	for _, l := range loaders {
		if l != nil {
			l(superCmd, cfg)
		}
	}
}
