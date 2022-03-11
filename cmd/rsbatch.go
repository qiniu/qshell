package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var batchStatCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "batchstat <Bucket> [-i <KeyListFile>]",
		Short: "Batch stat files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchStatType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchStatus(cfg, info)
		},
	}
	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	setBatchCmdWorkCountFlags(cmd, &info.BatchInfo)
	setBatchCmdSuccessExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdFailExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdForceFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchDeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [-i <KeyListFile>]",
		Short: "Batch delete files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchDeleteType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchDelete(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchChangeMimeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchChangeMimeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchgm <Bucket> [-i <KeyMimeMapFile>]",
		Short: "Batch change the mime type of files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchChangeMimeType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeMime(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchChangeTypeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchtype <Bucket> [-i <KeyFileTypeMapFile>]",
		Short: "Batch change the file type of files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchChangeType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeType(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchDeleteAfterCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [-i <KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchExpireType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchDeleteAfter(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchMoveCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchMoveInfo{}
	var cmd = &cobra.Command{
		Use:   "batchmove <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch move files from bucket to bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchMoveType
			if len(args) > 0 {
				info.SourceBucket = args[0]
			}
			if len(args) > 1 {
				info.DestBucket = args[1]
			}
			operations.BatchMove(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	setBatchCmdOverwriteFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchRenameCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchRenameInfo{}
	var cmd = &cobra.Command{
		Use:   "batchrename <Bucket> [-i <OldNewKeyMapFile>]",
		Short: "Batch rename files in the bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchRenameType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchRename(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	setBatchCmdOverwriteFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchCopyCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchCopyInfo{}
	var cmd = &cobra.Command{
		Use:   "batchcopy <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch copy files from bucket to bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchCopyType
			if len(args) > 0 {
				info.SourceBucket = args[0]
			}
			if len(args) > 1 {
				info.DestBucket = args[1]
			}
			operations.BatchCopy(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	setBatchCmdOverwriteFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchSignCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchPrivateUrlInfo{}
	var cmd = &cobra.Command{
		Use:   "batchsign [-i <ItemListFile>] [-e <Deadline>]",
		Short: "Batch create the private url from the public url list file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchSignType
			info.BatchInfo.Force = true
			operations.BatchPrivateUrl(cfg, info)
		},
	}
	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	cmd.Flags().StringVarP(&info.Deadline, "deadline", "e", "3600", "deadline in seconds, default 3600")
	return cmd
}

var batchFetchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var upHost = ""
	var info = operations.BatchFetchInfo{}
	var cmd = &cobra.Command{
		Use:   "batchfetch <Bucket> [-i <FetchUrlsFile>] [-c <WorkerCount>]",
		Short: "Batch fetch remoteUrls and save them in qiniu Bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchFetchType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(upHost) > 0 {
				cfg.CmdCfg.Hosts.Up = []string{upHost}
			}
			operations.BatchFetch(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	cmd.Flags().StringVarP(&upHost, "up-host", "u", "", "fetch uphost")
	return cmd
}

func setBatchCmdDefaultFlags(cmd *cobra.Command, info *batch.Info) {
	setBatchCmdInputFileFlags(cmd, info)
	setBatchCmdWorkCountFlags(cmd, info)
	setBatchCmdSuccessExportFileFlags(cmd, info)
	setBatchCmdFailExportFileFlags(cmd, info)
	setBatchCmdItemSeparateFlags(cmd, info)
	setBatchCmdForceFlags(cmd, info)
}
func setBatchCmdInputFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.InputFile, "input-file", "i", "", "input file, read from stdin if not set")
}
func setBatchCmdWorkCountFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().IntVarP(&info.WorkCount, "worker", "c", 1, "worker count")
}
func setBatchCmdSuccessExportFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.SuccessExportFilePath, "success-list", "s", "", "rename success list")
}
func setBatchCmdFailExportFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.FailExportFilePath, "failure-list", "e", "", "rename failure list")
}
func setBatchCmdItemSeparateFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields, default is \\t (tab)")
}
func setBatchCmdForceFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.Force, "force", "y", false, "force mode, default false")
}
func setBatchCmdOverwriteFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "w", false, "overwrite mode")
}

func init() {
	registerLoader(rsBatchCmdLoader)
}

func rsBatchCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		batchStatCmdBuilder(cfg),
		batchCopyCmdBuilder(cfg),
		batchMoveCmdBuilder(cfg),
		batchRenameCmdBuilder(cfg),
		batchDeleteCmdBuilder(cfg),
		batchDeleteAfterCmdBuilder(cfg),
		batchChangeMimeCmdBuilder(cfg),
		batchChangeTypeCmdBuilder(cfg),
		batchSignCmdBuilder(cfg),
		batchFetchCmdBuilder(cfg),
	)
}
