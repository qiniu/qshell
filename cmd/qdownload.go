package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"os"
)

var qDownloadCmd = &cobra.Command{
	Use:   "qdownload [<ThreadCount>] <LocalDownloadConfig>",
	Short: "Batch download files from the qiniu bucket",
	Long:  "By default qdownload use 5 goroutines to download, it can be customized use -c <count> flag",
	Args:  cobra.ExactArgs(1),
	Run:   QiniuDownload,
}

var downloadConfig qshell.DownloadConfig
var threadCount int

func init() {
	qDownloadCmd.Flags().StringVarP(&downloadConfig.DestDir, "dest", "t", ".", "dest directory")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.Prefix, "prefix", "p", "", "file prefix")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.Suffixes, "suffixes", "e", "", "file suffixes, separated by comma")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.CdnDomain, "domain", "n", "", "download through domain")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.Referer, "referer", "r", "", "referer needed when domain's hotlink protection is turned on")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.LogLevel, "log-level", "l", "info", "log level")
	qDownloadCmd.Flags().StringVarP(&downloadConfig.LogFile, "log-file", "f", "", "log file")
	qDownloadCmd.Flags().IntVarP(&downloadConfig.LogRotate, "log-rotate", "o", 1, "log rotate days")
	qDownloadCmd.Flags().BoolVarP(&downloadConfig.LogStdout, "log-stdout", "s", false, "log to file and stdout")
	qDownloadCmd.Flags().IntVarP(&threadCount, "thread", "c", 5, "num of threads to download files")

	RootCmd.AddCommand(qDownloadCmd)
}

func QiniuDownload(cmd *cobra.Command, params []string) {

	downloadConfig.Bucket = params[0]

	destFileInfo, err := os.Stat(downloadConfig.DestDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "stat %s: %v\n", downloadConfig.DestDir, err)
		return
	}

	if !destFileInfo.IsDir() {
		logs.Error("Download dest dir should be a directory")
		os.Exit(qshell.STATUS_HALT)
	}

	if threadCount < qshell.MIN_DOWNLOAD_THREAD_COUNT || threadCount > qshell.MAX_DOWNLOAD_THREAD_COUNT {
		logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			qshell.MIN_DOWNLOAD_THREAD_COUNT, qshell.MAX_DOWNLOAD_THREAD_COUNT)

		if threadCount < qshell.MIN_DOWNLOAD_THREAD_COUNT {
			threadCount = qshell.MIN_DOWNLOAD_THREAD_COUNT
		} else if threadCount > qshell.MAX_DOWNLOAD_THREAD_COUNT {
			threadCount = qshell.MAX_DOWNLOAD_THREAD_COUNT
		}
	}
	qshell.QiniuDownload(int(threadCount), &downloadConfig)
}
