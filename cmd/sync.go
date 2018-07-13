package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"os"
	"qiniu/api.v6/auth/digest"
	"qshell"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync <SrcResUrl> <Bucket> <Key> [<UpHostIp>]",
	Short: "Sync big file to qiniu bucket",
	Args:  cobra.RangeArgs(3, 4),
	Run:   Sync,
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func Sync(cmd *cobra.Command, params []string) {
	srcResUrl := params[0]
	bucket := params[1]
	key := params[2]
	upHostIp := ""
	if len(params) == 4 {
		upHostIp = params[3]
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		logs.Error(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}

	if HostFile == "" {
		//get bucket zone info
		bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
		if gErr != nil {
			fmt.Println("Get bucket region info error,", gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		//set up host
		qshell.SetZone(bucketInfo.Region)
	}

	//sync
	tStart := time.Now()
	syncRet, sErr := qshell.Sync(&mac, srcResUrl, bucket, key, upHostIp)
	if sErr != nil {
		logs.Error(sErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", srcResUrl, bucket, key, time.Since(tStart))
	fmt.Println("Hash:", syncRet.Hash)
	fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, FormatFsize(syncRet.Fsize))
	fmt.Println("Mime:", syncRet.MimeType)
}
