package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	storage2 "github.com/qiniu/qshell/v2/iqshell/storage"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"

	"github.com/spf13/cobra"
)

var listBucketCmdBuilder = func() *cobra.Command {
	var info = operations.ListInfo{
		StartDate:  "",
		EndDate:    "",
		AppendMode: false,
		Readable:   false,
		ApiInfo: rs.ListApiInfo{
			Delimiter: "",
			MaxRetry:  20,
		},
	}
	var cmd = &cobra.Command{
		Use:   "listbucket <Bucket>",
		Short: "List all the files in the bucket",
		Long:  "List all the files in the bucket to stdout if ListBucketResultFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.ApiInfo.Bucket = args[0]
			}
			operations.List(info)
		},
	}
	cmd.Flags().StringVarP(&info.ApiInfo.Marker, "marker", "m", "", "list marker")
	cmd.Flags().StringVarP(&info.ApiInfo.Prefix, "prefix", "p", "", "list by prefix")
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
				info.ApiInfo.Bucket = args[0]
			}
			operations.List(info)
		},
	}

	cmd.Flags().StringVarP(&info.ApiInfo.Marker, "marker", "m", "", "list marker")
	cmd.Flags().StringVarP(&info.ApiInfo.Prefix, "prefix", "p", "", "list by prefix")
	cmd.Flags().StringVarP(&info.Suffixes, "suffixes", "q", "", "list by key suffixes, separated by comma")
	cmd.Flags().IntVarP(&info.ApiInfo.MaxRetry, "max-retry", "x", -1, "max retries when error occurred")
	cmd.Flags().StringVarP(&info.SaveToFile, "out", "o", "", "output file")
	cmd.Flags().StringVarP(&info.StartDate, "start", "s", "", "start date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().StringVarP(&info.EndDate, "end", "e", "", "end date with format yyyy-mm-dd-hh-MM-ss")
	cmd.Flags().BoolVarP(&info.AppendMode, "append", "a", false, "append to file")
	cmd.Flags().BoolVarP(&info.Readable, "readable", "r", false, "present file size with human readable format")

	return cmd
}

var statCmdBuilder = func() *cobra.Command {
	var info = rs.StatusApiInfo{}
	var cmd = &cobra.Command{
		Use:   "stat <Bucket> <Key>",
		Short: "Get the basic info of a remote file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			operations.Status(info)
		},
	}
	return cmd
}

var forbiddenCmdBuilder = func() *cobra.Command {
	var info = operations.ForbiddenInfo{}
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
			operations.ForbiddenObject(info)
		},
	}
	cmd.Flags().BoolVarP(&info.UnForbidden, "reverse", "r", false, "unforbidden object in qiniu bucket")
	return cmd
}

var deleteCmdBuilder = func() *cobra.Command {
	var info = operations.DeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "delete <Bucket> <Key>",
		Short: "Delete a remote file in the bucket",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			operations.Delete(info)
		},
	}
	return cmd
}

var deleteAfterCmdBuilder = func() *cobra.Command {
	var info = operations.DeleteInfo{}
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
			operations.Delete(info)
		},
	}
	return cmd
}

var moveCmdBuilder = func() *cobra.Command {
	var info = rs.MoveApiInfo{}
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
			operations.Move(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket")
	return cmd
}

var copyCmdBuilder = func() *cobra.Command {
	var info = rs.CopyApiInfo{}
	var cmd = &cobra.Command{
		Use:   "copy <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Make a copy of a file and save in bucket",
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
			operations.Copy(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Force, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.DestKey, "key", "k", "", "filename saved in bucket")
	return cmd
}

var changeMimeCmdBuilder = func() *cobra.Command {
	var info = rs.ChangeMimeApiInfo{}
	var cmd = &cobra.Command{
		Use:   "chgm <Bucket> <Key> <NewMimeType>",
		Short: "Change the mime type of a file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.Mime = args[2]
			}
			operations.ChangeMime(info)
		},
	}
	return cmd
}

var changeTypeCmdBuilder = func() *cobra.Command {
	var info = operations.ChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Long:  "Change the file type of a file, file type must be in 0 or 1. And 0 means standard storage, while 1 means low frequency visit storage.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.Type = args[2]
			}
			operations.ChangeType(info)
		},
	}
	return cmd
}

var privateUrlCmdBuilder = func() *cobra.Command {
	var info = operations.PrivateUrlInfo{}
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
			operations.PrivateUrl(info)
		},
	}
	return cmd
}

var (
	qGetCmd = &cobra.Command{
		Use:   "get <Bucket> <Key>",
		Short: "Download a single file from bucket",
		Args:  cobra.ExactArgs(2),
		Run:   Get,
	}
	dirCacheCmd = &cobra.Command{
		Use:   "dircache <DirCacheRootPath>",
		Short: "Cache the directory structure of a file path",
		Long:  "Cache the directory structure of a file path to a file, \nif <DirCacheResultFile> not specified, cache to stdout",
		Args:  cobra.ExactArgs(1),
		Run:   DirCache,
	}
	fetchCmd = &cobra.Command{
		Use:   "fetch <RemoteResourceUrl> <Bucket> [-k <Key>]",
		Short: "Fetch a remote resource by url and save in bucket",
		Args:  cobra.ExactArgs(2),
		Run:   Fetch,
	}
	mirrorCmd = &cobra.Command{
		Use:   "mirrorupdate <Bucket> <Key>",
		Short: "Fetch and update the file in bucket using mirror storage",
		Args:  cobra.ExactArgs(2),
		Run:   Prefetch,
	}
	saveAsCmd = &cobra.Command{
		Use:   "saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>",
		Short: "Create a resource access url with fop and saveas",
		Args:  cobra.ExactArgs(3),
		Run:   Saveas,
	}
	m3u8DelCmd = &cobra.Command{
		Use:   "m3u8delete <Bucket> <M3u8Key>",
		Short: "Delete m3u8 playlist and the slices it references",
		Args:  cobra.ExactArgs(2),
		Run:   M3u8Delete,
	}
	m3u8RepCmd = &cobra.Command{
		Use:   "m3u8replace <Bucket> <M3u8Key> [<NewDomain>]",
		Short: "Replace m3u8 domain in the playlist",
		Args:  cobra.RangeArgs(2, 3),
		Run:   M3u8Replace,
	}
)

