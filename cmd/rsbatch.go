package cmd

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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
)

var (
	batchFetchCmd = &cobra.Command{
		Use:   "batchfetch <Bucket> [-i <FetchUrlsFile>] [-c <WorkerCount>]",
		Short: "Batch fetch remoteUrls and save them in qiniu Bucket",
		Args:  cobra.ExactArgs(1),
		Run:   BatchFetch,
	}
	batchStatCmd = &cobra.Command{
		Use:   "batchstat <Bucket> [-i <KeyListFile>]",
		Short: "Batch stat files in bucket",
		Long:  "Batch stat files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchStat,
	}
	batchDeleteCmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [-i <KeyListFile>]",
		Short: "Batch delete files in bucket",
		Long:  "Batch delete files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchDelete,
	}
	batchChgmCmd = &cobra.Command{
		Use:   "batchchgm <Bucket> [-i <KeyMimeMapFile>]",
		Short: "Batch change the mime type of files in bucket",
		Long:  "Batch change the mime type of files in bucket, read from stdin if KeyMimeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchChgm,
	}
	batchChtypeCmd = &cobra.Command{
		Use:   "batchchtype <Bucket> [-i <KeyFileTypeMapFile>]",
		Short: "Batch change the file type of files in bucket",
		Long:  "Batch change the file (storage) type of files in bucket, read from stdin if KeyFileTypeMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchChtype,
	}
	batchDelAfterCmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [-i <KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket",
		Long:  "Batch set the deleteAfterDays of the files in bucket, read from stdin if KeyDeleteAfterDaysMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchDeleteAfterDays,
	}
	batchRenameCmd = &cobra.Command{
		Use:   "batchrename <Bucket> [-i <OldNewKeyMapFile>]",
		Short: "Batch rename files in the bucket",
		Long:  "Batch rename files in the bucket, read from stdin if OldNewKeyMapFile not specified",
		Args:  cobra.ExactArgs(1),
		Run:   BatchRename,
	}
	batchMoveCmd = &cobra.Command{
		Use:   "batchmove <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch move files from bucket to bucket",
		Long:  "Batch move files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.ExactArgs(2),
		Run:   BatchMove,
	}
	batchCopyCmd = &cobra.Command{
		Use:   "batchcopy <SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]",
		Short: "Batch copy files from bucket to bucket",
		Long:  "Batch copy files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.ExactArgs(2),
		Run:   BatchCopy,
	}
	batchSignCmd = &cobra.Command{
		Use:   "batchsign [-i <UrlListFile>] [-e <Deadline>]",
		Short: "Batch create the private url from the public url list file",
		Args:  cobra.ExactArgs(0),
		Run:   BatchSign,
	}
)

func init() {
	batchFetchCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "urls list file")
	batchFetchCmd.Flags().IntVarP(&worker, "worker", "c", 5, "worker count")
	batchStatCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchCopyCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchMoveCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchRenameCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchDeleteCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDeleteCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchDeleteCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchDeleteCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "delete success list")
	batchDeleteCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "delete failure list")

	batchChgmCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChgmCmd.Flags().IntVarP(&worker, "worker", "c", 1, "woker count")
	batchChgmCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchChgmCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "change mimetype success list")
	batchChgmCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "change mimetype failure list")

	batchChtypeCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChtypeCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchChtypeCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchChtypeCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "change storage type success file list")
	batchChtypeCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "change storage type failure file list")

	batchDelAfterCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDelAfterCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchDelAfterCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchRenameCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchRenameCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchRenameCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchRenameCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "rename success list")
	batchRenameCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "rename failure list")

	batchMoveCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchMoveCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchMoveCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchMoveCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "move success list")
	batchMoveCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "move failure list")

	batchCopyCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchCopyCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchCopyCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchCopyCmd.Flags().StringVarP(&bsuccessFname, "success-list", "s", "", "copy success list")
	batchCopyCmd.Flags().StringVarP(&bfailureFname, "failure-list", "e", "", "copy failure list")

	batchSignCmd.Flags().IntVarP(&deadline, "deadline", "e", 3600, "deadline in seconds")
	batchSignCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	RootCmd.AddCommand(batchStatCmd, batchDeleteCmd, batchChgmCmd, batchChtypeCmd, batchDelAfterCmd,
		batchRenameCmd, batchMoveCmd, batchCopyCmd, batchSignCmd, batchFetchCmd)
}

