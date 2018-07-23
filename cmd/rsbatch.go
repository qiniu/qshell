package cmd

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"github.com/tonycai653/iqshell/qshell"
	"io"
	"net"
	"net/http"
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
)

var (
	batchStatCmd = &cobra.Command{
		Use:   "batchstat <Bucket> [<KeyListFile>]",
		Short: "Batch stat files in bucket",
		Long:  "Batch stat files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchStat,
	}
	batchDeleteCmd = &cobra.Command{
		Use:   "batchdelete <Bucket> [<KeyListFile>]",
		Short: "Batch delete files in bucket",
		Long:  "Batch delete files in bucket, read file list from stdin if KeyListFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchDelete,
	}
	batchChgmCmd = &cobra.Command{
		Use:   "batchchgm <Bucket> [<KeyMimeMapFile>]",
		Short: "Batch change the mime type of files in bucket",
		Long:  "Batch change the mime type of files in bucket, read from stdin if KeyMimeMapFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchChgm,
	}
	batchChtypeCmd = &cobra.Command{
		Use:   "batchchtype <Bucket> [<KeyFileTypeMapFile>]",
		Short: "Batch change the file type of files in bucket",
		Long:  "Batch change the file (storage) type of files in bucket, read from stdin if KeyFileTypeMapFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchChtype,
	}
	batchDelAfterCmd = &cobra.Command{
		Use:   "batchexpire <Bucket> [<KeyDeleteAfterDaysMapFile>]",
		Short: "Batch set the deleteAfterDays of the files in bucket",
		Long:  "Batch set the deleteAfterDays of the files in bucket, read from stdin if KeyDeleteAfterDaysMapFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchDeleteAfterDays,
	}
	batchRenameCmd = &cobra.Command{
		Use:   "batchrename <Bucket> [<OldNewKeyMapFile>]",
		Short: "Batch rename files in the bucket",
		Long:  "Batch rename files in the bucket, read from stdin if OldNewKeyMapFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   BatchRename,
	}
	batchMoveCmd = &cobra.Command{
		Use:   "batchmove <SrcBucket> <DestBucket> [<SrcDestKeyMapFile>]",
		Short: "Batch move files from bucket to bucket",
		Long:  "Batch move files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.RangeArgs(2, 3),
		Run:   BatchMove,
	}
	batchCopyCmd = &cobra.Command{
		Use:   "batchcopy <SrcBucket> <DestBucket> [<SrcDestKeyMapFile>]",
		Short: "Batch copy files from bucket to bucket",
		Long:  "Batch copy files from bucket to bucket, read from stdin if SrcDestKeyMapFile not specified",
		Args:  cobra.RangeArgs(2, 3),
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
	batchDeleteCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDeleteCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

	batchChgmCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChgmCmd.Flags().IntVarP(&worker, "worker", "c", 1, "woker count")

	batchChtypeCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchChtypeCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

	batchDelAfterCmd.Flags().BoolVarP(&forceFlag, "force", "y", false, "force mode")
	batchDelAfterCmd.Flags().IntVarP(&worker, "worker", "c", 1, "worker count")

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

	if len(params) == 2 {
		keyListFile = params[1]
	} else {
		keyListFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}
	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

	var fp io.ReadCloser
	var err error

	if keyListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := rs.EntryPath{
					bucket, key,
				}
				entries = append(entries, entry)
			}
		}
		//check 1000 limit
		if len(entries) == BATCH_ALLOW_MAX {
			batchStat(client, entries)
			//reset slice
			entries = make([]rs.EntryPath, 0)
		}
	}
	//stat the last batch
	if len(entries) > 0 {
		batchStat(client, entries)
	}
}

func batchStat(client rs.Client, entries []rs.EntryPath) {
	ret, err := qshell.BatchStat(client, entries)
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				fmt.Println(entry.Key + "\t" + item.Data.Error)
			} else {
				fmt.Println(fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%d", entry.Key,
					item.Data.Fsize, item.Data.Hash, item.Data.MimeType, item.Data.PutTime, item.Data.FileType))
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch stat error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch stat error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyListFile string

	if len(params) == 2 {
		keyListFile = params[1]
	} else {
		keyListFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

	var fp io.ReadCloser
	var err error

	if keyListFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) > 0 {
			key := items[0]
			if key != "" {
				entry := rs.EntryPath{
					bucket, key,
				}
				entries = append(entries, entry)
			}
		}
		//check limit
		if len(entries) == BATCH_ALLOW_MAX {
			toDeleteEntries := make([]rs.EntryPath, len(entries))
			copy(toDeleteEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDelete(client, toDeleteEntries)
			}
			entries = make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	//delete the last batch
	if len(entries) > 0 {
		toDeleteEntries := make([]rs.EntryPath, len(entries))
		copy(toDeleteEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchDelete(client, toDeleteEntries)
		}
	}

	batchWaitGroup.Wait()
}

