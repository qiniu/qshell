package cli

import (
	"github.com/qiniu/log"
	"qshell"
)

func AliListBucket(cmd string, params ...string) {
	if len(params) == 5 || len(params) == 6 {
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
		aliListBucket := qshell.AliListBucket{
			DataCenter:      dataCenter,
			Bucket:          bucket,
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
			Prefix:          prefix,
		}
		err := aliListBucket.ListBucket(listBucketResultFile)
		if err != nil {
			log.Error("List bucket error,", err)
		}
	} else {
		Help(cmd)
	}
}
