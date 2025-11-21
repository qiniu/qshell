package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
)

var batchStatCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "batchstat <Bucket> [-i <KeyListFile>]",
		Short: "Batch stat files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchStatType
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchStatus(cfg, info)
		},
	}
	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	setBatchCmdWorkerCountFlags(cmd, &info.BatchInfo)
	setBatchCmdMinWorkerCountFlags(cmd, &info.BatchInfo)
	setBatchCmdWorkerCountIncreasePeriodFlags(cmd, &info.BatchInfo)
	setBatchCmdSuccessExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdFailExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdResultExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdForceFlags(cmd, &info.BatchInfo)
	setBatchCmdEnableRecordFlags(cmd, &info.BatchInfo)
	setBatchCmdRecordRedoWhileErrorFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchForbiddenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchChangeStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "batchforbidden <Bucket> [-i <KeyListFile>] [-r]",
		Short: "Batch forbidden files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchForbiddenType
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeStatus(cfg, info)
		},
	}
	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	setBatchCmdWorkerCountFlags(cmd, &info.BatchInfo)
	setBatchCmdMinWorkerCountFlags(cmd, &info.BatchInfo)
	setBatchCmdWorkerCountIncreasePeriodFlags(cmd, &info.BatchInfo)
	setBatchCmdSuccessExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdFailExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdForceFlags(cmd, &info.BatchInfo)
	setBatchCmdEnableRecordFlags(cmd, &info.BatchInfo)
	setBatchCmdRecordRedoWhileErrorFlags(cmd, &info.BatchInfo)
	cmd.Flags().BoolVarP(&info.UnForbidden, "reverse", "r", false, "unforbidden object in qiniu bucket")
	return cmd
}

var batchDeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [-i <KeyListFile>]",
		Short: "Batch delete files in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchDeleteType
			info.BatchInfo.EnableStdin = true
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
			info.BatchInfo.EnableStdin = true
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
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeType(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchRestoreArCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchRestoreArchiveInfo{}
	var cmd = &cobra.Command{
		Use:   "batchrestorear <Bucket> <FreezeAfterDays>",
		Short: `Batch unfreeze archive file and file freeze after <FreezeAfterDays> days, <FreezeAfterDays> value should be between 1 and 7, include 1 and 7`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchRestoreArchiveType
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.FreezeAfterDays = args[1]
			}
			operations.BatchRestoreArchive(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchDeleteAfterCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [-i <KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket. DeleteAfterDays:great than or equal to 0, 0: cancel expiration time, unit: day",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchExpireType
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchDeleteAfter(cfg, info)
		},
	}
	setBatchCmdDefaultFlags(cmd, &info.BatchInfo)
	return cmd
}

var batchChangeLifecycleCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.BatchChangeLifecycleInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchlifecycle <Bucket> [-i <KeyFile>] [--to-ia-after-days <ToIAAfterDays>] [--to-archive-after-days <ToArchiveAfterDays>] [--to-deep-archive-after-days <ToDeepArchiveAfterDays>] [--delete-after-days <DeleteAfterDays>]",
		Short: "Set the lifecycle of some file.",
		Long: `Set the lifecycle of some file. <KeyFile> contain all file keys that need to set. one key per line.
Lifecycle value must great than or equal to -1, unit: day.
* less than  -1: there's no point and it won't trigger any effect
* equal to   -1: cancel lifecycle
* equal to    0: there's no point and it won't trigger any effect
* bigger than 0: set lifecycle`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchChangeLifecycle
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeLifecycle(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.ToIAAfterDays, "to-ia-after-days", "", 0, "to IA storage after some days. the range is -1 or bigger than 0. -1 means cancel to IA storage")
	cmd.Flags().IntVarP(&info.ToArchiveIRAfterDays, "to-archive-ir-after-days", "", 0, "to ARCHIVE_IR storage after some days. the range is -1 or bigger than 0. -1 means cancel to ARCHIVE_IR storage")
	cmd.Flags().IntVarP(&info.ToArchiveAfterDays, "to-archive-after-days", "", 0, "to ARCHIVE storage after some days. the range is -1 or bigger than 0. -1 means cancel to ARCHIVE storage")
	cmd.Flags().IntVarP(&info.ToDeepArchiveAfterDays, "to-deep-archive-after-days", "", 0, "to DEEP_ARCHIVE storage after some days. the range is -1 or bigger than 0. -1 means cancel to DEEP_ARCHIVE storage")
	cmd.Flags().IntVarP(&info.ToIntelligentTieringAfterDays, "to-intelligent-tiering-after-days", "", 0, "to INTELLIGENT_TIERING storage after some days. the range is -1 or bigger than 0. -1 means cancel to INTELLIGENT_TIERING storage")
	cmd.Flags().IntVarP(&info.DeleteAfterDays, "delete-after-days", "", 0, "delete after some days. the range is -1 or bigger than 0. -1 means cancel to delete")
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
			info.BatchInfo.EnableStdin = true
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
			info.BatchInfo.EnableStdin = true
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
			info.BatchInfo.EnableStdin = true
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
			info.BatchInfo.EnableStdin = true
			info.BatchInfo.Force = true
			operations.BatchPrivateUrl(cfg, info)
		},
	}
	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	setBatchCmdResultExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdEnableRecordFlags(cmd, &info.BatchInfo)
	setBatchCmdRecordRedoWhileErrorFlags(cmd, &info.BatchInfo)
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
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(upHost) > 0 {
				cfg.CmdCfg.Hosts.Up = []string{upHost}
			}
			operations.BatchFetch(cfg, info)
		},
	}

	setBatchCmdInputFileFlags(cmd, &info.BatchInfo)
	setBatchCmdEnableRecordFlags(cmd, &info.BatchInfo)
	setBatchCmdRecordRedoWhileErrorFlags(cmd, &info.BatchInfo)
	setBatchCmdSuccessExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdFailExportFileFlags(cmd, &info.BatchInfo)
	setBatchCmdItemSeparateFlags(cmd, &info.BatchInfo)
	setBatchCmdForceFlags(cmd, &info.BatchInfo)
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&upHost, "up-host", "u", "", "fetch uphost")
	return cmd
}