func BatchFetch(cmd *cobra.Command, params []string) {
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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)

	var (
		saveKey   string
		remoteUrl string
		pError    error
		fItemChan chan *iqshell.FetchItem = make(chan *iqshell.FetchItem)
	)
	defer close(fItemChan)

	go batchFetch(fItemChan)

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) <= 0 {
			continue
		}
		remoteUrl = items[0]
		if remoteUrl == "" {
			continue
		}
		if len(items) <= 1 {
			saveKey, pError = iqshell.KeyFromUrl(remoteUrl)
			if pError != nil {
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", remoteUrl, pError)
				continue
			}
		} else {
			saveKey = items[1]
		}
		item := iqshell.FetchItem{
			Bucket:    bucket,
			Key:       saveKey,
			RemoteUrl: remoteUrl,
		}
		fItemChan <- &item
	}
}

func batchFetch(fItemChan chan *iqshell.FetchItem) {
	for i := 0; i < worker; i++ {
		go func() {
			bm := iqshell.GetBucketManager()
			for fetchItem := range fItemChan {
				fetchResult, fErr := bm.Fetch(fetchItem.RemoteUrl, fetchItem.Bucket, fetchItem.Key)
				fmt.Println(fetchResult, fErr)
			}
		}()
	}
}

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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}

	bm := iqshell.GetBucketManager()
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := iqshell.EntryPath{
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
			entries = make([]iqshell.EntryPath, 0)
		}
	}
	//stat the last batch
	if len(entries) > 0 {
		batchStat(entries, bm)
	}
}

