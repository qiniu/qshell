package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

func asyncFetchCmdBuilder() *cobra.Command {
	info := operations.BatchAsyncFetchInfo{}
	cmd := &cobra.Command{
		Use:   "abfetch <Bucket> [-i <urlList>]",
		Short: "Async Batch fetch network resources to qiniu Bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchAsyncFetch(info)
		},
	}

	cmd.Flags().StringVarP(&info.Host, "host", "t", "", "download HOST header")
	cmd.Flags().StringVarP(&info.CallbackUrl, "callback-url", "a", "", "callback url")
	cmd.Flags().StringVarP(&info.CallbackBody, "callback-body", "b", "", "callback body")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "callback HOST")
	cmd.Flags().IntVarP(&info.FileType, "storage-type", "g", 0, "storage type")
	cmd.Flags().StringVarP(&info.InputFile, "input-file", "i", "", "input file with urls")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "thread-count", "c", 20, "thread count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "success fetch list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "error fetch list")

	return cmd
}

// NewCmdAsyncCheck 用来查询异步抓取的结果
func asyncCheckCmdBuilder() *cobra.Command {

	info := operations.CheckAsyncFetchStatusInfo{}
	cmd := &cobra.Command{
		Use:   "acheck <Bucket> <ID>",
		Short: "Check Async fetch status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Id = args[1]
			}
			operations.CheckAsyncFetchStatus(info)
		},
	}
	return cmd
}

func init() {
	RootCmd.AddCommand(asyncFetchCmdBuilder())
	RootCmd.AddCommand(asyncCheckCmdBuilder())
}
