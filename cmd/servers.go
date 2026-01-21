package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/servers/operations"
	"github.com/spf13/cobra"
)

var bucketsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ListInfo{}
	cmd := &cobra.Command{
		Use:   "buckets",
		Short: "Get all buckets of the account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BucketsType
			operations.List(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.Region, "region", "", "", "region of bucket; z0, z1, z2, as0, na0 etc")
	cmd.Flags().BoolVarP(&info.Detail, "detail", "", false, "print detail info for bucket")
	return cmd
}

func init() {
	registerLoader(serversCmdLoader)
}

func serversCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		bucketsCmdBuilder(cfg),
	)
}
