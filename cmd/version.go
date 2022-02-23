package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/cobra"
)

func versionCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, params []string) {
			log.Alert(data.Version)
		},
	}
	return cmd
}

func init() {
	registerLoader(versionCmdLoader)
}

func versionCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		versionCmdBuilder(cfg),
	)
}
