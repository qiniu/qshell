package cli

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/log"
	"qshell"
)

func GetBuckets(cmd string, params ...string) {
	if len(params) == 0 {
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		buckets, err := qshell.GetBuckets(&mac)
		if err != nil {
			log.Error("Get buckets error,", err)
		} else {
			for _, bucket := range buckets {
				fmt.Println(bucket)
			}
		}
	} else {
		CmdHelp(cmd)
	}
}
