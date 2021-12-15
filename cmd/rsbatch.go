package cmd

import (
	"bufio"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs/operations"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	storage2 "github.com/qiniu/qshell/v2/iqshell/storage"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
)

const (
	BATCH_ALLOW_MAX = 1000
)

func doBatchOperation(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

var (
	forceFlag     bool
	overwriteFlag bool
	worker        int
	inputFile     string
	deadline      int
	bsuccessFname string
	bfailureFname string
	sep           string
	bfetchUphost  string
)

func unescape(cmd *cobra.Command, args []string) {
	sep = utils.SimpleUnescape(&sep)
	if DebugFlag {
		fmt.Printf("forceFlag: %v, overwriteFlag: %v, worker: %v, inputFile: %q, deadline: %v, bsuccessFname: %q, bfailureFname: %q, sep: %q, bfetchUphost: %q\n", forceFlag, overwriteFlag, worker, inputFile, deadline, bsuccessFname, bfailureFname, sep, bfetchUphost)
	}
}

var batchStatCmdBuilder = func() *cobra.Command {
	var info = operations.BatchStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "batchstat <Bucket> [-i <KeyListFile>]",
		Short: "Batch stat files in bucket",
		Long:  "Batch stat files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchStatus(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchDeleteCmdBuilder = func() *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [-i <KeyListFile>]",
		Short: "Batch delete files in bucket",
		Long:  "Batch delete files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchDelete(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchChangeMimeCmdBuilder = func() *cobra.Command {
	var info = operations.BatchChangeMimeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchgm <Bucket> [-i <KeyMimeMapFile>]",
		Short: "Batch change the mime type of files in bucket",
		Long:  "Batch change the mime type of files in bucket, read from stdin if KeyMimeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeMime(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchChangeTypeCmdBuilder = func() *cobra.Command {
	var info = operations.BatchChangeTypeInfo{}
	var cmd = &cobra.Command{
		Use:   "batchchtype <Bucket> [-i <KeyFileTypeMapFile>]",
		Short: "Batch change the file type of files in bucket",
		Long:  "Batch change the file (storage) type of files in bucket, read from stdin if KeyFileTypeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchChangeType(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "delete success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "delete failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchDeleteAfterCmdBuilder = func() *cobra.Command {
	var info = operations.BatchDeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [-i <KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket",
		Long:  "Batch set the deleteAfterDays of the files in bucket, read from stdin if KeyDeleteAfterDaysMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchDeleteAfter(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchMoveCmdBuilder = func() *cobra.Command {
	var info = operations.BatchMoveInfo{}
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
			operations.BatchMove(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchRenameCmdBuilder = func() *cobra.Command {
	var info = operations.BatchRenameInfo{}
	var cmd = &cobra.Command{
		Use:   "batchrename <Bucket> [-i <OldNewKeyMapFile>]",
		Short: "Batch rename files in the bucket",
		Long:  "Batch rename files in the bucket, read from stdin if OldNewKeyMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			operations.BatchRename(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var batchCopyCmdBuilder = func() *cobra.Command {
	var info = operations.BatchCopyInfo{}
	var cmd = &cobra.Command{
		Use:   "batchcopy <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch copy files from bucket to bucket",
		Long:  "Batch copy files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.SourceBucket = args[0]
				info.DestBucket = args[1]
			}
			operations.BatchCopy(info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.InputFile, "input-file", "i", "", "input file")
	cmd.Flags().BoolVarP(&info.BatchInfo.Force, "force", "y", false, "force mode")
	cmd.Flags().BoolVarP(&info.BatchInfo.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().IntVarP(&info.BatchInfo.Worker, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "rename success list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "e", "", "rename failure list")
	cmd.Flags().StringVarP(&info.BatchInfo.ItemSeparate, "sep", "F", "\t", "Separator used for split line fields")
	return cmd
}

var (
	batchFetchCmd = &cobra.Command{
		Use:   "batchfetch <Bucket> [-i <FetchUrlsFile>] [-c <WorkerCount>]",
		Short: "Batch fetch remoteUrls and save them in qiniu Bucket",
		Args:  cobra.ExactArgs(1),
		Run:   BatchFetch,
	}
	batchSignCmd = &cobra.Command{
		Use:   "batchsign [-i <ItemListFile>] [-e <Deadline>]",
		Short: "Batch create the private url from the public url list file",
		Args:  cobra.ExactArgs(0),
		Run:   BatchSign,
	}
)

func init() {
	batchFetchCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "urls list file")
	batchFetchCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchFetchCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "file to save batch fetch success list")
	batchFetchCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "file to save batch fetch failure list")
	batchFetchCmd.Flags().StringVarP(&bfetchUphost, "up-host", "u", "", "fetch uphost")
	batchFetchCmd.Flags().StringVarP(&sep, "sep", "F", "\t", "Separator used for split line fields")

	batchSignCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchSignCmd.Flags().IntVarP(&deadline, "deadline", "e", 3600, "deadline in seconds")
	batchSignCmd.Flags().StringVarP(&sep, "sep", "F", "\t", "Separator used for split line fields")

	cmds := []*cobra.Command{
		batchStatCmdBuilder(),
		batchCopyCmdBuilder(),
		batchMoveCmdBuilder(),
		batchRenameCmdBuilder(),
		batchDeleteCmdBuilder(),
		batchDeleteAfterCmdBuilder(),
		batchChangeMimeCmdBuilder(),
		batchChangeTypeCmdBuilder(),
		batchSignCmd, batchFetchCmd,
	}
	RootCmd.AddCommand(cmds...)
	for _, cmd := range cmds {
		cmd.PersistentPreRun = unescape
	}
}

type fetchConfig struct {
	upHost string

	threadCount int

	successFname   string
	failureFname   string
	overwriteFname string

	fileExporter *storage2.FileExporter
	bm           *storage2.BucketManager
}

// initFileExporter需要在主goroutine中调用， 原因同initBucketManager
func (fc *fetchConfig) initFileExporter() {
	fileExporter, fErr := storage2.NewFileExporter(fc.successFname, fc.failureFname, "")
	if fErr != nil {
		fmt.Fprintf(os.Stderr, "create file exporter: %v\n", fErr)
		os.Exit(1)
	}
	fc.fileExporter = fileExporter
}

// GetBucketManagerWithConfig 会使用os.Exit推出，因此该方法需要在main gouroutine中调用
func (fc *fetchConfig) initBucketManager() {

	cfg := workspace.GetConfig()
	region := (&cfg).GetRegion()
	if len(fc.upHost) > 0 {
		region.SrcUpHosts = []string{fc.upHost}
		region.CdnUpHosts = nil
	}

	fc.bm = storage2.GetBucketManagerWithConfig(&storage.Config{
		Region: region,
	})
}

// initUpHost需要在主goroutine中调用
func (fc *fetchConfig) initUpHost(bucket string) {
	if bfetchUphost == "" {
		acc, aerr := account.GetAccount()
		if aerr != nil {
			fmt.Fprintf(os.Stderr, "failed to get accessKey")
			os.Exit(1)
		}
		region, rErr := storage.GetRegion(acc.AccessKey, bucket)
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "failed getting fetch host for bucket: %s: %v\n", bucket, rErr)
			os.Exit(1)
		}
		bfetchUphost = region.IovipHost
	}
	fc.upHost = bfetchUphost
}

// 批量抓取网络资源到七牛存储空间
func BatchFetch(cmd *cobra.Command, params []string) {
	if worker <= 0 || worker > 1000 {
		fmt.Fprintf(os.Stderr, "threads count: %d is too large, must be (0, 1000]", worker)
		os.Exit(1)
	}
	bucket := params[0]
	var urlsListFile string

	if inputFile == "" {
		urlsListFile = "stdin"
	} else {
		urlsListFile = inputFile
	}
	var fp io.ReadCloser
	var err error

	if urlsListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(urlsListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open urls list file: %s : %v\n", urlsListFile, err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)

	var (
		saveKey   string
		remoteUrl string
		pError    error
		fItemChan chan *data.FetchItem = make(chan *data.FetchItem)
	)
	defer close(fItemChan)

	itemc := make(chan *data.FetchItem)
	donec := make(chan struct{})

	fconfig := fetchConfig{
		threadCount:  worker,
		successFname: bsuccessFname,
		failureFname: bfailureFname,
	}

	fconfig.initUpHost(bucket)
	fconfig.initBucketManager()
	fconfig.initFileExporter()

	go fetchChannel(itemc, donec, &fconfig)

	for scanner.Scan() {
		line := scanner.Text()
		items := utils.SplitString(line, sep)
		if len(items) <= 0 {
			continue
		}
		remoteUrl = items[0]
		if remoteUrl == "" {
			continue
		}
		if len(items) <= 1 {
			saveKey, pError = utils.KeyFromUrl(remoteUrl)
			if pError != nil {
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", remoteUrl, pError)
				continue
			}
		} else {
			saveKey = items[1]
		}
		item := data.FetchItem{
			Bucket:    bucket,
			Key:       saveKey,
			RemoteUrl: remoteUrl,
		}
		itemc <- &item
	}
	close(itemc)

	<-donec
}

func fetchChannel(c chan *data.FetchItem, donec chan struct{}, fconfig *fetchConfig) {

	fileExporter := fconfig.fileExporter
	bm := fconfig.bm

	limitc := make(chan struct{}, fconfig.threadCount)
	wg := sync.WaitGroup{}

	for item := range c {
		limitc <- struct{}{}
		wg.Add(1)

		go func(item *data.FetchItem) {
			_, fErr := bm.Fetch(item.RemoteUrl, item.Bucket, item.Key)
			if fErr != nil {
				fmt.Fprintf(os.Stderr, "fetch %s => %s:%s failed\n", item.RemoteUrl, item.Bucket, item.Key)
				if fileExporter != nil {
					fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%s\t%v\n", item.RemoteUrl, item.Key, fErr))
				}
			} else {
				fmt.Printf("fetch %s => %s:%s success\n", item.RemoteUrl, item.Bucket, item.Key)
				if fileExporter != nil {
					fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", item.RemoteUrl, item.Key))
				}
			}
			<-limitc
			wg.Done()
		}(item)
	}
	wg.Wait()

	donec <- struct{}{}
}

// 批量获取文件列表的信息
func BatchStat(cmd *cobra.Command, params []string) {
	bucket := params[0]

	var keyListFile string

	if inputFile == "" {
		keyListFile = "stdin"
	} else {
		keyListFile = inputFile
	}

	var fp io.ReadCloser
	var err error

	if keyListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open key list file: %s, error: %v\n", keyListFile, err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}

	bm := storage2.GetBucketManager()
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, sep)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := storage2.EntryPath{
					Bucket: bucket,
					Key:    key,
				}
				entries = append(entries, entry)
			}
		}
		//check 1000 limit
		if len(entries) == BATCH_ALLOW_MAX {
			batchStat(entries, bm)
			//reset slice
			entries = make([]storage2.EntryPath, 0)
		}
	}
	//stat the last batch
	if len(entries) > 0 {
		batchStat(entries, bm)
	}
}