var (
	outFile                  string
	cOverwrite               bool
	finalKey                 string
	tsUrlRemoveSparePreSlash bool
)

func init() {
	dirCacheCmd.Flags().StringVarP(&outFile, "outfile", "o", "", "output filepath")
	qGetCmd.Flags().StringVarP(&outFile, "outfile", "o", "", "save file as specified by this option")

	fetchCmd.Flags().StringVarP(&finalKey, "key", "k", "", "filename saved in bucket")
	m3u8RepCmd.Flags().BoolVarP(&tsUrlRemoveSparePreSlash, "remove-spare-pre-slash", "r", true, "remove spare prefix slash(/) , only keep one slash if ts path has prefix / ")

	RootCmd.AddCommand(
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
	)

	RootCmd.AddCommand(qGetCmd, dirCacheCmd, fetchCmd, mirrorCmd,
		saveAsCmd, m3u8DelCmd, m3u8RepCmd)
}

// 【dircache】扫描本地文件目录， 形成一个关于文件信息的文本文件
func DirCache(cmd *cobra.Command, params []string) {
	var cacheResultFile string
	cacheRootPath := params[0]

	cacheResultFile = outFile
	if cacheResultFile == "" {
		cacheResultFile = "stdout"
	}
	_, retErr := utils.DirCache(cacheRootPath, cacheResultFile)
	if retErr != nil {
		os.Exit(data.STATUS_ERROR)
	}
}

// 【get】下载七牛存储中的一个文件， 该命令不需要存储空间绑定有可访问的CDN域名
func Get(cmd *cobra.Command, params []string) {

	bucket := params[0]
	key := params[1]

	destFile := key
	if outFile != "" {
		destFile = outFile
	}

	bm := storage2.GetBucketManager()
	err := bm.Get(bucket, key, destFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
}

// 【fetch】通过http链接抓取网上的资源到七牛存储空间
func Fetch(cmd *cobra.Command, params []string) {
	remoteResUrl := params[0]
	bucket := params[1]

	var err error
	if finalKey == "" {
		finalKey, err = utils.KeyFromUrl(remoteResUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get key from url failed: %v\n", err)
			os.Exit(data.STATUS_ERROR)
		}
	}

	bm := storage2.GetBucketManager()
	fetchResult, err := bm.Fetch(remoteResUrl, bucket, finalKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fetch error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	} else {
		fmt.Println("Key:", fetchResult.Key)
		fmt.Println("Hash:", fetchResult.Hash)
		fmt.Printf("Fsize: %d (%s)\n", fetchResult.Fsize, utils.FormatFileSize(fetchResult.Fsize))
		fmt.Println("Mime:", fetchResult.MimeType)
	}
}

// 【cdnprefetch】CDN文件预取, 预取文件到CDN节点和父层节点
func Prefetch(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]

	bm := storage2.GetBucketManager()
	err := bm.Prefetch(bucket, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prefetch error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
}

// 【saveas】打印输出主动saveas链接
func Saveas(cmd *cobra.Command, params []string) {
	publicUrl := params[0]
	saveBucket := params[1]
	saveKey := params[2]

	bm := storage2.GetBucketManager()
	url, err := bm.Saveas(publicUrl, saveBucket, saveKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Saveas error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	} else {
		fmt.Println(url)
	}
}

// 【m3u8delete】删除m3u8文件，包括m3u8文件本身和分片文件
func M3u8Delete(cmd *cobra.Command, params []string) {
	bucket := params[0]
	m3u8Key := params[1]

	bm := storage2.GetBucketManager()
	m3u8FileList, err := bm.M3u8FileList(bucket, m3u8Key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get m3u8 file list error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
	entryCnt := len(m3u8FileList)
	if entryCnt == 0 {
		fmt.Fprintln(os.Stderr, "no m3u8 slices found")
		os.Exit(data.STATUS_ERROR)
	}
	fileExporter, nErr := storage2.NewFileExporter("", "", "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	if entryCnt <= BATCH_ALLOW_MAX {
		batchDelete(m3u8FileList, bm, fileExporter)
	} else {
		batchCnt := entryCnt / BATCH_ALLOW_MAX
		for i := 0; i < batchCnt; i++ {
			end := (i + 1) * BATCH_ALLOW_MAX
			if end > entryCnt {
				end = entryCnt
			}
			entriesToDelete := m3u8FileList[i*BATCH_ALLOW_MAX : end]
			batchDelete(entriesToDelete, bm, fileExporter)
		}
	}
}

// 【m3u8replace】替换m3u8文件中的域名信息
func M3u8Replace(cmd *cobra.Command, params []string) {
	bucket := params[0]
	m3u8Key := params[1]
	var newDomain string
	if len(params) == 3 {
		newDomain = params[2]
	}
	bm := storage2.GetBucketManager()
	err := bm.M3u8ReplaceDomain(bucket, m3u8Key, newDomain, tsUrlRemoveSparePreSlash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "m3u8 replace domain error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
}

