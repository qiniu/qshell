package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

func asyncFetchCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchAsyncFetchInfo{}
	cmd := &cobra.Command{
		Use:   "abfetch <Bucket> [-i <urlList>]",
		Short: "Async Batch fetch network resources to qiniu Bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ABFetch
			info.GroupInfo.ItemSeparate = "\t" // 此处用户不可定义
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchAsyncFetch(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.Host, "host", "t", "", "the host when download from fetch url")
	cmd.Flags().StringVarP(&info.CallbackUrl, "callback-url", "a", "", "callback url")
	cmd.Flags().StringVarP(&info.CallbackBody, "callback-body", "b", "", "callback body")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "callback HOST")
	cmd.Flags().IntVarP(&info.FileType, "storage-type", "g", 0, "storage type")
	cmd.Flags().StringVarP(&info.GroupInfo.InputFile, "input-file", "i", "", "input file with urls")
	cmd.Flags().IntVarP(&info.GroupInfo.WorkCount, "thread-count", "c", 20, "thread count")
	cmd.Flags().StringVarP(&info.GroupInfo.SuccessExportFilePath, "success-list", "s", "", "success fetch list")
	cmd.Flags().StringVarP(&info.GroupInfo.FailExportFilePath, "failure-list", "e", "", "error fetch list")

	return cmd
}

// NewCmdAsyncCheck 用来查询异步抓取的结果
func asyncCheckCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	info := operations.CheckAsyncFetchStatusInfo{}
	cmd := &cobra.Command{
		Use:   "acheck <Bucket> <ID>",
		Short: "Check Async fetch status",
		Run: func(cmd *cobra.Command, args []string) {
			 cfg.CmdCfg.CmdId = docs.ACheckType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Id = args[1]
			}
				operations.CheckAsyncFetchStatus(cfg, info)
		},
	}
	return cmd
}

func init() {
	registerLoader(asyncFetchCmdLoader)
}

func asyncFetchCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config)  {
	superCmd.AddCommand(
		asyncFetchCmdBuilder(cfg),
		asyncCheckCmdBuilder(cfg),
		)
}
