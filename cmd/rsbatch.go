package cmd

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
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
)

var (
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
		Use:   "batchsign <UrlListFile> [<Deadline>]",
		Short: "Batch create the private url from the public url list file",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchSign,
	}
)

func init() {
	batchStatCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchCopyCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchMoveCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	batchRenameCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchDeleteCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDeleteCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchDeleteCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchChgmCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChgmCmd.Flags().IntVarP(&worker, "worker", "c", 1, "woker count")
	batchChgmCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchChtypeCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChtypeCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchChtypeCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchDelAfterCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDelAfterCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")
	batchDelAfterCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")

	batchRenameCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchRenameCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchRenameCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

	batchMoveCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchMoveCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchMoveCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

	batchCopyCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchCopyCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "w", false, "overwrite mode")
	batchCopyCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

	RootCmd.AddCommand(batchStatCmd, batchDeleteCmd, batchChgmCmd, batchChtypeCmd, batchDelAfterCmd,
		batchRenameCmd, batchMoveCmd, batchCopyCmd, batchSignCmd)
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
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}

	bm := qshell.GetBucketManager()
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := qshell.EntryPath{
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
			entries = make([]qshell.EntryPath, 0)
		}
	}
	//stat the last batch
	if len(entries) > 0 {
		batchStat(entries, bm)
	}
}

