package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download/operations"
	"github.com/spf13/cobra"
	"os"
)

var downloadCmdBuilder = func() *cobra.Command {
	info := operations.BatchDownloadInfo{}
	cmd := &cobra.Command{
		Use:   "qdownload [-c <ThreadCount>] <LocalDownloadConfig>",
		Short: "Batch download files from the qiniu bucket",
		Long: `By default qdownload use 5 goroutines to download, it can be customized use -c <count> flag.
And qdownload will use batch stat api or list api to get files info so that it have knowledge to tell whether files
have already in local disk and need to skip download or not.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.QDownloadType
			if len(args) > 0 {
				cfg.DownloadConfigFile = args[0]
			}
			cfg.CmdCfg.Download.LogSetting = &config.LogSetting{
				LogLevel:  config.InfoKey,
				LogFile:   "",
				LogRotate: 0,
				LogStdout: data.TrueString,
			}
			if prepare(cmd, &info) {
				if len(args) == 0 {
					fmt.Fprintln(os.Stdout, "LocalDownloadConfig can't empty")
					return
				}
				operations.BatchDownload(info)
			}
		},
	}
	cmd.Flags().IntVarP(&info.GroupInfo.WorkCount, "thread", "c", 5, "num of threads to download files")
	return cmd
}

var getCmdBuilder = func() *cobra.Command {
	info := operations.DownloadInfo{}
	var cmd = &cobra.Command{
		Use:   "get <Bucket> <Key>",
		Short: "Download a single file from bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.GetType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if prepare(cmd, &info) {
				operations.DownloadFile(info)
			}
		},
	}

	cmd.Flags().StringVarP(&info.ToFile, "outfile", "o", "", "save file as specified by this option")
	cmd.Flags().StringVarP(&info.Domain, "domain", "", "", "domain of server")

	return cmd
}

func init() {
	rootCmd.AddCommand(getCmdBuilder())
	rootCmd.AddCommand(downloadCmdBuilder())
}
