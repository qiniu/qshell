package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/operations"
	"github.com/spf13/cobra"
)

var syncCmdBuilder = func() *cobra.Command {
	info := operations.SyncUploadInfo{}
	cmd := &cobra.Command{
		Use:   "sync <SrcResUrl> <Buckets> [-k <Key>]",
		Short: "Sync big file to qiniu bucket",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.ResourceUrl = args[0]
				info.Bucket = args[1]
			}
		},
	}
	cmd.Flags().BoolVarP(&info.IsResumeV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().StringVarP(&info.UpHostIp, "uphost", "u", "", "upload host")
	cmd.Flags().StringVarP(&info.Key, "key", "k", "", "save as <key> in bucket")
	return cmd
}

func init() {
	RootCmd.AddCommand(syncCmdBuilder())
}