func setBatchCmdDefaultFlags(cmd *cobra.Command, info *batch.Info) {
	setBatchCmdInputFileFlags(cmd, info)
	setBatchCmdWorkerCountFlags(cmd, info)
	setBatchCmdMinWorkerCountFlags(cmd, info)
	setBatchCmdWorkerCountIncreasePeriodFlags(cmd, info)
	setBatchCmdEnableRecordFlags(cmd, info)
	setBatchCmdRecordRedoWhileErrorFlags(cmd, info)
	setBatchCmdSuccessExportFileFlags(cmd, info)
	setBatchCmdFailExportFileFlags(cmd, info)
	setBatchCmdItemSeparateFlags(cmd, info)
	setBatchCmdForceFlags(cmd, info)
}
func setBatchCmdInputFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.InputFile, "input-file", "i", "", "input file, read from stdin if not set")
}
func setBatchCmdForceFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.Force, "force", "y", false, "force mode, default false")
}
func setBatchCmdWorkerCountFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().IntVarP(&info.WorkerCount, "worker", "c", 4, "worker count. 1 means the number of objects in one operation is 250 and if configured as 10 , the number of objects in one operation is 2500. This value needs to be consistent with the upper limit of Qiniu’s operation, otherwise unexpected errors will occur. Under normal circumstances you do not need to adjust this value and if you need please carefully.")
}
func setBatchCmdMinWorkerCountFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().IntVarP(&info.MinWorkerCount, "min-worker", "", 1, "min worker count. 1 means the number of objects in one operation is 1000 and if configured as 3 , the number of objects in one operation is 3000. for more, please refer to worker")
}
func setBatchCmdWorkerCountIncreasePeriodFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().IntVarP(&info.WorkerCountIncreasePeriod, "worker-count-increase-period", "", 60, "worker count increase period. when the worker count is too big, an overrun error will be triggered. In order to alleviate this problem, qshell will automatically reduce the worker count. In order to complete the operation as quickly as possible, qshell will periodically increase the worker count. unit: second")
}
func setBatchCmdItemSeparateFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields, default is \\t (tab)")
}
func setBatchCmdEnableRecordFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.EnableRecord, "enable-record", "", false, "record work progress, and do from last progress while retry")
}
func setBatchCmdRecordRedoWhileErrorFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.RecordRedoWhileError, "record-redo-while-error", "", false, "when re-executing the command and checking the command task progress record, if a task has already been done and failed, the task will be re-executed. The default is false, and the task will not be re-executed when it detects that the task fails")
}
func setBatchCmdSuccessExportFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.SuccessExportFilePath, "success-list", "s", "", "specifies the file path where the successful file list is saved")
}
func setBatchCmdFailExportFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.FailExportFilePath, "failure-list", "e", "", "specifies the file path where the failure file list is saved")
}
func setBatchCmdOverwriteFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "w", false, "overwrite mode")
	_ = cmd.Flags().MarkShorthandDeprecated("overwrite", "deprecated and use --overwrite instead")
}
func setBatchCmdResultExportFileFlags(cmd *cobra.Command, info *batch.Info) {
	cmd.Flags().StringVarP(&info.ResultExportFilePath, "outfile", "o", "", "specifies the file path where the results is saved")
}

func init() {
	registerLoader(rsBatchCmdLoader)
}

func rsBatchCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		batchStatCmdBuilder(cfg),
		batchForbiddenCmdBuilder(cfg),
		batchCopyCmdBuilder(cfg),
		batchMoveCmdBuilder(cfg),
		batchRenameCmdBuilder(cfg),
		batchDeleteCmdBuilder(cfg),
		batchChangeLifecycleCmdBuilder(cfg),
		batchDeleteAfterCmdBuilder(cfg),
		batchChangeMimeCmdBuilder(cfg),
		batchChangeTypeCmdBuilder(cfg),
		batchRestoreArCmdBuilder(cfg),
		batchSignCmdBuilder(cfg),
		batchFetchCmdBuilder(cfg),
	)
}
