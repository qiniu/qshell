package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/ali"
	"github.com/spf13/cobra"
)

var aliCmdBuilder = func() *cobra.Command {
	var info = ali.ListBucketInfo{}
	var cmd = &cobra.Command{
		Use:   "alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccesskeySecret> [Prefix] <ListBucketResultFile>",
		Short: "List all the file in the bucket of aliyun oss by prefix",
		Args:  cobra.RangeArgs(5, 6),
		Run: func(cmd *cobra.Command, args []string) {
			info.DataCenter = args[0]
			info.Bucket = args[1]
			info.AccessKey = args[2]
			info.SecretKey = args[3]
			if len(args) == 6 {
				info.Prefix = args[4]
				info.SaveToFile = args[5]
			} else {
				info.SaveToFile = args[4]
			}
			ali.ListBucket(info)
		},
	}
	return cmd
}


func init() {
	RootCmd.AddCommand(aliCmdBuilder())
}
