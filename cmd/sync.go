package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/config"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync <SrcResUrl> <Buckets> [-k <Key>]",
	Short: "Sync big file to qiniu bucket",
	Args:  cobra.ExactArgs(2),
	Run:   Sync,
}

var (
	upHostIp string
	saveKey  string
)

func init() {
	syncCmd.Flags().BoolVarP(&isResumeV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	syncCmd.Flags().StringVarP(&upHostIp, "uphost", "u", "", "upload host")
	syncCmd.Flags().StringVarP(&saveKey, "key", "k", "", "save as <key> in bucket")
	RootCmd.AddCommand(syncCmd)
}

// 抓取网络中的资源到七牛存储空间，使用与比较大的文件
// 小文件抓取可以使用fetch接口, qshell fetch命令
func Sync(cmd *cobra.Command, params []string) {
	srcResUrl := params[0]
	bucket := params[1]
	var key string
	var kErr error

	if saveKey != "" {
		key = saveKey
	} else {

		key, kErr = iqshell.KeyFromUrl(srcResUrl)
		if kErr != nil {
			fmt.Fprintf(os.Stderr, "get path as key: %v\n", kErr)
			os.Exit(config.STATUS_ERROR)
		}
	}

	bm := iqshell.GetBucketManager()
	//sync
	tStart := time.Now()
	syncRet, sErr := bm.Sync(srcResUrl, bucket, key, upHostIp, isResumeV2)
	if sErr != nil {
		logs.Error(sErr)
		os.Exit(config.STATUS_ERROR)
	}

	fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", srcResUrl, bucket, key, time.Since(tStart))
	fmt.Println("Hash:", syncRet.Hash)
	fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, FormatFsize(syncRet.Fsize))
	fmt.Println("Mime:", syncRet.MimeType)
}
