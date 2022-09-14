package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
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
			cfg.CmdCfg.Log.LogLevel = data.NewString(config.DebugKey)
			info.Force = true
			if len(args) > 0 {
				info.LocalDownloadConfig = args[0]
			}
			operations.BatchDownloadWithConfig(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.WorkerCount, "thread", "c", 5, "num of threads to download files")
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
	cmd.Flags().BoolVarP(&info.CheckHash, "check-hash", "", false, "check the consistency of the hash of the local file and the server file after downloading.")
	cmd.Flags().BoolVarP(&info.UseGetFileApi, "get-file-api", "", false, "public storage cloud not support, private storage cloud support when has getfile api.")
	cmd.Flags().BoolVarP(&info.EnableSlice, "enable-slice", "", false, "file download using slices, you need to pay attention to the setting of --slice-file-size-threshold. default is false")
	cmd.Flags().Int64VarP(&info.SliceFileSizeThreshold, "slice-file-size-threshold", "", 40*utils.MB, "the file size threshold that download using slices. when you use --enable-slice option, files larger than this size will be downloaded using slices. Unit: B")
	cmd.Flags().Int64VarP(&info.SliceSize, "slice-size", "", 4*utils.MB, "slice size that download using slices. when you use --enable-slice option, the file will be cut into data blocks according to the slice size, then the data blocks will be downloaded concurrently, and finally these data blocks will be spliced into a file. Unit: B")
	cmd.Flags().IntVarP(&info.SliceConcurrentCount, "slice-concurrent-count", "", 10, "the count of concurrently downloaded slices.")
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
