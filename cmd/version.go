package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/spf13/cobra"
)

func versionCmdBuilder() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, params []string) {
			fmt.Println(data.Version)
		},
	}
	return cmd
}

func init() {
	RootCmd.AddCommand(versionCmdBuilder())
}
