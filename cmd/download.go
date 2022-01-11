package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download/operations"
	"github.com/spf13/cobra"
)

var downloadCmdBuilder = func() *cobra.Command {
	info := operations.BatchDownloadInfo{}
	cmd := &cobra.Command{
		Use:   "qdownload [-c <ThreadCount>] <LocalDownloadConfig>",
		Short: "Batch download files from the qiniu bucket",
		Long: `By default qdownload use 5 goroutines to download, it can be customized use -c <count> flag.
And qdownload will use batch stat api or list api to get files info so that it have knowledge to tell whether files
have already in local disk and need to skip download or not.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cfg.DownloadConfigFile = args[0]
			}
			loadConfig()
			operations.BatchDownload(info)
		},
	}
	cmd.Flags().IntVarP(&info.GroupInfo.WorkCount, "thread", "c", 5, "num of threads to download files")
	return cmd
}

func init() {
	rootCmd.AddCommand(downloadCmdBuilder())
}
