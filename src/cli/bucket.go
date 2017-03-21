package cli

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"qiniu/api.v6/auth/digest"
	"qshell"
)

func GetBuckets(cmd string, params ...string) {
	if len(params) == 0 {
		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}
		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}
		buckets, err := qshell.GetBuckets(&mac)
		if err != nil {
			logs.Error("Get buckets error,", err)
			os.Exit(qshell.STATUS_ERROR)
		} else {
			if len(buckets) == 0 {
				fmt.Println("No buckets found")
			} else {
				for _, bucket := range buckets {
					fmt.Println(bucket)
				}
			}
		}
	} else {
		CmdHelp(cmd)
	}
}

func GetDomainsOfBucket(cmd string, params ...string) {
	if len(params) == 1 {
		bucket := params[0]
		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}
		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}
		domains, err := qshell.GetDomainsOfBucket(&mac, bucket)
		if err != nil {
			logs.Error("Get domains error,", err)
			os.Exit(qshell.STATUS_ERROR)
		} else {
			if len(domains) == 0 {
				fmt.Printf("No domains found for bucket `%s`\n", bucket)
			} else {
				for _, domain := range domains {
					fmt.Println(domain)
				}
			}
		}
	} else {
		CmdHelp(cmd)
	}
}

func GetFileFromBucket(cmd string, params ...string) {

}
