package cmd

import (
	operations2 "github.com/qiniu/qshell/v2/iqshell/storage/servers/operations"
	"github.com/spf13/cobra"
)

var bucketsCmdBuilder = func() *cobra.Command {
	var info = operations2.ListInfo{}
	var cmd = &cobra.Command{
		Use:   "buckets",
		Short: "Get all buckets of the account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			loadConfig()
			operations2.List(info)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(
		bucketsCmdBuilder(), // 列举所有 bucket
	)
}