func batchStat(entries []storage2.EntryPath, bm *storage2.BucketManager) {
	ret, err := bm.BatchStat(entries)
	if err != nil && len(ret) <= 0 {
		fmt.Fprintf(os.Stderr, "Batch stat error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fmt.Fprintln(os.Stderr, entry.Key+"\t"+item.Data.Error)
			} else {
				fmt.Println(fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d", entry.Key,
					item.Data.Fsize, item.Data.Hash, item.Data.MimeType, item.Data.PutTime, item.Data.Type))
			}
		}
	}
}

// 批量删除七牛存储空间中的文件
func BatchDelete(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
		} else {
			fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyListFile string

	if inputFile == "" {
		keyListFile = "stdin"
	} else {
		keyListFile = inputFile
	}

	bm := storage2.GetBucketManager()

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	var fp io.ReadCloser
	var err error

	if keyListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyListFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open key list file error", err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.EntryPath, 0, BATCH_ALLOW_MAX)
	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, sep)
		if len(items) <= 0 {
			continue
		}
		key := items[0]
		if key != "" {
			putTime := ""
			if len(items) > 1 {
				putTime = items[1]
			}
			entry := storage2.EntryPath{
				Bucket: bucket, Key: key, PutTime: putTime,
			}
			entries = append(entries, entry)
		}
		//check limit
		if len(entries) == BATCH_ALLOW_MAX {
			toDeleteEntries := make([]storage2.EntryPath, len(entries))
			copy(toDeleteEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDelete(toDeleteEntries, bm, fileExporter)
			}
			entries = make([]storage2.EntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	//delete the last batch
	if len(entries) > 0 {
		toDeleteEntries := make([]storage2.EntryPath, len(entries))
		copy(toDeleteEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDelete(toDeleteEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchDelete(entries []storage2.EntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchDelete(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch delete error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]

			if item.Code != 200 || item.Data.Error != "" {
				fmt.Fprintf(os.Stderr, "Delete '%s' => '%s' failed, Code: %d, Error: %s\n", entry.Bucket, entry.Key, item.Code, item.Data.Error)
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
			} else {
				fmt.Printf("Delete '%s' => '%s' success\n", entry.Bucket, entry.Key)
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
			}
		}
	}
}

// 批量修改存储在七牛存储空间中文件的MimeType信息
func BatchChgm(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyMimeMapFile string
	if inputFile == "" {
		keyMimeMapFile = "stdin"
	} else {
		keyMimeMapFile = inputFile
	}

	bm := storage2.GetBucketManager()
	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})
	var fp io.ReadCloser
	var err error
	if keyMimeMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyMimeMapFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open key mime map file error: %v\n", err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]storage2.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, sep)
		if len(items) == 2 {
			key := items[0]
			mimeType := items[1]
			if key != "" && mimeType != "" {
				entry := storage2.ChgmEntryPath{
					EntryPath: storage2.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					MimeType: mimeType,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChgmEntries := make([]storage2.ChgmEntryPath, len(entries))
			copy(toChgmEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChgm(toChgmEntries, bm, fileExporter)
			}
			entries = make([]storage2.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChgmEntries := make([]storage2.ChgmEntryPath, len(entries))
		copy(toChgmEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChgm(toChgmEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchChgm(entries []storage2.ChgmEntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchChgm(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chgm error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
				fmt.Fprintf(os.Stderr, "Chgm '%s' => '%s' Failed, Code: %d, Error: %s\n", entry.Key, entry.MimeType, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
				fmt.Printf("Chgm '%s' => '%s' success\n", entry.Key, entry.MimeType)
			}
		}
	}
}

// 批量修改存储在七牛存储空间中文件的存储类型信息（标准存储-》低频存储，低频-》标准存储)
func BatchChtype(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyTypeMapFile string
	if inputFile == "" {
		keyTypeMapFile = "stdin"
	} else {
		keyTypeMapFile = inputFile
	}
	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	bm := storage2.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyTypeMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyTypeMapFile)
		if err != nil {
			fmt.Printf("Open key file type map file error: %v\n", err)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)

	var key, line string
	var fileType int
	var items []string
	var entry storage2.ChtypeEntryPath

	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line = scanner.Text()
		items = strings.Split(line, sep)

		if len(items) == 2 {
			fileType, _ = strconv.Atoi(items[1])
		} else if len(items) == 1 {
			fileType = 1
		}
		key = items[0]
		if key != "" {
			entry = storage2.ChtypeEntryPath{
				EntryPath: storage2.EntryPath{
					Bucket: bucket,
					Key:    key,
				},
				FileType: fileType,
			}
			entries = append(entries, entry)
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChtypeEntries := make([]storage2.ChtypeEntryPath, len(entries))
			copy(toChtypeEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChtype(toChtypeEntries, bm, fileExporter)
			}
			entries = make([]storage2.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChtypeEntries := make([]storage2.ChtypeEntryPath, len(entries))
		copy(toChtypeEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChtype(toChtypeEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()

}

func batchChtype(entries []storage2.ChtypeEntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchChtype(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chtype error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
				fmt.Fprintf(os.Stderr, "Chtype '%s' => '%d' Failed, Code: %d, Error: %s\n", entry.Key, entry.FileType, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
				fmt.Printf("Chtype '%s' => '%d' success\n", entry.Key, entry.FileType)
			}
		}
	}
	return
}

// 批量设置七牛存储空间中的删除标志（多少天后删除）
func BatchDeleteAfterDays(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	bucket := params[0]
	var keyExpireMapFile string

	if inputFile == "" {
		keyExpireMapFile = "stdin"
	} else {
		keyExpireMapFile = inputFile
	}

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	bm := storage2.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyExpireMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyExpireMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open key expire map file error")
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, sep)
		if len(items) == 2 {
			key := items[0]
			days, _ := strconv.Atoi(items[1])
			if key != "" {
				entry := storage2.DeleteAfterDaysEntryPath{
					EntryPath: storage2.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					DeleteAfterDays: days,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toExpireEntries := make([]storage2.DeleteAfterDaysEntryPath, len(entries))
			copy(toExpireEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDeleteAfterDays(toExpireEntries, bm)
			}
			entries = make([]storage2.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toExpireEntries := make([]storage2.DeleteAfterDaysEntryPath, len(entries))
		copy(toExpireEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDeleteAfterDays(toExpireEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchDeleteAfterDays(entries []storage2.DeleteAfterDaysEntryPath, bm *storage2.BucketManager) {
	ret, err := bm.BatchDeleteAfterDays(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch expire error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fmt.Fprintf(os.Stderr, "Expire '%s' => '%d' Failed, Code: %d, Error: %s\n", entry.Key, entry.DeleteAfterDays, item.Code, item.Data.Error)
			} else {
				fmt.Printf("Expire '%s' => '%d' success\n", entry.Key, entry.DeleteAfterDays)
			}
		}
	}
}

// 批量重命名七牛存储空间中的文件
func BatchRename(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	bucket := params[0]
	var oldNewKeyMapFile string

	if inputFile == "" {
		oldNewKeyMapFile = "stdin"
	} else {
		oldNewKeyMapFile = inputFile
	}

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	bm := storage2.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if oldNewKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open old new key map file error")
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.RenameEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := utils.SplitString(line, sep)
		if len(items) == 2 {
			oldKey := items[0]
			newKey := items[1]
			if oldKey != "" && newKey != "" {
				entry := storage2.RenameEntryPath{
					SrcEntry: storage2.EntryPath{
						Bucket: bucket,
						Key:    oldKey,
					},
					DstEntry: storage2.EntryPath{
						Bucket: bucket,
						Key:    newKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toRenameEntries := make([]storage2.RenameEntryPath, len(entries))
			copy(toRenameEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchRename(toRenameEntries, bm, fileExporter)
			}
			entries = make([]storage2.RenameEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toRenameEntries := make([]storage2.RenameEntryPath, len(entries))
		copy(toRenameEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchRename(toRenameEntries, bm, fileExporter)
		}
	}
	batchWaitGroup.Wait()
}

func batchRename(entries []storage2.RenameEntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchRename(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch rename error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%s\t%d\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key, item.Code, item.Data.Error))
				fmt.Fprintf(os.Stderr, "Rename '%s' => '%s' Failed, Code: %d, Error: %s\n", entry.SrcEntry.Key, entry.DstEntry.Key, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				fmt.Printf("Rename '%s' => '%s' success\n", entry.SrcEntry.Key, entry.DstEntry.Key)
			}
		}
	}
}

// 批量移动七牛存储空间中的文件
func BatchMove(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(data.STATUS_HALT)
		}
	}

	srcBucket := params[0]
	destBucket := params[1]
	var srcDestKeyMapFile string

	if inputFile == "" {
		srcDestKeyMapFile = "stdin"
	} else {
		srcDestKeyMapFile = inputFile
	}

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	bm := storage2.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]storage2.MoveEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := utils.SplitString(line, sep)
		if len(items) == 1 || len(items) == 2 {
			srcKey := items[0]
			destKey := srcKey
			if len(items) == 2 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				entry := storage2.MoveEntryPath{
					SrcEntry: storage2.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: storage2.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toMoveEntries := make([]storage2.MoveEntryPath, len(entries))
			copy(toMoveEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchMove(toMoveEntries, bm, fileExporter)
			}
			entries = make([]storage2.MoveEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toMoveEntries := make([]storage2.MoveEntryPath, len(entries))
		copy(toMoveEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchMove(toMoveEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchMove(entries []storage2.MoveEntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchMove(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch move error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%s\t%d\t%s\n", entry.SrcEntry.Key,
					entry.DstEntry.Key, item.Code, item.Data.Error))
				fmt.Fprintf(os.Stderr, "Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				fmt.Printf("Move '%s:%s' => '%s:%s' success\n",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key, entry.DstEntry.Bucket, entry.DstEntry.Key)
			}
		}
	}
}

// 批量拷贝七牛存储中的文件
func BatchCopy(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := utils.CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Verification code is not valid")
			os.Exit(data.STATUS_HALT)
		}
	}

	srcBucket := params[0]
	destBucket := params[1]

	var srcDestKeyMapFile string

	if inputFile == "" {
		srcDestKeyMapFile = "stdin"
	} else {
		srcDestKeyMapFile = inputFile
	}

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	bm := storage2.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]storage2.CopyEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := storage2.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := utils.SplitString(line, sep)
		if len(items) == 1 || len(items) == 2 {
			srcKey := items[0]
			destKey := srcKey
			if len(items) == 2 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				entry := storage2.CopyEntryPath{
					SrcEntry: storage2.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: storage2.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toCopyEntries := make([]storage2.CopyEntryPath, len(entries))
			copy(toCopyEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchCopy(toCopyEntries, bm, fileExporter)
			}
			entries = make([]storage2.CopyEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toCopyEntries := make([]storage2.CopyEntryPath, len(entries))
		copy(toCopyEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchCopy(toCopyEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchCopy(entries []storage2.CopyEntryPath, bm *storage2.BucketManager, fileExporter *storage2.FileExporter) {
	ret, err := bm.BatchCopy(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch copy error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%s\t%d\t%s\n", entry.SrcEntry.Key,
					entry.DstEntry.Key, item.Code, item.Data.Error))
				fmt.Fprintf(os.Stderr, "Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s\n",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				fmt.Printf("Copy '%s:%s' => '%s:%s' success\n",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key)
			}
		}
	}
}

// 批量签名存储空间中的文件
func BatchSign(cmd *cobra.Command, params []string) {
	if deadline <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid <Deadline>: deadline must be int and greater than 0\n")
		os.Exit(data.STATUS_HALT)
	}
	d := time.Now().Add(time.Second * time.Duration(deadline) * 24 * 365).Unix()

	var bReader io.Reader

	bm := storage2.GetBucketManager()

	if inputFile != "" {
		fp, openErr := os.Open(inputFile)
		if openErr != nil {
			fmt.Fprintln(os.Stderr, "Open url list file error,", openErr)
			os.Exit(data.STATUS_HALT)
		}
		defer fp.Close()
		bReader = fp
	} else {
		bReader = os.Stdin
	}

	var url string
	scanner := bufio.NewScanner(bReader)
	for scanner.Scan() {
		line := scanner.Text()
		items := utils.SplitString(line, sep)
		if len(items) <= 0 {
			continue
		}
		url = items[0]
		if url == "" {
			continue
		}
		urlToSign := strings.TrimSpace(url)
		if urlToSign == "" {
			continue
		}
		signedUrl, _ := bm.PrivateUrl(urlToSign, d)
		fmt.Println(signedUrl)
	}
}
