package cli

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"qiniu/api.v6/auth/digest"
	"qshell"
	"time"
)

func Sync(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
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
			return
		}

		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}
		//get bucket zone info
		bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
		if gErr != nil {
			fmt.Println("Get bucket region info error,", gErr)
			return
		}

		//set up host
		qshell.SetZone(bucketInfo.Region)

		//sync
		tStart := time.Now()
		syncRet, sErr := qshell.Sync(&mac, srcResUrl, bucket, key, upHostIp)
		if sErr != nil {
			logs.Error(sErr)
			return
		}

		fmt.Printf("Sync %s => %s:%s Success, Duration: %s!\n", srcResUrl, bucket, key, time.Since(tStart))
		fmt.Println("Hash:", syncRet.Hash)
		fmt.Printf("Fsize: %d (%s)\n", syncRet.Fsize, FormatFsize(syncRet.Fsize))
		fmt.Println("Mime:", syncRet.MimeType)
	} else {
		CmdHelp(cmd)
	}
}
