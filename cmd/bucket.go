package cmd

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qshell"
	"os"
)

var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "Get all buckets of the account",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("%q accepts no args", cmd.CommandPath())
		}
		return nil
	},
	Run: GetBuckets,
}

var domainsCmd = &cobra.Command{
	Use:   "domains <Bucket>",
	Short: "Get all domains of the bucket",
	Args:  cobra.ExactArgs(1),
	Run:   GetDomainsOfBucket,
}

func init() {
	RootCmd.AddCommand(bucketsCmd, domainsCmd)
}

func GetBuckets(cmd *cobra.Command, params []string) {
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
}

func GetDomainsOfBucket(cmd *cobra.Command, params []string) {
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
}
