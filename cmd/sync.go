package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qshell"
	"os"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync <SrcResUrl> <Bucket> <Key>",
	Short: "Sync big file to qiniu bucket",
	Args:  cobra.RangeArgs(3, 4),
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
	key := params[2]

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		logs.Error(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}
	qshell.SetUpHost(&mac, bucket, upHostIp)

	//sync
	tStart := time.Now()
	syncRet, sErr := qshell.Sync(&mac, srcResUrl, bucket, key, qshell.UpHost())
	if sErr != nil {
		logs.Error(sErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", srcResUrl, bucket, key, time.Since(tStart))
	fmt.Println("Hash:", syncRet.Hash)
	fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, FormatFsize(syncRet.Fsize))
	fmt.Println("Mime:", syncRet.MimeType)
}