func batchStat(entries []iqshell.EntryPath, bm *iqshell.BucketManager) {
	ret, err := bm.BatchStat(entries)
	if err != nil {
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

func BatchDelete(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
		} else {
			fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyListFile string

	if inputFile == "" {
		keyListFile = "stdin"
	} else {
		keyListFile = inputFile
	}

	bm := iqshell.GetBucketManager()

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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.EntryPath, 0, BATCH_ALLOW_MAX)
	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := iqshell.EntryPath{
					bucket, key,
				}
				entries = append(entries, entry)
			}
		}
		//check limit
		if len(entries) == BATCH_ALLOW_MAX {
			toDeleteEntries := make([]iqshell.EntryPath, len(entries))
			copy(toDeleteEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDelete(toDeleteEntries, bm, fileExporter)
			}
			entries = make([]iqshell.EntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	//delete the last batch
	if len(entries) > 0 {
		toDeleteEntries := make([]iqshell.EntryPath, len(entries))
		copy(toDeleteEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDelete(toDeleteEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchDelete(entries []iqshell.EntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
	ret, err := bm.BatchDelete(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch delete error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]

			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Delete '%s' => '%s' failed, Code: %d, Error: %s", entry.Bucket, entry.Key, item.Code, item.Data.Error)
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
			} else {
				logs.Debug("Delete '%s' => '%s' success", entry.Bucket, entry.Key)
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
			}
		}
	}
}

func BatchChgm(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyMimeMapFile string
	if inputFile == "" {
		keyMimeMapFile = "stdin"
	} else {
		keyMimeMapFile = inputFile
	}

	bm := iqshell.GetBucketManager()
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
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]iqshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) == 2 {
			key := items[0]
			mimeType := items[1]
			if key != "" && mimeType != "" {
				entry := iqshell.ChgmEntryPath{
					EntryPath: iqshell.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					MimeType: mimeType,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChgmEntries := make([]iqshell.ChgmEntryPath, len(entries))
			copy(toChgmEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChgm(toChgmEntries, bm, fileExporter)
			}
			entries = make([]iqshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChgmEntries := make([]iqshell.ChgmEntryPath, len(entries))
		copy(toChgmEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChgm(toChgmEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchChgm(entries []iqshell.ChgmEntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
	ret, err := bm.BatchChgm(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chgm error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
				logs.Error("Chgm '%s' => '%s' Failed, Code: %d, Error: %s", entry.Key, entry.MimeType, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
				logs.Debug("Chgm '%s' => '%s' success", entry.Key, entry.MimeType)
			}
		}
	}
}

func BatchChtype(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
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

	bm := iqshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyTypeMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyTypeMapFile)
		if err != nil {
			fmt.Printf("Open key file type map file error: %v\n", err)
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)

	var key, line string
	var fileType int
	var items []string
	var entry iqshell.ChtypeEntryPath

	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line = scanner.Text()
		items = strings.Fields(line)

		if len(items) == 2 {
			fileType, _ = strconv.Atoi(items[1])
		} else if len(items) == 1 {
			fileType = 1
		}
		key = items[0]
		if key != "" {
			entry = iqshell.ChtypeEntryPath{
				EntryPath: iqshell.EntryPath{
					Bucket: bucket,
					Key:    key,
				},
				FileType: fileType,
			}
			entries = append(entries, entry)
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChtypeEntries := make([]iqshell.ChtypeEntryPath, len(entries))
			copy(toChtypeEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChtype(toChtypeEntries, bm, fileExporter)
			}
			entries = make([]iqshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChtypeEntries := make([]iqshell.ChtypeEntryPath, len(entries))
		copy(toChtypeEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChtype(toChtypeEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()

}

func batchChtype(entries []iqshell.ChtypeEntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
	ret, err := bm.BatchChtype(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chtype error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", entry.Key, item.Code, item.Data.Error))
				logs.Error("Chtype '%s' => '%d' Failed, Code: %d, Error: %s\n", entry.Key, entry.FileType, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\n", entry.Key))
				logs.Debug("Chtype '%s' => '%s' success", entry.Key, entry.FileType)
			}
		}
	}
	return
}

func BatchDeleteAfterDays(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
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

	bm := iqshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyExpireMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyExpireMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open key expire map file error")
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			key := items[0]
			days, _ := strconv.Atoi(items[1])
			if key != "" {
				entry := iqshell.DeleteAfterDaysEntryPath{
					EntryPath: iqshell.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					DeleteAfterDays: days,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toExpireEntries := make([]iqshell.DeleteAfterDaysEntryPath, len(entries))
			copy(toExpireEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDeleteAfterDays(toExpireEntries, bm)
			}
			entries = make([]iqshell.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toExpireEntries := make([]iqshell.DeleteAfterDaysEntryPath, len(entries))
		copy(toExpireEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDeleteAfterDays(toExpireEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchDeleteAfterDays(entries []iqshell.DeleteAfterDaysEntryPath, bm *iqshell.BucketManager) {
	ret, err := bm.BatchDeleteAfterDays(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch expire error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Expire '%s' => '%d' Failed, Code: %d, Error: %s", entry.Key, entry.DeleteAfterDays, item.Code, item.Data.Error)
			} else {
				logs.Debug("Expire '%s' => '%d' success", entry.Key, entry.DeleteAfterDays)
			}
		}
	}
}

func BatchRename(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
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

	bm := iqshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if oldNewKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open old new key map file error")
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			oldKey := items[0]
			newKey := items[1]
			if oldKey != "" && newKey != "" {
				entry := iqshell.RenameEntryPath{
					SrcEntry: iqshell.EntryPath{
						Bucket: bucket,
						Key:    oldKey,
					},
					DstEntry: iqshell.EntryPath{
						Bucket: bucket,
						Key:    newKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toRenameEntries := make([]iqshell.RenameEntryPath, len(entries))
			copy(toRenameEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchRename(toRenameEntries, bm, fileExporter)
			}
			entries = make([]iqshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toRenameEntries := make([]iqshell.RenameEntryPath, len(entries))
		copy(toRenameEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchRename(toRenameEntries, bm, fileExporter)
		}
	}
	batchWaitGroup.Wait()
}

func batchRename(entries []iqshell.RenameEntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
	ret, err := bm.BatchRename(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch rename error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%s\t%d\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key, item.Code, item.Data.Error))
				logs.Error("Rename '%s' => '%s' Failed, Code: %d, Error: %s", entry.SrcEntry.Key, entry.DstEntry.Key, item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				logs.Debug("Rename '%s' => '%s' success", entry.SrcEntry.Key, entry.DstEntry.Key)
			}
		}
	}
}

func BatchMove(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Task quit!")
			os.Exit(iqshell.STATUS_HALT)
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

	bm := iqshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]iqshell.MoveEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 1 || len(items) == 2 {
			srcKey := items[0]
			destKey := srcKey
			if len(items) == 2 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				entry := iqshell.MoveEntryPath{
					SrcEntry: iqshell.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: iqshell.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toMoveEntries := make([]iqshell.MoveEntryPath, len(entries))
			copy(toMoveEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchMove(toMoveEntries, bm, fileExporter)
			}
			entries = make([]iqshell.MoveEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toMoveEntries := make([]iqshell.MoveEntryPath, len(entries))
		copy(toMoveEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchMove(toMoveEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchMove(entries []iqshell.MoveEntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
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
				logs.Error("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				logs.Debug("Move '%s:%s' => '%s:%s' success",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key, entry.DstEntry.Bucket, entry.DstEntry.Key)
			}
		}
	}
}

func BatchCopy(cmd *cobra.Command, params []string) {
	if !forceFlag {
		//confirm
		rcode := CreateRandString(6)

		rcode2 := ""
		if runtime.GOOS == "windows" {
			fmt.Printf("<DANGER> Input %s to confirm operation: ", rcode)
		} else {
			fmt.Printf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode)
		}
		fmt.Scanln(&rcode2)

		if rcode != rcode2 {
			fmt.Fprintln(os.Stderr, "Verification code is not valid")
			os.Exit(iqshell.STATUS_HALT)
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

	bm := iqshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]iqshell.CopyEntryPath, 0, BATCH_ALLOW_MAX)

	fileExporter, nErr := iqshell.NewFileExporter(bsuccessFname, bfailureFname, "")
	if nErr != nil {
		fmt.Fprintf(os.Stderr, "create FileExporter: %v\n", nErr)
		os.Exit(1)
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 1 || len(items) == 2 {
			srcKey := items[0]
			destKey := srcKey
			if len(items) == 2 {
				destKey = items[1]
			}
			if srcKey != "" && destKey != "" {
				entry := iqshell.CopyEntryPath{
					SrcEntry: iqshell.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: iqshell.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toCopyEntries := make([]iqshell.CopyEntryPath, len(entries))
			copy(toCopyEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchCopy(toCopyEntries, bm, fileExporter)
			}
			entries = make([]iqshell.CopyEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toCopyEntries := make([]iqshell.CopyEntryPath, len(entries))
		copy(toCopyEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchCopy(toCopyEntries, bm, fileExporter)
		}
	}

	batchWaitGroup.Wait()
}

func batchCopy(entries []iqshell.CopyEntryPath, bm *iqshell.BucketManager, fileExporter *iqshell.FileExporter) {
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
				logs.Error("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
				fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", entry.SrcEntry.Key, entry.DstEntry.Key))
				logs.Debug("Copy '%s:%s' => '%s:%s' success",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key)
			}
		}
	}
}

func BatchSign(cmd *cobra.Command, params []string) {
	if deadline <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid <Deadline>: deadline must be int and greater than 0\n")
		os.Exit(iqshell.STATUS_HALT)
	}
	d := time.Now().Add(time.Second * time.Duration(deadline) * 24 * 365).Unix()

	var bReader io.Reader

	bm := iqshell.GetBucketManager()

	if inputFile != "" {
		fp, openErr := os.Open(inputFile)
		if openErr != nil {
			fmt.Fprintln(os.Stderr, "Open url list file error,", openErr)
			os.Exit(iqshell.STATUS_HALT)
		}
		defer fp.Close()
		bReader = fp
	} else {
		bReader = os.Stdin
	}

	scanner := bufio.NewScanner(bReader)
	for scanner.Scan() {
		urlToSign := strings.TrimSpace(scanner.Text())
		if urlToSign == "" {
			continue
		}
		signedUrl, _ := bm.PrivateUrl(urlToSign, d)
		fmt.Println(signedUrl)
	}
}
