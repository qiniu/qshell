package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/operations"
	"github.com/spf13/cobra"
)

var domainsCmdBuilder = func() *cobra.Command {
	var info = operations.ListDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "domains <Bucket>",
		Short: "Get all domains of the bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			prepare(cmd, nil)
			operations.ListDomains(info)
		},
	}
	return cmd
}

var listBucketCmdBuilder = func() *cobra.Command {
	var info = operations.ListInfo{
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
		Long:  "List all the files in the bucket to stdout if ListBucketResultFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			prepare(cmd, nil)
			operations.List(info)
		},
	}
	cmd.Flags().StringVarP(&info.Marker, "marker", "m", "", "list marker")
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list by prefix")
	cmd.Flags().StringVarP(&info.SaveToFile, "out", "o", "", "output file")
	return cmd
}

var listBucketCmd2Builder = func() *cobra.Command {
	var info = operations.ListInfo{}
	var cmd = &cobra.Command{
		Use:   "listbucket2 <Bucket>",
		Short: "List all the files in the bucket using v2/list interface",
		Long:  "List all the files in the bucket to stdout if ListBucketResultFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			prepare(cmd, nil)
			operations.List(info)
		},
	}

	cmd.Flags().StringVarP(&info.Marker, "marker", "m", "", "list marker")
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list by prefix")
	cmd.Flags().StringVarP(&info.Suffixes, "suffixes", "q", "", "list by key suffixes, separated by comma")
	cmd.Flags().IntVarP(&info.MaxRetry, "max-retry", "x", -1, "max retries when error occurred")
	cmd.Flags().StringVarP(&info.SaveToFile, "out", "o", "", "output file")
	cmd.Flags().StringVarP(&info.StartDate, "start", "s", "", "start date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().StringVarP(&info.EndDate, "end", "e", "", "end date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().BoolVarP(&info.AppendMode, "append", "a", false, "append to file")
	cmd.Flags().BoolVarP(&info.Readable, "readable", "r", false, "present file size with human readable format")

	return cmd
}

func init() {
	rootCmd.AddCommand(
		domainsCmdBuilder(), // 列举某个 bucket 的 domain
	)
}
