package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
)

var statCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.StatusInfo{}
	cmd := &cobra.Command{
		Use:   "stat <Bucket> <Key>",
		Short: "Get the basic info of a remote file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.StatType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.Status(cfg, info)
		},
	}
	return cmd
}

var forbiddenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ForbiddenInfo{}
	cmd := &cobra.Command{
		Use:   "forbidden <Bucket> <Key>",
		Short: "Forbidden file in qiniu bucket",
		Long:  "Forbidden object in qiniu bucket, when used with -r option, unforbidden the object",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ForbiddenType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.ForbiddenObject(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.UnForbidden, "reverse", "r", false, "unforbidden object in qiniu bucket")
	return cmd
}

var deleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DeleteInfo{}
	cmd := &cobra.Command{
		Use:   "delete <Bucket> <Key>",
		Short: "Delete a remote file in the bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.DeleteType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.Delete(cfg, info)
		},
	}
	return cmd
}

var changeLifecycleCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := &operations.ChangeLifecycleInfo{}
	cmd := &cobra.Command{
		Use:   "chlifecycle <Bucket> <Key> [--to-ia-after-days <ToIAAfterDays>] [--to-archive-after-days <ToArchiveAfterDays>] [--to-deep-archive-after-days <ToDeepArchiveAfterDays>] [--to-intelligent-tiering-after-days <ToIntelligentTieringAfterDays>] [--delete-after-days <DeleteAfterDays>]",
		Short: "Set the lifecycle of a file.",
		Long: `Set the lifecycle of a file. Lifecycle value must great than or equal to -1, unit: day.
* less than  -1: there's no point and it won't trigger any effect
* equal to   -1: cancel lifecycle
* equal to    0: there's no point and it won't trigger any effect
* bigger than 0: set lifecycle`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ChangeLifecycle
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.ChangeLifecycle(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.ToIAAfterDays, "to-ia-after-days", "", 0, "to IA storage after some days. the range is -1 or bigger than 0. -1 means cancel to IA storage")
	cmd.Flags().IntVarP(&info.ToArchiveIRAfterDays, "to-archive-ir-after-days", "", 0, "to ARCHIVE_IR storage after some days. the range is -1 or bigger than 0. -1 means cancel to ARCHIVE_IR storage")
	cmd.Flags().IntVarP(&info.ToArchiveAfterDays, "to-archive-after-days", "", 0, "to ARCHIVE storage after some days. the range is -1 or bigger than 0. -1 means cancel to ARCHIVE storage")
	cmd.Flags().IntVarP(&info.ToDeepArchiveAfterDays, "to-deep-archive-after-days", "", 0, "to DEEP_ARCHIVE storage after some days. the range is -1 or bigger than 0. -1 means cancel to DEEP_ARCHIVE storage")
	cmd.Flags().IntVarP(&info.ToIntelligentTieringAfterDays, "to-intelligent-tiering-after-days", "", 0, "to INTELLIGENT_TIERING storage after some days. the range is -1 or bigger than 0. -1 means cancel to INTELLIGENT_TIERING storage")
	cmd.Flags().IntVarP(&info.DeleteAfterDays, "delete-after-days", "", 0, "delete after some days. the range is -1 or bigger than 0. -1 means cancel to delete")
	return cmd
}

var deleteAfterCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DeleteAfterInfo{}
	cmd := &cobra.Command{
		Use:   "expire <Bucket> <Key> <DeleteAfterDays>",
		Short: "Set the deleteAfterDays of a file. DeleteAfterDays:great than or equal to 0, 0: cancel expiration time, unit: day",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ExpireType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.AfterDays = args[2]
			}
			operations.DeleteAfter(cfg, info)
		},
	}
	return cmd
}

var moveCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.MoveInfo{}
	cmd := &cobra.Command{
		Use:   "move <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Move/Rename a file and save in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.MoveType
			if len(args) > 0 {
				info.SourceBucket = args[0]
			}
			if len(args) > 1 {
				info.SourceKey = args[1]
			}
			if len(args) > 2 {
				info.DestBucket = args[2]
			}
			if len(info.DestKey) == 0 {
				info.DestKey = info.SourceKey
			}
			operations.Move(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	_ = cmd.Flags().MarkShorthandDeprecated("overwrite", "deprecated and use --overwrite instead")

	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket")
	return cmd
}

var renameCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.RenameInfo{}
	cmd := &cobra.Command{
		Use:   "rename <SrcBucket> <SrcKey> <DestKey>",
		Short: "Make a rename of a file and save in the bucket",
		Example: `rename A.png(bucket:bucketA key:A.png) to B.png(bucket:bucketA key:B.png):
	qshell rename bucketA A.png B.png
you can check if B.png has exists by:
	qshell stat bucketB B.png
`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RenameType
			if len(args) > 0 {
				info.SourceBucket = args[0]
				info.DestBucket = args[0]
			}
			if len(args) > 1 {
				info.SourceKey = args[1]
			}
			if len(args) > 2 {
				info.DestKey = args[2]
			}
			operations.Rename(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	_ = cmd.Flags().MarkShorthandDeprecated("overwrite", "deprecated and use --overwrite instead")

	return cmd
}

var copyCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.CopyInfo{}
	cmd := &cobra.Command{
		Use:   "copy <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Make a copy of a file and save in bucket",
		Example: `copy A.png(bucket:bucketA key:A.png) to B.png(bucket:bucketB key:B.png):
	qshell copy bucketA A.png bucketB -k B.png
you can check if B.png has exists by:
	qshell stat bucketB B.png
`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.CopyType
			if len(args) > 0 {
				info.SourceBucket = args[0]
			}
			if len(args) > 1 {
				info.SourceKey = args[1]
			}
			if len(args) > 2 {
				info.DestBucket = args[2]
			}
			operations.Copy(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	_ = cmd.Flags().MarkShorthandDeprecated("overwrite", "deprecated and use --overwrite instead")

	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket, use <SrcKey> while omitted")
	return cmd
}

var changeMimeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ChangeMimeInfo{}
	cmd := &cobra.Command{
		Use:   "chgm <Bucket> <Key> <NewMimeType>",
		Short: "Change the mime type of a file",
		Example: `change mimetype of A.png(bucket:bucketA key:A.png) to image/jpeg
	qshell chgm bucketA A.png image/jpeg
and you can check result by command:
	qshell stat bucketA A.png`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ChangeMimeType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.Mime = args[2]
			}
			operations.ChangeMime(cfg, info)
		},
	}
	return cmd
}

var changeTypeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ChangeTypeInfo{}
	cmd := &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Long: `Change the file type of a file, file type must be one of 0, 1, 2, 3. 
And 0 means STANDARD storage, 
while 1 means IA storage,
while 2 means ARCHIVE storage.
while 3 means DEEP_ARCHIVE storage.
while 4 means ARCHIVE_IR storage.
while 5 means INTELLIGENT_TIERING storage`,
		Example: `change storage type of A.png(bucket:bucketA key:A.png) to ARCHIVE storage
	qshell chtype bucketA A.png 2
and you can check result by command:
	qshell stat bucketA A.png`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ChangeType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.Type = args[2]
			}
			operations.ChangeType(cfg, info)
		},
	}
	return cmd
}

var restoreArCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.RestoreArchiveInfo{}
	cmd := &cobra.Command{
		Use:   "restorear <Bucket> <Key> <FreezeAfterDays>",
		Short: `Unfreeze archive file and file freeze after <FreezeAfterDays> days, <FreezeAfterDays> value should be between 1 and 7, include 1 and 7`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RestoreArchiveType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.FreezeAfterDays = args[2]
			}
			operations.RestoreArchive(cfg, info)
		},
	}

	return cmd
}

var privateUrlCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.PrivateUrlInfo{}
	cmd := &cobra.Command{
		Use:   "privateurl <PublicUrl> [<Deadline>]",
		Short: "Create private resource access url",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.PrivateUrlType
			if len(args) > 0 {
				info.PublicUrl = args[0]
			}
			if len(args) > 1 {
				info.Deadline = args[1]
			}
			operations.PrivateUrl(cfg, info)
		},
	}
	return cmd
}

var saveAsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.SaveAsInfo{}
	cmd := &cobra.Command{
		Use:   "saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>",
		Short: "Create a resource access url with fop and saveas",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SaveAsType
			if len(args) > 0 {
				info.PublicUrl = args[0]
			}
			if len(args) > 1 {
				info.SaveBucket = args[1]
			}
			if len(args) > 2 {
				info.SaveKey = args[2]
			}
			operations.SaveAs(cfg, info)
		},
	}
	return cmd
}

var mirrorUpdateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.MirrorUpdateInfo{}
	cmd := &cobra.Command{
		Use:   "mirrorupdate <Bucket> <Key>",
		Short: "Fetch and update the file in bucket using mirror storage",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.MirrorUpdateType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.MirrorUpdate(cfg, info)
		},
	}
	return cmd
}

var prefetchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.MirrorUpdateInfo{}
	cmd := &cobra.Command{
		Use:   "prefetch <Bucket> <Key>",
		Short: "Fetch and update the file in bucket using mirror storage",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.PrefetchType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.MirrorUpdate(cfg, info)
		},
	}
	return cmd
}

var fetchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info operations.FetchInfo
	cmd := &cobra.Command{
		Use:   "fetch <RemoteResourceUrl> <Bucket> [-k <Key>]",
		Short: "Fetch a remote resource by url and save in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.FetchType
			if len(args) > 0 {
				info.FromUrl = args[0]
			}
			if len(args) > 1 {
				info.Bucket = args[1]
			}
			operations.Fetch(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.Key, "key", "k", "", "filename saved in bucket")

	return cmd
}

func init() {
	registerLoader(rsCmdLoader)
}

func rsCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		statCmdBuilder(cfg),
		forbiddenCmdBuilder(cfg),
		changeLifecycleCmdBuilder(cfg),
		deleteCmdBuilder(cfg),
		deleteAfterCmdBuilder(cfg),
		moveCmdBuilder(cfg),
		renameCmdBuilder(cfg),
		copyCmdBuilder(cfg),
		changeMimeCmdBuilder(cfg),
		changeTypeCmdBuilder(cfg),
		restoreArCmdBuilder(cfg),
		privateUrlCmdBuilder(cfg),
		saveAsCmdBuilder(cfg),
		mirrorUpdateCmdBuilder(cfg),
		fetchCmdBuilder(cfg),
		prefetchCmdBuilder(cfg),
	)
}
