package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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
	cmd.Flags().BoolVarP(&info.Detail, "detail", "", false, "print detail info for domain")
	return cmd
}

var bucketCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.GetBucketInfo{}
	var cmd = &cobra.Command{
		Use:   "bucket <Bucket>",
		Short: "Get bucket info",
		Long:  `Get bucket info`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BucketType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.GetBucket(cfg, info)
		},
	}
	return cmd
}

var mkBucketCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.CreateInfo{}
	var cmd = &cobra.Command{
		Use:   "mkbucket <Bucket>",
		Short: "Create a bucket in region",
		Long: `Create a bucket in region; 

The Bucket name is required to be unique within the scope of the object storage system, consists of 3 to 63 characters, supports lowercase letters, dashes(-) and numbers, and must start and end with lowercase letters or numbers.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.MkBucketType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.Create(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.RegionId, "region", "", "z0", "z0:huadong; z1:huabei; z2:huanan; na0:North America; as0:Southeast Asia;for more:https://developer.qiniu.com/kodo/1671/region-endpoint-fq")
	cmd.Flags().BoolVarP(&info.Private, "private", "", false, "create private bucket")
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
		Long:  "List all the files in the bucket to stdout if output file not specified. Each row of data information is displayed in the following order:\n Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser",
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
		Long:  "List all the files in the bucket using v2/list interface to stdout if output file not specified. Each row of data information is displayed in the following order by default:\n Key\tFileSize\tHash\tPutTime\tMimeType\tStorageType\tEndUser",
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
	cmd.Flags().IntVarP(&info.OutputLimit, "limit", "", -1, "max count of items to output")
	cmd.Flags().Int64VarP(&info.OutputFileMaxLines, "output-file-max-lines", "", 0, "maximum number of lines per output file.")
	cmd.Flags().Int64VarP(&info.OutputFileMaxSize, "output-file-max-size", "", 0, "maximum size of each output file")
	cmd.Flags().StringVarP(&info.StartDate, "start", "s", "", "start date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().StringVarP(&info.EndDate, "end", "e", "", "end date with format yyyy-mm-dd-hh-MM-ss")

	cmd.Flags().StringVarP(&info.StorageTypes, "storages", "", "", "Specify storage type, separated by comma. 0:STANDARD storage, 1:IA storage, 2 means ARCHIVE storage. 3:DEEP_ARCHIVE storage")
	cmd.Flags().StringVarP(&info.MimeTypes, "mimetypes", "", "", "Specify mimetype, separated by comma.")
	cmd.Flags().StringVarP(&info.MinFileSize, "min-file-size", "", "", "Specify min file size")
	cmd.Flags().StringVarP(&info.MaxFileSize, "max-file-size", "", "", "Specify max file size")

	cmd.Flags().BoolVarP(&info.AppendMode, "append", "a", false, "result append to file instead of overwriting")
	cmd.Flags().BoolVarP(&info.Readable, "readable", "r", false, "present file size with human readable format")
	cmd.Flags().StringVarP(&info.ApiVersion, "api-version", "", "v2", "list api version, one of v1 and v2.")
	cmd.Flags().IntVarP(&info.V1Limit, "api-v1-limit", "", 800, "when using the v1 api, the number of entries per enumeration, in the range 1-1000.")
	cmd.Flags().BoolVarP(&info.EnableRecord, "enable-record", "", false, "enable record, ")

	cmd.Flags().StringVarP(&info.OutputFieldsSep, "output-fields-sep", "", data.DefaultLineSeparate, "Each line needs to display the delimiter of the file information.")
	cmd.Flags().StringVarP(&info.ShowFields, "show-fields", "", "", "The file attributes to be displayed on each line, separated by commas. Optional range: Key, Hash, FileSize, PutTime, MimeType, StorageType, EndUser.")

	return cmd
}

func init() {
	registerLoader(bucketCmdLoader)
}

func bucketCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		bucketCmdBuilder(cfg),
		mkBucketCmdBuilder(cfg),
		listBucketCmdBuilder(cfg),
		listBucketCmd2Builder(cfg),
		domainsCmdBuilder(cfg),
	)
}
