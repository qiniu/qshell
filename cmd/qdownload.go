package cmd

import (
	"encoding/json"
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

var (
	threadCount int
)

func init() {
	qDownloadCmd.Flags().IntVarP(&threadCount, "thread", "c", 5, "num of threads to download files")

	RootCmd.AddCommand(qDownloadCmd)
}

func QiniuDownload(cmd *cobra.Command, params []string) {

	var downloadConfig qshell.DownloadConfig

	configFile := params[0]

	cfh, oErr := os.Open(configFile)
	if oErr != nil {
		fmt.Fprintf(os.Stderr, "open file: %s: %v\n", configFile, oErr)
		os.Exit(1)
	}
	dErr := json.NewDecoder(cfh).Decode(&downloadConfig)
	if dErr != nil {
		fmt.Fprintf(os.Stderr, "decode configFile content: %v\n", dErr)
		os.Exit(1)
	}

	destFileInfo, err := os.Stat(downloadConfig.DestDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "stat %s: %v\n", downloadConfig.DestDir, err)
		os.Exit(1)
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
