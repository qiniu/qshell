package cmd

import (
	operations2 "github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var batchStatCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "batchstat <Bucket> [-i <KeyListFile>]",
		Short: "Batch stat files in bucket",
		Long:  "Batch stat files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchStatus(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchDeleteCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [-i <KeyListFile>]",
		Short: "Batch delete files in bucket",
		Long:  "Batch delete files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchDelete(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchChangeMimeCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchChangeMimeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchgm <Bucket> [-i <KeyMimeMapFile>]",
		Short: "Batch change the mime type of files in bucket",
		Long:  "Batch change the mime type of files in bucket, read from stdin if KeyMimeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchChangeMime(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchChangeTypeCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchtype <Bucket> [-i <KeyFileTypeMapFile>]",
		Short: "Batch change the file type of files in bucket",
		Long:  "Batch change the file (storage) type of files in bucket, read from stdin if KeyFileTypeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchChangeType(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchDeleteAfterCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [-i <KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket",
		Long:  "Batch set the deleteAfterDays of the files in bucket, read from stdin if KeyDeleteAfterDaysMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchDeleteAfter(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchMoveCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchMoveInfo{}
	var cmd = &cobra.Command{
		Use:   "batchmove <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch move files from bucket to bucket",
		Long:  "Batch move files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.SourceBucket = args[0]
				info.DestBucket = args[1]
			}
			loadConfig()
			operations2.BatchMove(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchRenameCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchRenameInfo{}
	var cmd = &cobra.Command{
		Use:   "batchrename <Bucket> [-i <OldNewKeyMapFile>]",
		Short: "Batch rename files in the bucket",
		Long:  "Batch rename files in the bucket, read from stdin if OldNewKeyMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			loadConfig()
			operations2.BatchRename(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchCopyCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchCopyInfo{}
	var cmd = &cobra.Command{
		Use:   "batchcopy <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch copy files from bucket to bucket",
		Long:  `Batch copy files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified.
SrcDestKeyMapFile content: line was an copy item
line style:[fromBucketKey][Separator][ToBucketKey]
[ToBucketKey] while use [fromBucketKey] while omitted
`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.SourceBucket = args[0]
				info.DestBucket = args[1]
			}
			loadConfig()
			operations2.BatchCopy(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchSignCmdBuilder = func() *cobra.Command {
	var info = operations2.BatchPrivateUrlInfo{}
	var cmd = &cobra.Command{
		Use:   "batchsign [-i <ItemListFile>] [-e <Deadline>]",
		Short: "Batch create the private url from the public url list file",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			loadConfig()
			operations2.BatchPrivateUrl(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().StringVarP(&info.Deadline, "deadline", "e", "3600", "deadline in seconds")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchFetchCmdBuilder = func() *cobra.Command {
	var upHost = ""
	var info = operations2.BatchFetchInfo{}
	var cmd = &cobra.Command{
		Use:   "batchfetch <Bucket> [-i <FetchUrlsFile>] [-c <WorkerCount>]",
		Short: "Batch fetch remoteUrls and save them in qiniu Bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(upHost) > 0 {
				cfg.CmdCfg.Hosts.Up = []string{upHost}
			}
			loadConfig()
			operations2.BatchFetch(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&upHost, "up-host", "u", "", "fetch uphost")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

func init() {
	rootCmd.AddCommand(
		batchStatCmdBuilder(),
		batchCopyCmdBuilder(),
		batchMoveCmdBuilder(),
		batchRenameCmdBuilder(),
		batchDeleteCmdBuilder(),
		batchDeleteAfterCmdBuilder(),
		batchChangeMimeCmdBuilder(),
		batchChangeTypeCmdBuilder(),
		batchSignCmdBuilder(),
		batchFetchCmdBuilder(),
	)
}
