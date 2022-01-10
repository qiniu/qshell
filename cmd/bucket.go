package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/operations"
	"github.com/spf13/cobra"
)

var domainsCmdBuilder = func() *cobra.Command {
	var info = operations.ListDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "domains <Bucket>",
		Short: "Get all domains of the bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations.ListDomains(info)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(
		domainsCmdBuilder(), // 列举某个 bucket 的 domain
	)
}
