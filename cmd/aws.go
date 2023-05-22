package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/aws"
	"github.com/spf13/cobra"
)

// NewCmdAwsFetch 返回一个cobra.Command指针
func awsFetchCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	info := aws.FetchInfo{}
	cmd := &cobra.Command{
		Use:   "awsfetch [-p <Prefix>] [-n <maxKeys>] [-m <ContinuationToken>] [-c <threadCount>][-u <Qiniu UpHost>] -S <AwsSecretKey> -A <AwsID> <awsBucket> <awsRegion> <qiniuBucket>",
		Short: "Copy data from AWS bucket to qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.AwsFetch
			info.BatchInfo.Force = true
			if len(args) > 0 {
				info.AwsBucketInfo.Bucket = args[0]
			}
			if len(args) > 1 {
				info.AwsBucketInfo.Region = args[1]
			}
			if len(args) > 2 {
				info.QiniuBucket = args[2]
			}
			aws.Fetch(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.AwsBucketInfo.Prefix, "prefix", "p", "", "list AWS bucket with this prefix if set")
	cmd.Flags().Int64VarP(&info.AwsBucketInfo.MaxKeys, "max-keys", "n", 1000, "list AWS bucket with numbers of keys returned each time limited by this number if set")
	cmd.Flags().StringVarP(&info.AwsBucketInfo.CToken, "continuation-token", "m", "", "AWS list continuation token")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "thread-count", "c", 20, "maximum of fetch thread")

	// 没有实际作用
	cmd.Flags().StringVarP(&info.Host, "up-host", "u", "", "Qiniu fetch up host, deprecated")
	_ = cmd.Flags().MarkDeprecated("up-host", "deprecated")

	cmd.Flags().StringVarP(&info.AwsBucketInfo.SecretKey, "aws-secret-key", "S", "", "AWS secret key")
	cmd.Flags().StringVarP(&info.AwsBucketInfo.Id, "aws-id", "A", "", "AWS ID")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "success fetch key list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "error fetch key list")
	cmd.Flags().BoolVarP(&info.BatchInfo.EnableRecord, "enable-record", "", false, "record work progress, and do from last progress while retry")
	cmd.Flags().BoolVarP(&info.BatchInfo.RecordRedoWhileError, "record-redo-while-error", "", false, "when re-executing the command and checking the command task progress record, if a task has already been done and failed, the task will be re-executed. The default is false, and the task will not be re-executed when it detects that the task fails")

	return cmd
}

// NewCmdAwsList 返回一个cobra.Command指针
// 该命令列举亚马逊存储空间中的文件, 会忽略目录
func awsListCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	info := aws.ListBucketInfo{}
	cmd := &cobra.Command{
		Use:   "awslist [-p <Prefix>] [-n <maxKeys>] [-m <ContinuationToken>] -S <AwsSecretKey> -A <AwsID> <awsBucket> <awsRegion>",
		Short: "List Objects in AWS bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.AwsList
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Region = args[1]
			}
			aws.ListBucket(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list AWS bucket with this prefix if set")
	cmd.Flags().Int64VarP(&info.MaxKeys, "max-keys", "n", 1000, "list AWS bucket with numbers of keys returned each time limited by this number if set")
	cmd.Flags().StringVarP(&info.CToken, "continuation-token", "m", "", "AWS list continuation token")
	cmd.Flags().StringVarP(&info.SecretKey, "aws-secret-key", "S", "", "AWS secret key")
	cmd.Flags().StringVarP(&info.Id, "aws-id", "A", "", "AWS ID")
	return cmd
}

func init() {
	registerLoader(awsCmdLoader)
}

func awsCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		awsFetchCmdBuilder(cfg),
		awsListCmdBuilder(cfg),
	)
}
