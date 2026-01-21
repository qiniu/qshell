package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var matchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.MatchInfo{}
	cmd := &cobra.Command{
		Use:   "match <Bucket> <Key> <LocalFile>",
		Short: "Verify that the local file matches the Qiniu cloud storage file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.MatchType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.LocalFile = args[2]
			}
			operations.Match(cfg, info)
		},
	}
	return cmd
}

var batchMatchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchMatchInfo{}
	cmd := &cobra.Command{
		Use:   "batchmatch <Bucket> <LocalFileDir>",
		Short: "Batch Verify that the local file matches the Qiniu cloud storage file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.BatchMatchType
			info.BatchInfo.Force = true
			info.BatchInfo.EnableStdin = true
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.LocalFileDir = args[1]
			}
			operations.BatchMatch(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file, read from stdin if not set")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields, default is \\t (tab)")
	cmd.Flags().BoolVarP(&info.BatchInfo.EnableRecord, "enable-record", "", false, "record work progress, and do from last progress while retry")
	cmd.Flags().BoolVarP(&info.BatchInfo.RecordRedoWhileError, "record-redo-while-error", "", false, "when re-executing the command and checking the command task progress record, if a task has already been done and failed, the task will be re-executed. The default is false, and the task will not be re-executed when it detects that the task fails")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	return cmd
}

func init() {
	registerLoader(matchCmdLoader)
}

func matchCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		matchCmdBuilder(cfg),
		batchMatchCmdBuilder(cfg),
	)
}
