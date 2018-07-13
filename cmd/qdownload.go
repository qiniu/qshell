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
	qDownloadCmd.Flags().StringVar(&downloadConfig.DestDir, "dest", ".", "dest directory")
	qDownloadCmd.Flags().StringVar(&downloadConfig.Prefix, "prefix", "", "file prefix")
	qDownloadCmd.Flags().StringVar(&downloadConfig.Suffixes, "suffixes", "", "file suffixes, separated by comma")
	qDownloadCmd.Flags().StringVar(&downloadConfig.CdnDomain, "domain", "", "download through domain")
	qDownloadCmd.Flags().StringVar(&downloadConfig.Referer, "referer", "", "referer needed when domain's hotlink protection is turned on")
	qDownloadCmd.Flags().StringVar(&downloadConfig.LogLevel, "ll", "info", "log level")
	qDownloadCmd.Flags().StringVar(&downloadConfig.LogFile, "lf", "", "log file")
	qDownloadCmd.Flags().IntVar(&downloadConfig.LogRotate, "lr", 1, "log rotate days")
	qDownloadCmd.Flags().BoolVar(&downloadConfig.LogStdout, "ls", false, "log to file and stdout")
	qDownloadCmd.Flags().IntVar(&threadCount, "thread", 5, "num of threads to download files")

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

	downloadConfig.IsHostFileSpecified = (HostFile != "")
	qshell.QiniuDownload(int(threadCount), &downloadConfig)
}
