package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/cobra"
)

func versionCmdBuilder() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, params []string) {
			loadConfig()
			log.Alert(data.Version)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(versionCmdBuilder())
}

type Status struct {
	isCancel bool
}
