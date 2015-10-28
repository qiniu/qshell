package cli

import (
	"fmt"
	"qiniu/api.v6/auth/digest"
	"qiniu/log"
	"qshell"
)

func GetBuckets(cmd string, params ...string) {
	if len(params) == 0 {
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
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

func GetDomainsOfBucket(cmd string, params ...string) {
	if len(params) == 1 {
		bucket := params[0]
		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		domains, err := qshell.GetDomainsOfBucket(&mac, bucket)
		if err != nil {
			log.Error("Get domains error,", err)
		} else {
			for _, domain := range domains {
				fmt.Println(domain)
			}
		}
	} else {
		CmdHelp(cmd)
	}
}
