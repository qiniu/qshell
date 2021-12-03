package cmd

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage"

	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
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

// 【buckets】获取一个用户的所有的存储空间
func GetBuckets(cmd *cobra.Command, params []string) {

	bm := storage.GetBucketManager()
	buckets, err := bm.Buckets(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get buckets error: %v\n", err)
		os.Exit(1)
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

// 【domains】获取一个空间绑定的CDN域名
func GetDomainsOfBucket(cmd *cobra.Command, params []string) {
	bucket := params[0]
	bm := storage.GetBucketManager()
	domains, err := bm.DomainsOfBucket(bucket)

	if err != nil {
		logs.Error("Get domains error: ", err)
		os.Exit(data.STATUS_ERROR)
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