func batchStat(entries []qshell.EntryPath, bm *qshell.BucketManager) {
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
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyListFile string

	if inputFile == "" {
		keyListFile = "stdin"
	} else {
		keyListFile = inputFile
	}

	bm := qshell.GetBucketManager()

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
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := qshell.EntryPath{
					bucket, key,
				}
				entries = append(entries, entry)
			}
		}
		//check limit
		if len(entries) == BATCH_ALLOW_MAX {
			toDeleteEntries := make([]qshell.EntryPath, len(entries))
			copy(toDeleteEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDelete(toDeleteEntries, bm)
			}
			entries = make([]qshell.EntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	//delete the last batch
	if len(entries) > 0 {
		toDeleteEntries := make([]qshell.EntryPath, len(entries))
		copy(toDeleteEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDelete(toDeleteEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchDelete(entries []qshell.EntryPath, bm *qshell.BucketManager) {
	ret, err := bm.BatchDelete(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch delete error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]

			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Delete '%s' => '%s' failed, Code: %d, Error: %s", entry.Bucket, entry.Key, item.Code, item.Data.Error)
			} else {
				logs.Debug("Delete '%s' => '%s' success", entry.Bucket, entry.Key)
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
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyMimeMapFile string
	if inputFile == "" {
		keyMimeMapFile = "stdin"
	} else {
		keyMimeMapFile = inputFile
	}

	bm := qshell.GetBucketManager()
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
			fmt.Printf("Open key mime map file error: %v\n", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]qshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) == 2 {
			key := items[0]
			mimeType := items[1]
			if key != "" && mimeType != "" {
				entry := qshell.ChgmEntryPath{
					EntryPath: qshell.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					MimeType: mimeType,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChgmEntries := make([]qshell.ChgmEntryPath, len(entries))
			copy(toChgmEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChgm(toChgmEntries, bm)
			}
			entries = make([]qshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChgmEntries := make([]qshell.ChgmEntryPath, len(entries))
		copy(toChgmEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChgm(toChgmEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchChgm(entries []qshell.ChgmEntryPath, bm *qshell.BucketManager) {
	ret, err := bm.BatchChgm(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chgm error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Chgm '%s' => '%s' Failed, Code: %d, Error: %s", entry.Key, entry.MimeType, item.Code, item.Data.Error)
			} else {
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
			os.Exit(qshell.STATUS_HALT)
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

	bm := qshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyTypeMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyTypeMapFile)
		if err != nil {
			fmt.Printf("Open key file type map file error: %v\n", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)

	var key, line string
	var fileType int
	var items []string
	var entry qshell.ChtypeEntryPath

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
			entry = qshell.ChtypeEntryPath{
				EntryPath: qshell.EntryPath{
					Bucket: bucket,
					Key:    key,
				},
				FileType: fileType,
			}
			entries = append(entries, entry)
		}
	}

	var errEntries []qshell.ChtypeEntryPath

	var batches = len(entries)/BATCH_ALLOW_MAX + 1
	var completed int

	for i := 0; i < batches; i++ {
		var batch []qshell.ChtypeEntryPath

		if i == batches-1 {
			batch = entries
		} else {
			batch = entries[:BATCH_ALLOW_MAX]
			entries = entries[BATCH_ALLOW_MAX:]
		}

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			for _, entry := range batchChtype(batch, bm) {
				errEntries = append(errEntries, entry)
			}
			completed += 1
			fmt.Printf("\rComplete: %%%.1f", float64(completed)/float64(batches)*100)
		}
	}
	batchWaitGroup.Wait()

	fmt.Println()
	if len(errEntries) > 0 {
		fmt.Fprintf(os.Stderr, "Total %d entries failed: \n", len(errEntries))
		fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("-", 30))
		for _, entry := range errEntries {
			fmt.Fprintf(os.Stderr, "%s\t%d\n", entry.Key, entry.FileType)
		}
	}
}

func batchChtype(entries []qshell.ChtypeEntryPath, bm *qshell.BucketManager) (errEntries []qshell.ChtypeEntryPath) {
	ret, err := bm.BatchChtype(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch chtype error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fmt.Fprintf(os.Stderr, "Chtype '%s' => '%d' Failed, Code: %d, Error: %s\n", entry.Key, entry.FileType, item.Code, item.Data.Error)
				errEntries = append(errEntries, entry)
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
			os.Exit(qshell.STATUS_HALT)
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

	bm := qshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if keyExpireMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyExpireMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open key expire map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			key := items[0]
			days, _ := strconv.Atoi(items[1])
			if key != "" {
				entry := qshell.DeleteAfterDaysEntryPath{
					EntryPath: qshell.EntryPath{
						Bucket: bucket,
						Key:    key,
					},
					DeleteAfterDays: days,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toExpireEntries := make([]qshell.DeleteAfterDaysEntryPath, len(entries))
			copy(toExpireEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDeleteAfterDays(toExpireEntries, bm)
			}
			entries = make([]qshell.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toExpireEntries := make([]qshell.DeleteAfterDaysEntryPath, len(entries))
		copy(toExpireEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDeleteAfterDays(toExpireEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchDeleteAfterDays(entries []qshell.DeleteAfterDaysEntryPath, bm *qshell.BucketManager) {
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
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]
	var oldNewKeyMapFile string

	if inputFile == "" {
		oldNewKeyMapFile = "stdin"
	} else {
		oldNewKeyMapFile = params[1]
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

	bm := qshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if oldNewKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open old new key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			oldKey := items[0]
			newKey := items[1]
			if oldKey != "" && newKey != "" {
				entry := qshell.RenameEntryPath{
					SrcEntry: qshell.EntryPath{
						Bucket: bucket,
						Key:    oldKey,
					},
					DstEntry: qshell.EntryPath{
						Bucket: bucket,
						Key:    newKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toRenameEntries := make([]qshell.RenameEntryPath, len(entries))
			copy(toRenameEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchRename(toRenameEntries, bm)
			}
			entries = make([]qshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toRenameEntries := make([]qshell.RenameEntryPath, len(entries))
		copy(toRenameEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchRename(toRenameEntries, bm)
		}
	}
	batchWaitGroup.Wait()
}

func batchRename(entries []qshell.RenameEntryPath, bm *qshell.BucketManager) {
	ret, err := bm.BatchRename(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch rename error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Rename '%s' => '%s' Failed, Code: %d, Error: %s", entry.SrcEntry.Key, entry.DstEntry.Key, item.Code, item.Data.Error)
			} else {
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
			os.Exit(qshell.STATUS_HALT)
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

	bm := qshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	entries := make([]qshell.MoveEntryPath, 0, BATCH_ALLOW_MAX)
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
				entry := qshell.MoveEntryPath{
					SrcEntry: qshell.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: qshell.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toMoveEntries := make([]qshell.MoveEntryPath, len(entries))
			copy(toMoveEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchMove(toMoveEntries, bm)
			}
			entries = make([]qshell.MoveEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toMoveEntries := make([]qshell.MoveEntryPath, len(entries))
		copy(toMoveEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchMove(toMoveEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchMove(entries []qshell.MoveEntryPath, bm *qshell.BucketManager) {
	ret, err := bm.BatchMove(entries)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch move error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
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
			os.Exit(qshell.STATUS_HALT)
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

	bm := qshell.GetBucketManager()
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open src dest key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]qshell.CopyEntryPath, 0, BATCH_ALLOW_MAX)
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
				entry := qshell.CopyEntryPath{
					SrcEntry: qshell.EntryPath{
						Bucket: srcBucket,
						Key:    srcKey,
					},
					DstEntry: qshell.EntryPath{
						Bucket: destBucket,
						Key:    destKey,
					},
					Force: overwriteFlag,
				}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toCopyEntries := make([]qshell.CopyEntryPath, len(entries))
			copy(toCopyEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchCopy(toCopyEntries, bm)
			}
			entries = make([]qshell.CopyEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toCopyEntries := make([]qshell.CopyEntryPath, len(entries))
		copy(toCopyEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchCopy(toCopyEntries, bm)
		}
	}

	batchWaitGroup.Wait()
}

func batchCopy(entries []qshell.CopyEntryPath, bm *qshell.BucketManager) {
	ret, err := bm.BatchCopy(entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Batch copy error: %v\n", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key,
					item.Code, item.Data.Error)
			} else {
				logs.Debug("Copy '%s:%s' => '%s:%s' success",
					entry.SrcEntry.Bucket, entry.SrcEntry.Key,
					entry.DstEntry.Bucket, entry.DstEntry.Key)
			}
		}
	}
}

func BatchSign(cmd *cobra.Command, params []string) {
	urlListFile := params[0]
	var deadline int64
	if len(params) == 2 {
		if val, err := strconv.ParseInt(params[1], 10, 64); err != nil {
			fmt.Fprintln(os.Stderr, "Invalid <Deadline>")
			os.Exit(qshell.STATUS_HALT)
		} else {
			deadline = val
		}
	} else {
		deadline = time.Now().Add(time.Second * 3600 * 24 * 365).Unix()
	}

	bm := qshell.GetBucketManager()
	fp, openErr := os.Open(urlListFile)
	if openErr != nil {
		fmt.Fprintln(os.Stderr, "Open url list file error,", openErr)
		os.Exit(qshell.STATUS_HALT)
	}
	defer fp.Close()

	bReader := bufio.NewScanner(fp)
	for bReader.Scan() {
		urlToSign := strings.TrimSpace(bReader.Text())
		if urlToSign == "" {
			continue
		}
		signedUrl, _ := bm.PrivateUrl(urlToSign, deadline)
		fmt.Println(signedUrl)
	}
}
