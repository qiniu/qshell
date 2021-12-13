package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/operations"
	"github.com/spf13/cobra"
)

var bucketsCmdBuilder = func() *cobra.Command {
	var info = operations.ListInfo{}
	var cmd = &cobra.Command{
		Use:   "buckets",
		Short: "Get all buckets of the account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			operations.List(info)
		},
	}
	return cmd
}

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
			operations.ListDomains(info)
		},
	}
	return cmd
}

func init() {
	RootCmd.AddCommand(
		bucketsCmdBuilder(), // 列举所有 bucket
		domainsCmdBuilder(), // 列举某个 bucket 的 domain
	)
}
