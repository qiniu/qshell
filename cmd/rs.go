package cmd

import (
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	lsBucketCmd = &cobra.Command{
		Use:   "listbucket <Bucket>",
		Short: "List all the files in the bucket",
		Long:  "List all the files in the bucket to stdout if ListBucketResultFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   ListBucket,
	}
	lsBucketCmd2 = &cobra.Command{
		Use:   "listbucket2 <Bucket>",
		Short: "List all the files in the bucket using v2/list interface",
		Long:  "List all the files in the bucket to stdout if ListBucketResultFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   ListBucket2,
	}
	statCmd = &cobra.Command{
		Use:   "stat <Bucket> <Key>",
		Short: "Get the basic info of a remote file",
		Args:  cobra.ExactArgs(2),
		Run:   Stat,
	}
	delCmd = &cobra.Command{
		Use:   "delete <Bucket> <Key>",
		Short: "Delete a remote file in the bucket",
		Args:  cobra.ExactArgs(2),
		Run:   Delete,
	}
	moveCmd = &cobra.Command{
		Use:   "move <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Move/Rename a file and save in bucket",
		Args:  cobra.ExactArgs(3),
		Run:   Move,
	}
	copyCmd = &cobra.Command{
		Use:   "copy <SrcBucket> <SrcKey> <DestBucket> [-k <DestKey>]",
		Short: "Make a copy of a file and save in bucket",
		Args:  cobra.ExactArgs(3),
		Run:   Copy,
	}
	chgmCmd = &cobra.Command{
		Use:   "chgm <Bucket> <Key> <NewMimeType>",
		Short: "Change the mime type of a file",
		Args:  cobra.ExactArgs(3),
		Run:   Chgm,
	}
	chtypeCmd = &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Long:  "Change the file type of a file, file type must be in 0 or 1. And 0 means standard storage, while 1 means low frequency visit storage.",
		Args:  cobra.ExactArgs(3),
		Run:   Chtype,
	}
	delafterCmd = &cobra.Command{
		Use:   "expire <Bucket> <Key> <DeleteAfterDays>",
		Short: "Set the deleteAfterDays of a file",
		Args:  cobra.ExactArgs(3),
		Run:   DeleteAfterDays,
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
	privateUrlCmd = &cobra.Command{
		Use:   "privateurl <PublicUrl> [<Deadline>]",
		Short: "Create private resource access url",
		Args:  cobra.RangeArgs(1, 2),
		Run:   PrivateUrl,
	}
)

var (
	outFile    string
	listMarker string
	prefix     string
	suffixes   string
	mOverwrite bool
	cOverwrite bool
	startDate  string
	endDate    string
	maxRetry   int
	finalKey   string
	appendMode bool
)

func init() {
	dirCacheCmd.Flags().StringVarP(&outFile, "outfile", "o", "", "output filepath")
	qGetCmd.Flags().StringVarP(&outFile, "outfile", "o", "", "save file as specified by this option")

	lsBucketCmd.Flags().StringVarP(&listMarker, "marker", "m", "", "list marker")
	lsBucketCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "list by prefix")
	lsBucketCmd.Flags().StringVarP(&outFile, "out", "o", "", "output file")

	lsBucketCmd2.Flags().StringVarP(&listMarker, "marker", "m", "", "list marker")
	lsBucketCmd2.Flags().StringVarP(&prefix, "prefix", "p", "", "list by prefix")
	lsBucketCmd2.Flags().StringVarP(&suffixes, "suffixes", "q", "", "list by key suffixes, separated by comma")
	lsBucketCmd2.Flags().IntVarP(&maxRetry, "max-retry", "x", -1, "max retries when error occurred")
	lsBucketCmd2.Flags().StringVarP(&outFile, "out", "o", "", "output file")
	lsBucketCmd2.Flags().StringVarP(&startDate, "start", "s", "", "start date with format yyyy-mm-dd-hh-MM-ss")
	lsBucketCmd2.Flags().StringVarP(&endDate, "end", "e", "", "end date with format yyyy-mm-dd-hh-MM-ss")
	lsBucketCmd2.Flags().BoolVarP(&appendMode, "append", "a", false, "append to file")

	moveCmd.Flags().BoolVarP(&mOverwrite, "overwrite", "w", false, "overwrite mode")
	moveCmd.Flags().StringVarP(&finalKey, "key", "k", "", "filename saved in bucket")
	copyCmd.Flags().BoolVarP(&cOverwrite, "overwrite", "w", false, "overwrite mode")
	copyCmd.Flags().StringVarP(&finalKey, "key", "k", "", "filename saved in bucket")
	fetchCmd.Flags().StringVarP(&finalKey, "key", "k", "", "filename saved in bucket")

	RootCmd.AddCommand(qGetCmd, dirCacheCmd, lsBucketCmd, statCmd, delCmd, moveCmd,
		copyCmd, chgmCmd, chtypeCmd, delafterCmd, fetchCmd, mirrorCmd,
		saveAsCmd, m3u8DelCmd, m3u8RepCmd, privateUrlCmd, lsBucketCmd2)
}

