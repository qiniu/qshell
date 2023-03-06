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
		Short: "Batch asynchronous fetch network resources to qiniu Bucket",
		Long: `Batch asynchronous fetch in two steps:
1. Initiate an asynchronous fetch request. The success of the request does not mean that the fetch is successful. Step 2 is required to detect whether the fetch is really successful.
2. Check if the fetch is successful, you can skip this step with the long option: --disable-check-fetch-result`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ABFetch
			info.BatchInfo.ItemSeparate = "\t" // 此处用户不可定义
			info.BatchInfo.EnableStdin = true
			info.BatchInfo.OperationCountPerRequest = 1
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
	cmd.Flags().IntVarP(&info.FileType, "storage-type", "g", 0, "storage type, same to --file-type")
	cmd.Flags().IntVarP(&info.FileType, "file-type", "", 0, "storage type, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage")
	cmd.Flags().BoolVar(&info.Overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file with urls")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "thread-count", "c", 20, "thread count")
	cmd.Flags().BoolVarP(&info.BatchInfo.EnableRecord, "enable-record", "", false, "record work progress, and do from last progress while retry")
	cmd.Flags().BoolVarP(&info.BatchInfo.RecordRedoWhileError, "record-redo-while-error", "", false, "when re-executing the command and checking the command task progress record, if a task has already been done and failed, the task will be re-executed. The default is false, and the task will not be re-executed when it detects that the task fails")
	cmd.Flags().BoolVarP(&info.DisableCheckFetchResult, "disable-check-fetch-result", "", false, "not check async result after fetch")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "success fetch list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "error fetch list")

	// 废弃
	_ = cmd.Flags().MarkDeprecated("storage-type", "use --file-type instead")
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

func asyncFetchCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		asyncFetchCmdBuilder(cfg),
		asyncCheckCmdBuilder(cfg),
	)
}