func batchDelete(client rs.Client, entries []rs.EntryPath) {
	ret, err := qshell.BatchDelete(client, entries)

	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]

			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Delete '%s' => '%s' failed, Code: %d, Error: %s", entry.Bucket, entry.Key, item.Code, item.Data.Error)
			} else {
				logs.Debug("Delete '%s' => '%s' success", entry.Bucket, entry.Key)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch delete error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch delete error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyMimeMapFile string
	if len(params) == 2 {
		keyMimeMapFile = params[1]
	} else {
		keyMimeMapFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

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
				entry := qshell.ChgmEntryPath{bucket, key, mimeType}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChgmEntries := make([]qshell.ChgmEntryPath, len(entries))
			copy(toChgmEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChgm(client, toChgmEntries)
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
			batchChgm(client, toChgmEntries)
		}
	}

	batchWaitGroup.Wait()
}

func batchChgm(client rs.Client, entries []qshell.ChgmEntryPath) {
	ret, err := qshell.BatchChgm(client, entries)
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Chgm '%s' => '%s' Failed, Code: %d, Error: %s", entry.Key, entry.MimeType, item.Code, item.Data.Error)
			} else {
				logs.Debug("Chgm '%s' => '%s' success", entry.Key, entry.MimeType)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch chgm error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch chgm error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]

	var keyTypeMapFile string
	if len(params) == 2 {
		keyTypeMapFile = params[1]
	} else {
		keyTypeMapFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

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
	scanner.Split(bufio.ScanLines)
	entries := make([]qshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Fields(line)
		if len(items) == 2 {
			key := items[0]
			fileType, _ := strconv.Atoi(items[1])
			if key != "" {
				entry := qshell.ChtypeEntryPath{bucket, key, fileType}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toChtypeEntries := make([]qshell.ChtypeEntryPath, len(entries))
			copy(toChtypeEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChtype(client, toChtypeEntries)
			}
			entries = make([]qshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(entries) > 0 {
		toChtypeEntries := make([]qshell.ChtypeEntryPath, len(entries))
		copy(toChtypeEntries, entries)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChtype(client, toChtypeEntries)
		}
	}

	batchWaitGroup.Wait()
}

func batchChtype(client rs.Client, entries []qshell.ChtypeEntryPath) {
	ret, err := qshell.BatchChtype(client, entries)
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Chtype '%s' => '%d' Failed, Code: %d, Error: %s", entry.Key, entry.FileType, item.Code, item.Data.Error)
			} else {
				logs.Debug("Chtype '%s' => '%d' success", entry.Key, entry.FileType)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch chtype error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch chtype error,", err)
			}
		}
	}
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]
	var keyExpireMapFile string

	if len(params) == 2 {
		keyExpireMapFile = params[1]
	} else {
		keyExpireMapFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

	var fp io.ReadCloser
	var err error

	if keyExpireMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(keyExpireMapFile)
		if err != nil {
			fmt.Println("Open key expire map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]qshell.DeleteAfterDaysEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			key := items[0]
			days, _ := strconv.Atoi(items[1])
			if key != "" {
				entry := qshell.DeleteAfterDaysEntryPath{bucket, key, days}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toExpireEntries := make([]qshell.DeleteAfterDaysEntryPath, len(entries))
			copy(toExpireEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchDeleteAfterDays(client, toExpireEntries)
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
			batchDeleteAfterDays(client, toExpireEntries)
		}
	}

	batchWaitGroup.Wait()
}

func batchDeleteAfterDays(client rs.Client, entries []qshell.DeleteAfterDaysEntryPath) {
	ret, err := qshell.BatchDeleteAfterDays(client, entries)
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Expire '%s' => '%d' Failed, Code: %d, Error: %s", entry.Key, entry.DeleteAfterDays, item.Code, item.Data.Error)
			} else {
				logs.Debug("Expire '%s' => '%d' success", entry.Key, entry.DeleteAfterDays)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch expire error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch expire error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	bucket := params[0]
	var oldNewKeyMapFile string

	if len(params) == 2 {
		oldNewKeyMapFile = "stdin"
	} else {
		oldNewKeyMapFile = params[1]
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

	var fp io.ReadCloser
	var err error

	if oldNewKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Println("Open old new key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	entries := make([]qshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			oldKey := items[0]
			newKey := items[1]
			if oldKey != "" && newKey != "" {
				entry := qshell.RenameEntryPath{bucket, oldKey, newKey}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toRenameEntries := make([]qshell.RenameEntryPath, len(entries))
			copy(toRenameEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchRename(client, toRenameEntries, overwriteFlag)
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
			batchRename(client, toRenameEntries, overwriteFlag)
		}
	}
	batchWaitGroup.Wait()
}

func batchRename(client rs.Client, entries []qshell.RenameEntryPath, overwrite bool) {
	ret, err := qshell.BatchRename(client, entries, overwrite)

	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Rename '%s' => '%s' Failed, Code: %d, Error: %s", entry.OldKey, entry.NewKey, item.Code, item.Data.Error)
			} else {
				logs.Debug("Rename '%s' => '%s' success", entry.OldKey, entry.NewKey)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch rename error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch rename error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	srcBucket := params[0]
	destBucket := params[1]
	var srcDestKeyMapFile string

	if len(params) == 3 {
		srcDestKeyMapFile = params[2]
	} else {
		srcDestKeyMapFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")

	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
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
				entry := qshell.MoveEntryPath{srcBucket, destBucket, srcKey, destKey}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toMoveEntries := make([]qshell.MoveEntryPath, len(entries))
			copy(toMoveEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchMove(client, toMoveEntries, overwriteFlag)
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
			batchMove(client, toMoveEntries, overwriteFlag)
		}
	}

	batchWaitGroup.Wait()
}

func batchMove(client rs.Client, entries []qshell.MoveEntryPath, overwrite bool) {
	ret, err := qshell.BatchMove(client, entries, overwrite)

	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code, item.Data.Error)
			} else {
				logs.Debug("Move '%s:%s' => '%s:%s' success",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch move error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch move error,", err)
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
			fmt.Println("Task quit!")
			os.Exit(qshell.STATUS_HALT)
		}
	}

	srcBucket := params[0]
	destBucket := params[1]

	var srcDestKeyMapFile string

	if len(params) == 3 {
		srcDestKeyMapFile = params[2]
	} else {
		srcDestKeyMapFile = "stdin"
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
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

	client := rs.NewMacEx(&mac, &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(60) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: time.Second * 60 * 10,
	}, "")
	var fp io.ReadCloser
	var err error

	if srcDestKeyMapFile == "stdin" {
		fp = os.Stdin
	} else {
		fp, err = os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
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
				entry := qshell.CopyEntryPath{srcBucket, destBucket, srcKey, destKey}
				entries = append(entries, entry)
			}
		}
		if len(entries) == BATCH_ALLOW_MAX {
			toCopyEntries := make([]qshell.CopyEntryPath, len(entries))
			copy(toCopyEntries, entries)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchCopy(client, toCopyEntries, overwriteFlag)
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
			batchCopy(client, toCopyEntries, overwriteFlag)
		}
	}

	batchWaitGroup.Wait()
}