func DirCache(cmd *cobra.Command, params []string) {
	var cacheResultFile string
	cacheRootPath := params[0]

	cacheResultFile = outFile
	if cacheResultFile == "" {
		cacheResultFile = "stdout"
	}
	_, retErr := iqshell.DirCache(cacheRootPath, cacheResultFile)
	if retErr != nil {
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func ListBucket2(cmd *cobra.Command, params []string) {
	bucket := params[0]

	var dateParser = func(datestr string) (time.Time, error) {
		var dttm [6]int

		if datestr == "" {
			return time.Time{}, nil
		}
		fields := strings.Split(datestr, "-")
		if len(fields) > 6 {
			return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second\n")
		}
		for ind, field := range fields {
			field, err := strconv.Atoi(field)
			if err != nil {
				return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second, each field must be integer\n")
			}
			dttm[ind] = field
		}
		return time.Date(dttm[0], time.Month(dttm[1]), dttm[2], dttm[3], dttm[4], dttm[5], 0, time.Local), nil
	}
	start, err := dateParser(startDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "date parse error: %v\n", err)
		os.Exit(1)
	}

	end, err := dateParser(endDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "date parse error: %v\n", err)
		os.Exit(1)
	}

	sf := make([]string, 0)
	for _, s := range strings.Split(suffixes, ",") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			sf = append(sf, strings.TrimSpace(s))
		}
	}
	bm := iqshell.GetBucketManager()
	retErr := bm.ListBucket2(bucket, prefix, listMarker, outFile, "", start, end, sf, maxRetry, appendMode)
	if retErr != nil {
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func ListBucket(cmd *cobra.Command, params []string) {
	bucket := params[0]

	bm := iqshell.GetBucketManager()
	retErr := bm.ListFiles(bucket, prefix, listMarker, outFile)
	if retErr != nil {
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Get(cmd *cobra.Command, params []string) {

	bucket := params[0]
	key := params[1]

	destFile := key
	if outFile != "" {
		destFile = outFile
	}

	bm := iqshell.GetBucketManager()
	err := bm.Get(bucket, key, destFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Stat(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]

	bm := iqshell.GetBucketManager()
	fileInfo, err := bm.Stat(bucket, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stat error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	} else {
		printStat(bucket, key, fileInfo)
	}
}

func Delete(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]

	bm := iqshell.GetBucketManager()
	err := bm.Delete(bucket, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Move(cmd *cobra.Command, params []string) {
	srcBucket := params[0]
	srcKey := params[1]
	destBucket := params[2]

	if finalKey == "" {
		finalKey = srcKey
	}

	bm := iqshell.GetBucketManager()
	err := bm.Move(srcBucket, srcKey, destBucket, finalKey, mOverwrite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Move error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Copy(cmd *cobra.Command, params []string) {
	srcBucket := params[0]
	srcKey := params[1]
	destBucket := params[2]
	if finalKey == "" {
		finalKey = srcKey
	}

	bm := iqshell.GetBucketManager()
	err := bm.Copy(srcBucket, srcKey, destBucket, finalKey, cOverwrite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Copy error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Chgm(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	newMimeType := params[2]

	bm := iqshell.GetBucketManager()
	err := bm.ChangeMime(bucket, key, newMimeType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Change mimetype error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Chtype(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	fileTypeStr := params[2]
	fileType, cErr := strconv.Atoi(fileTypeStr)
	if cErr != nil || (fileType != 0 && fileType != 1) {
		fmt.Println("Invalid file type:", fileTypeStr, ", fileType must be 0(standard) or 1(low frequency storage)")
		os.Exit(iqshell.STATUS_HALT)
		return
	}

	bm := iqshell.GetBucketManager()
	err := bm.ChangeType(bucket, key, fileType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Change file type error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func DeleteAfterDays(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	expireStr := params[2]
	expire, cErr := strconv.Atoi(expireStr)
	if cErr != nil {
		fmt.Fprintln(os.Stderr, "Invalid deleteAfterDays: ", expireStr)
		os.Exit(iqshell.STATUS_HALT)
		return
	}

	bm := iqshell.GetBucketManager()
	err := bm.DeleteAfterDays(bucket, key, expire)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Set file deleteAfterDays error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Fetch(cmd *cobra.Command, params []string) {
	remoteResUrl := params[0]
	bucket := params[1]

	var err error
	if finalKey == "" {
		finalKey, err = iqshell.KeyFromUrl(remoteResUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get key from url failed: %v\n", err)
			os.Exit(iqshell.STATUS_ERROR)
		}
	}

	bm := iqshell.GetBucketManager()
	fetchResult, err := bm.Fetch(remoteResUrl, bucket, finalKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fetch error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	} else {
		fmt.Println("Key:", fetchResult.Key)
		fmt.Println("Hash:", fetchResult.Hash)
		fmt.Printf("Fsize: %d (%s)\n", fetchResult.Fsize, FormatFsize(fetchResult.Fsize))
		fmt.Println("Mime:", fetchResult.MimeType)
	}
}

func Prefetch(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]

	bm := iqshell.GetBucketManager()
	err := bm.Prefetch(bucket, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prefetch error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func Saveas(cmd *cobra.Command, params []string) {
	publicUrl := params[0]
	saveBucket := params[1]
	saveKey := params[2]

	bm := iqshell.GetBucketManager()
	url, err := bm.Saveas(publicUrl, saveBucket, saveKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Saveas error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	} else {
		fmt.Println(url)
	}
}

func M3u8Delete(cmd *cobra.Command, params []string) {
	bucket := params[0]
	m3u8Key := params[1]

	bm := iqshell.GetBucketManager()
	m3u8FileList, err := bm.M3u8FileList(bucket, m3u8Key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get m3u8 file list error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
	entryCnt := len(m3u8FileList)
	if entryCnt == 0 {
		fmt.Fprintln(os.Stderr, "no m3u8 slices found")
		os.Exit(iqshell.STATUS_ERROR)
	}
	fileExporter, nErr := iqshell.NewFileExporter("", "", "")
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

func M3u8Replace(cmd *cobra.Command, params []string) {
	bucket := params[0]
	m3u8Key := params[1]
	var newDomain string
	if len(params) == 3 {
		newDomain = strings.TrimRight(params[2], "/")
	}

	bm := iqshell.GetBucketManager()
	err := bm.M3u8ReplaceDomain(bucket, m3u8Key, newDomain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "m3u8 replace domain error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
}

func PrivateUrl(cmd *cobra.Command, params []string) {
	publicUrl := params[0]
	var deadline int64
	if len(params) == 2 {
		if val, err := strconv.ParseInt(params[1], 10, 64); err != nil {
			fmt.Fprintln(os.Stderr, "Invalid <Deadline>")
			os.Exit(iqshell.STATUS_HALT)
		} else {
			deadline = val
		}
	} else {
		deadline = time.Now().Add(time.Second * 3600).Unix()
	}

	bm := iqshell.GetBucketManager()
	url, _ := bm.PrivateUrl(publicUrl, deadline)
	fmt.Println(url)
}

func printStat(bucket string, key string, entry storage.FileInfo) {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", entry.Fsize, FormatFsize(entry.Fsize))

	putTime := time.Unix(0, entry.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", entry.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", entry.MimeType)
	if entry.Type == 0 {
		statInfo += fmt.Sprintf("%-20s%d -> 标准存储\r\n", "FileType:", entry.Type)
	} else {
		statInfo += fmt.Sprintf("%-20s%d -> 低频存储\r\n", "FileType:", entry.Type)
	}
	fmt.Println(statInfo)
}
