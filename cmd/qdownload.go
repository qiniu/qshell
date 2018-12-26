package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var qDownloadCmd = &cobra.Command{
	Use:   "qdownload [-c <ThreadCount>] <LocalDownloadConfig>",
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

	var downloadConfig iqshell.DownloadConfig

	configFile := params[0]

	cfh, oErr := os.Open(configFile)
	if oErr != nil {
		fmt.Fprintf(os.Stderr, "open file: %s: %v\n", configFile, oErr)
		os.Exit(1)
	}
	content, rErr := ioutil.ReadAll(cfh)
	if rErr != nil {
		fmt.Fprintf(os.Stderr, "read configFile content: %v\n", rErr)
		os.Exit(1)
	}

	// remove windows utf-8 BOM
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	uErr := json.Unmarshal(content, &downloadConfig)

	if uErr != nil {
		fmt.Fprintf(os.Stderr, "decode configFile content: %v\n", uErr)
		os.Exit(1)
	}

	destFileInfo, err := os.Stat(downloadConfig.DestDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "stat %s: %v\n", downloadConfig.DestDir, err)
		os.Exit(1)
	}

	if !destFileInfo.IsDir() {
		logs.Error("Download dest dir should be a directory")
		os.Exit(iqshell.STATUS_HALT)
	}

	if threadCount < iqshell.MIN_DOWNLOAD_THREAD_COUNT || threadCount > iqshell.MAX_DOWNLOAD_THREAD_COUNT {
		logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			iqshell.MIN_DOWNLOAD_THREAD_COUNT, iqshell.MAX_DOWNLOAD_THREAD_COUNT)

		if threadCount < iqshell.MIN_DOWNLOAD_THREAD_COUNT {
			threadCount = iqshell.MIN_DOWNLOAD_THREAD_COUNT
		} else if threadCount > iqshell.MAX_DOWNLOAD_THREAD_COUNT {
			threadCount = iqshell.MAX_DOWNLOAD_THREAD_COUNT
		}
	}
	iqshell.QiniuDownload(int(threadCount), &downloadConfig)
}