func batchCopy(client rs.Client, entries []qshell.CopyEntryPath, overwrite bool) {
	ret, err := qshell.BatchCopy(client, entries, overwrite)

	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Code != 200 || item.Data.Error != "" {
				logs.Error("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code, item.Data.Error)
			} else {
				logs.Debug("Copy '%s:%s' => '%s:%s' success",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey)
			}
		}
	} else {
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Batch copy error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Batch copy error,", err)
			}
		}
	}
}

func BatchSign(cmd *cobra.Command, params []string) {
	urlListFile := params[0]
	var deadline int64
	if len(params) == 2 {
		if val, err := strconv.ParseInt(params[1], 10, 64); err != nil {
			fmt.Println("Invalid <Deadline>")
			os.Exit(qshell.STATUS_HALT)
		} else {
			deadline = val
		}
	} else {
		deadline = time.Now().Add(time.Second * 3600 * 24 * 365).Unix()
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{
		account.AccessKey,
		[]byte(account.SecretKey),
	}

	fp, openErr := os.Open(urlListFile)
	if openErr != nil {
		fmt.Println("Open url list file error,", openErr)
		os.Exit(qshell.STATUS_HALT)
	}
	defer fp.Close()

	bReader := bufio.NewScanner(fp)
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		urlToSign := strings.TrimSpace(bReader.Text())
		if urlToSign == "" {
			continue
		}
		signedUrl, _ := qshell.PrivateUrl(&mac, urlToSign, deadline)
		fmt.Println(signedUrl)
	}
}
