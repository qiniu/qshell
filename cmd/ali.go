package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/ali"
	"github.com/spf13/cobra"
)

var aliCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = ali.ListBucketInfo{}
	var cmd = &cobra.Command{
		Use:   "alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccessKeySecret> [Prefix] <ListBucketResultFile>",
		Short: "List all the file in the bucket of aliyun oss by prefix",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.AliListBucket
			if len(args) > 0 {
				info.DataCenter = args[0]
			}
			if len(args) > 1 {
				info.Bucket = args[1]
			}
			if len(args) > 2 {
				info.AccessKey = args[2]
			}
			if len(args) > 3 {
				info.SecretKey = args[3]
			}
			if len(args) > 4 {
				if len(args) == 6 {
					info.Prefix = args[4]
					info.SaveToFile = args[5]
				} else {
					info.SaveToFile = args[4]
				}
			}
			ali.ListBucket(cfg, info)
		},
	}
	return cmd
}

func init() {
	registerLoader(aliCmdLoader)
}

func aliCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config)  {
	superCmd.AddCommand(aliCmdBuilder(cfg))
}