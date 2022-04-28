package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download/operations"
	"github.com/spf13/cobra"
)

var downloadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchDownloadWithConfigInfo{}
	cmd := &cobra.Command{
		Use:   "qdownload [-c <ThreadCount>] <LocalDownloadConfig>",
		Short: "Batch download files from the qiniu bucket",
		Long: `By default qdownload use 5 goroutines to download, it can be customized use -c <count> flag.
And qdownload will use batch stat api or list api to get files info so that it have knowledge to tell whether files
have already in local disk and need to skip download or not.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QDownloadType
			info.BatchInfo.Force = true
			if len(args) > 0 {
				info.LocalDownloadConfig = args[0]
			}
			operations.BatchDownloadWithConfig(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "thread", "c", 5, "num of threads to download files")
	return cmd
}

var getCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DownloadInfo{}
	var cmd = &cobra.Command{
		Use:   "get <Bucket> <Key>",
		Short: "Download a single file from bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.GetType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.DownloadFile(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.ToFile, "outfile", "o", "", "save file as specified by this option")
	cmd.Flags().StringVarP(&info.Domain, "domain", "", "", "domain of request")
	cmd.Flags().BoolVarP(&info.UseGetFileApi, "get-file-api", "", false, "public storage cloud not support, private storage cloud support when has getfile api.")
	cmd.Flags().BoolVarP(&info.RemoveTempWhileError, "remove-temp-while-error", "", false, "remove download temp file while error happened, default is false")
	return cmd
}

func init() {
	registerLoader(downloadCmdLoader)
}

func downloadCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		getCmdBuilder(cfg),
		downloadCmdBuilder(cfg),
	)
}
