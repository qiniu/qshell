package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var statCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.StatusInfo{}
	var cmd = &cobra.Command{
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
	var info = operations.ForbiddenInfo{}
	var cmd = &cobra.Command{
		Use:   "forbidden <Bucket> <Key>",
		Short: "Forbidden file in qiniu bucket",
		Long:  "Forbidden object in qiniu bucket, when used with -r option, unforbidden the object",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			operations.ForbiddenObject(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.UnForbidden, "reverse", "r", false, "unforbidden object in qiniu bucket")
	return cmd
}

var deleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.DeleteInfo{}
	var cmd = &cobra.Command{
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

var deleteAfterCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.DeleteAfterInfo{}
	var cmd = &cobra.Command{
		Use:   "expire <Bucket> <Key> <DeleteAfterDays>",
		Short: "Set the deleteAfterDays of a file",
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
	var info = operations.MoveInfo{}
	var cmd = &cobra.Command{
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
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket")
	return cmd
}

var renameCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.RenameInfo{}
	var cmd = &cobra.Command{
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
	return cmd
}

var copyCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.CopyInfo{}
	var cmd = &cobra.Command{
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
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket, use <SrcKey> while omitted")
	return cmd
}

var changeMimeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ChangeMimeInfo{}
	var cmd = &cobra.Command{
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
	var info = operations.ChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Long: `Change the file type of a file, file type must be one of 0, 1, 2. 
And 0 means STANDARD storage, 
while 1 means IA storage,
while 2 means ARCHIVE storage.
while 3 means DEEP_ARCHIVE storage.`,
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
	var info = operations.RestoreArchiveInfo{}
	var cmd = &cobra.Command{
		Use: "restorear <Bucket> <Key> <FreezeAfterDays>",
		Short: `Unfreeze archive file and file freeze after <FreezeAfterDays> days,
<FreezeAfterDays> value should be between 1 and 7, include 1 and 7`,
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
	var info = operations.PrivateUrlInfo{}
	var cmd = &cobra.Command{
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
	var info = operations.SaveAsInfo{}
	var cmd = &cobra.Command{
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
	var info = operations.MirrorUpdateInfo{}
	var cmd = &cobra.Command{
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
	var info = operations.MirrorUpdateInfo{}
	var cmd = &cobra.Command{
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
	var cmd = &cobra.Command{
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
