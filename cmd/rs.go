package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket/operations"
	operations2 "github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

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
			loadConfig()
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
			loadConfig()
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

var statCmdBuilder = func() *cobra.Command {
	var info = operations2.StatusInfo{}
	var cmd = &cobra.Command{
		Use:   "stat <Bucket> <Key>",
		Short: "Get the basic info of a remote file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			loadConfig()
			operations2.Status(info)
		},
	}
	return cmd
}

var forbiddenCmdBuilder = func() *cobra.Command {
	var info = operations2.ForbiddenInfo{}
	var cmd = &cobra.Command{
		Use:   "forbidden <Bucket> <Key>",
		Short: "forbidden file in qiniu bucket",
		Long:  "forbidden object in qiniu bucket, when used with -r option, unforbidden the object",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			loadConfig()
			operations2.ForbiddenObject(info)
		},
	}
	cmd.Flags().BoolVarP(&info.UnForbidden, "reverse", "r", false, "unforbidden object in qiniu bucket")
	return cmd
}

var deleteCmdBuilder = func() *cobra.Command {
	var info = operations2.DeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "delete <Bucket> <Key>",
		Short: "Delete a remote file in the bucket",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			loadConfig()
			operations2.Delete(info)
		},
	}
	return cmd
}

var deleteAfterCmdBuilder = func() *cobra.Command {
	var info = operations2.DeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "expire <Bucket> <Key> <DeleteAfterDays>",
		Short: "Set the deleteAfterDays of a file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.AfterDays = args[2]
			}
			loadConfig()
			operations2.Delete(info)
		},
	}
	return cmd
}

var moveCmdBuilder = func() *cobra.Command {
	var info = operations2.MoveInfo{}
	var cmd = &cobra.Command{
		Use:   "move <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Move/Rename a file and save in bucket",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.SourceBucket = args[0]
				info.SourceKey = args[1]
				info.DestBucket = args[2]
			}
			if len(info.DestKey) == 0 {
				info.DestKey = info.SourceKey
			}
			loadConfig()
			operations2.Move(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket")
	return cmd
}

var copyCmdBuilder = func() *cobra.Command {
	var info = operations2.CopyInfo{}
	var cmd = &cobra.Command{
		Use:   "copy <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Make a copy of a file and save in bucket",
		Example: `copy A.png(bucket:bucketA key:A.png) to B.png(bucket:bucketB key:B.png):
	qshell copy bucketA A.png bucketB -k B.png
you can check if B.png has exists by:
	qshell stat bucketB B.png
`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.SourceBucket = args[0]
				info.SourceKey = args[1]
				info.DestBucket = args[2]
			}
			if len(info.DestKey) == 0 {
				info.DestKey = info.SourceKey
			}
			loadConfig()
			operations2.Copy(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket, use <SrcKey> while omitted")
	return cmd
}

var changeMimeCmdBuilder = func() *cobra.Command {
	var info = operations2.ChangeMimeInfo{}
	var cmd = &cobra.Command{
		Use:   "chgm <Bucket> <Key> <NewMimeType>",
		Short: "Change the mime type of a file",
		Example: `change mimetype of A.png(bucket:bucketA key:A.png) to image/jpeg
	qshell chgm bucketA A.png image/jpeg
and you can check result by command:
	qshell stat bucketA A.png`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.Mime = args[2]
			}
			loadConfig()
			operations2.ChangeMime(info)
		},
	}
	return cmd
}

var changeTypeCmdBuilder = func() *cobra.Command {
	var info = operations2.ChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Long:  `Change the file type of a file, file type must be in 0/1 or 2. 
And 0 means standard storage, 
while 1 means low frequency visit storage,
while 2 means low archive storage.`,
		Example: `change storage type of A.png(bucket:bucketA key:A.png) to archive storage
	qshell chtype bucketA A.png 2
and you can check result by command:
	qshell stat bucketA A.png`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.Type = args[2]
			}
			loadConfig()
			operations2.ChangeType(info)
		},
	}
	return cmd
}

var privateUrlCmdBuilder = func() *cobra.Command {
	var info = operations2.PrivateUrlInfo{}
	var cmd = &cobra.Command{
		Use:   "privateurl <PublicUrl> [<Deadline>]",
		Short: "Create private resource access url",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.PublicUrl = args[0]
			}
			if len(args) > 1 {
				info.Deadline = args[1]
			}
			loadConfig()
			operations2.PrivateUrl(info)
		},
	}
	return cmd
}

var saveAsCmdBuilder = func() *cobra.Command {
	var info = operations2.SaveAsInfo{}
	var cmd = &cobra.Command{
		Use:   "saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>",
		Short: "Create a resource access url with fop and saveas",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.PublicUrl = args[0]
				info.SaveBucket = args[1]
				info.SaveKey = args[2]
			}
			loadConfig()
			operations2.SaveAs(info)
		},
	}
	return cmd
}

var mirrorUpdateCmdBuilder = func() *cobra.Command {
	var info = operations2.MirrorUpdateInfo{}
	var cmd = &cobra.Command{
		Use:   "mirrorupdate <Bucket> <Key>",
		Short: "Fetch and update the file in bucket using mirror storage",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			loadConfig()
			operations2.MirrorUpdate(info)
		},
	}
	return cmd
}

var getCmdBuilder = func() *cobra.Command {
	var info operations2.GetInfo
	var cmd = &cobra.Command{
		Use:   "get <Bucket> <Key>",
		Short: "Download a single file from bucket",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			loadConfig()
			operations2.GetObject(info)
		},
	}

	cmd.Flags().StringVarP(&info.SaveToFile, "outfile", "o", "", "save file as specified by this option")

	return cmd
}

var fetchCmdBuilder = func() *cobra.Command {
	var info operations2.FetchInfo
	var cmd = &cobra.Command{
		Use:   "fetch <RemoteResourceUrl> <Bucket> [-k <Key>]",
		Short: "Fetch a remote resource by url and save in bucket",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.FromUrl = args[0]
				info.Bucket = args[1]
			}
			loadConfig()
			operations2.Fetch(info)
		},
	}

	cmd.Flags().StringVarP(&info.Key, "key", "k", "", "filename saved in bucket")

	return cmd
}

func init() {

	rootCmd.AddCommand(
		listBucketCmdBuilder(),
		listBucketCmd2Builder(),
		statCmdBuilder(),
		forbiddenCmdBuilder(),
		deleteCmdBuilder(),
		deleteAfterCmdBuilder(),
		moveCmdBuilder(),
		copyCmdBuilder(),
		changeMimeCmdBuilder(),
		changeTypeCmdBuilder(),
		privateUrlCmdBuilder(),
		saveAsCmdBuilder(),
		mirrorUpdateCmdBuilder(),
		getCmdBuilder(),
		fetchCmdBuilder(),
	)
}
