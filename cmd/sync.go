package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/qshell/iqshell"
	"os"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync <SrcResUrl> <Buckets[<Key>]",
	Short: "Sync big file to qiniu bucket",
	Args:  cobra.RangeArgs(2, 3),
	Run:   Sync,
}

var upHostIp string

func init() {
	syncCmd.Flags().StringVarP(&upHostIp, "uphost", "u", "", "upload host")
	RootCmd.AddCommand(syncCmd)
}

func Sync(cmd *cobra.Command, params []string) {
	srcResUrl := params[0]
	bucket := params[1]
	var key string
	var kErr error

	if len(params) == 3 {
		key = params[2]
	} else {
		key, kErr = iqshell.KeyFromUrl(srcResUrl)
		if kErr != nil {
			fmt.Fprintf(os.Stderr, "get path as key: %v\n", kErr)
			os.Exit(iqshell.STATUS_ERROR)
		}
	}

	bm := iqshell.GetBucketManager()
	//sync
	tStart := time.Now()
	syncRet, sErr := bm.Sync(srcResUrl, bucket, key)
	if sErr != nil {
		logs.Error(sErr)
		os.Exit(iqshell.STATUS_ERROR)
	}

	fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", srcResUrl, bucket, key, time.Since(tStart))
	fmt.Println("Hash:", syncRet.Hash)
	fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, FormatFsize(syncRet.Fsize))
	fmt.Println("Mime:", syncRet.MimeType)
}
