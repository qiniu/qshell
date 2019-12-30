package cmd

import (
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
)

var aliCmd = &cobra.Command{
	Use:   "alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccesskeySecret> [Prefix] <ListBucketResultFile>",
	Short: "List all the file in the bucket of aliyun oss by prefix",
	Args:  cobra.RangeArgs(5, 6),
	Run:   AliListBucket,
}

func init() {
	RootCmd.AddCommand(aliCmd)
}

// 【alilistbucket】列举阿里空间中的文件列表
func AliListBucket(cmd *cobra.Command, params []string) {
	dataCenter := params[0]
	bucket := params[1]
	accessKeyId := params[2]
	accessKeySecret := params[3]
	listBucketResultFile := ""
	prefix := ""
	if len(params) == 6 {
		prefix = params[4]
		listBucketResultFile = params[5]
	} else {
		listBucketResultFile = params[4]
	}
	aliListBucket := iqshell.AliListBucket{
		DataCenter:      dataCenter,
		Bucket:          bucket,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Prefix:          prefix,
	}
	err := aliListBucket.ListBucket(listBucketResultFile)
	if err != nil {
		logs.Error("List bucket error,", err)
		return
	}
	logs.Info("List bucket done!")
}
