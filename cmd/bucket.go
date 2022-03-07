package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/operations"
	"github.com/spf13/cobra"
)

var domainsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ListDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "domains <Bucket>",
		Short: "Get all domains of the bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.DomainsType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.ListDomains(cfg, info)
		},
	}
	return cmd
}

var listBucketCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ListInfo{
		Marker:     "",
		StartDate:  "",
		EndDate:    "",
		AppendMode: false,
		Readable:   false,
		Delimiter:  "",
		MaxRetry:   20,
	}
	var cmd = &cobra.Command{
		Use:   "listbucket <Bucket>",
		Short: "List all the files in the bucket",
		Long:  "List all the files in the bucket to stdout if output file not specified",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ListBucketType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.List(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list by prefix")
	cmd.Flags().StringVarP(&info.SaveToFile, "out", "o", "", "output file")
	return cmd
}

var listBucketCmd2Builder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ListInfo{}
	var cmd = &cobra.Command{
		Use:   "listbucket2 <Bucket>",
		Short: "List all the files in the bucket using v2/list interface",
		Long:  "List all the files in the bucket using v2/list interface to stdout if output file not specified",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ListBucket2Type
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.List(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.Marker, "marker", "m", "", "list marker")
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list by prefix")
	cmd.Flags().StringVarP(&info.Suffixes, "suffixes", "q", "", "list by key suffixes, separated by comma")
	cmd.Flags().IntVarP(&info.MaxRetry, "max-retry", "x", -1, "max retries when error occurred")
	cmd.Flags().StringVarP(&info.SaveToFile, "out", "o", "", "output file")
	cmd.Flags().StringVarP(&info.StartDate, "start", "s", "", "start date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().StringVarP(&info.EndDate, "end", "e", "", "end date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().BoolVarP(&info.AppendMode, "append", "a", false, "result append to file instead of overwriting")
	cmd.Flags().BoolVarP(&info.Readable, "readable", "r", false, "present file size with human readable format")
	cmd.Flags().IntVarP(&info.Limit, "limit", "", -1, "max count of items to output")

	return cmd
}

func init() {
	registerLoader(bucketCmdLoader)
}

func bucketCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config)  {
	superCmd.AddCommand(
		listBucketCmdBuilder(cfg),
		listBucketCmd2Builder(cfg),
		domainsCmdBuilder(cfg),
	)
}