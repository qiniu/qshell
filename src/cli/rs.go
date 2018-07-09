package cli

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"net/http"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
	"qshell"
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

func printStat(bucket string, key string, entry rs.Entry) {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", entry.Fsize, FormatFsize(entry.Fsize))

	putTime := time.Unix(0, entry.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", entry.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", entry.MimeType)
	if entry.FileType == 0 {
		statInfo += fmt.Sprintf("%-20s%d -> 标准存储\r\n", "FileType:", entry.FileType)
	} else {
		statInfo += fmt.Sprintf("%-20s%d -> 低频存储\r\n", "FileType:", entry.FileType)
	}
	fmt.Println(statInfo)
}

func DirCache(cmd string, params ...string) {
	if len(params) == 2 {
		cacheRootPath := params[0]
		cacheResultFile := params[1]
		_, retErr := qshell.DirCache(cacheRootPath, cacheResultFile)
		if retErr != nil {
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func ListBucket(cmd string, params ...string) {
	var listMarker string
	flagSet := flag.NewFlagSet("listbucket", flag.ExitOnError)
	flagSet.StringVar(&listMarker, "marker", "", "list marker")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 || len(cmdParams) == 3 {
		bucket := cmdParams[0]
		prefix := ""
		listResultFile := ""
		if len(cmdParams) == 2 {
			listResultFile = cmdParams[1]
		} else if len(cmdParams) == 3 {
			prefix = cmdParams[1]
			listResultFile = cmdParams[2]
		}

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}

		if !IsHostFileSpecified {
			//get zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Failed to get region info of bucket", bucket, gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set zone
			qshell.SetZone(bucketInfo.Region)
		}

		retErr := qshell.ListBucket(&mac, bucket, prefix, listMarker, listResultFile)
		if retErr != nil {
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func ListBucket2(cmd string, params ...string) {
	var listMarker string
	flagSet := flag.NewFlagSet("listbucket", flag.ExitOnError)
	flagSet.StringVar(&listMarker, "marker", "", "list marker")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 || len(cmdParams) == 3 {
		bucket := cmdParams[0]
		prefix := ""
		listResultFile := ""
		if len(cmdParams) == 2 {
			listResultFile = cmdParams[1]
		} else if len(cmdParams) == 3 {
			prefix = cmdParams[1]
			listResultFile = cmdParams[2]
		}

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}

		if !IsHostFileSpecified {
			//get zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Failed to get region info of bucket", bucket, gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set zone
			qshell.SetZone(bucketInfo.Region)
		}

		nextMarker, retErr := qshell.ListBucketV2(&mac, bucket, prefix, listMarker, listResultFile)
		if nextMarker != "" {
			fmt.Println("Next Marker:", nextMarker)
		}
		if retErr != nil {
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Stat(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]

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

		entry, err := client.Stat(nil, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Stat error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Stat error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		} else {
			printStat(bucket, key, entry)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Delete(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]

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

		err := client.Delete(nil, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Delete error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Delete error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Move(cmd string, params ...string) {
	var overwrite bool
	flagSet := flag.NewFlagSet("move", flag.ExitOnError)
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite mode")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 3 || len(cmdParams) == 4 {
		srcBucket := cmdParams[0]
		srcKey := cmdParams[1]
		destBucket := cmdParams[2]
		destKey := srcKey
		if len(cmdParams) == 4 {
			destKey = cmdParams[3]
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

		err := client.Move(nil, srcBucket, srcKey, destBucket, destKey, overwrite)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Move error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Move error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Copy(cmd string, params ...string) {
	var overwrite bool
	flagSet := flag.NewFlagSet("copy", flag.ExitOnError)
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite mode")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 3 || len(cmdParams) == 4 {
		srcBucket := cmdParams[0]
		srcKey := cmdParams[1]
		destBucket := cmdParams[2]
		destKey := srcKey
		if len(cmdParams) == 4 {
			destKey = cmdParams[3]
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

		err := client.Copy(nil, srcBucket, srcKey, destBucket, destKey, overwrite)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Copy error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Copy error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Chgm(cmd string, params ...string) {
	if len(params) == 3 {
		bucket := params[0]
		key := params[1]
		newMimeType := params[2]

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

		err := client.ChangeMime(nil, bucket, key, newMimeType)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Change mimetype error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Change mimetype error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Chtype(cmd string, params ...string) {
	if len(params) == 3 {
		bucket := params[0]
		key := params[1]
		fileTypeStr := params[2]
		fileType, cErr := strconv.Atoi(fileTypeStr)
		if cErr != nil {
			fmt.Println("Invalid file type")
			os.Exit(qshell.STATUS_HALT)
			return
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

		err := client.ChangeType(nil, bucket, key, fileType)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Change file type error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Change file type error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func DeleteAfterDays(cmd string, params ...string) {
	if len(params) == 3 {
		bucket := params[0]
		key := params[1]
		expireStr := params[2]
		expire, cErr := strconv.Atoi(expireStr)
		if cErr != nil {
			fmt.Println("Invalid deleteAfterDays")
			os.Exit(qshell.STATUS_HALT)
			return
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

		err := client.DeleteAfterDays(nil, bucket, key, expire)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Set file deleteAfterDays error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Set file deleteAfterDays error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Fetch(cmd string, params ...string) {
	if len(params) == 2 || len(params) == 3 {
		remoteResUrl := params[0]
		bucket := params[1]
		key := ""
		if len(params) == 3 {
			key = params[2]
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

		if !IsHostFileSpecified {
			//get bucket zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Get bucket region info error,", gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set up host
			qshell.SetZone(bucketInfo.Region)
		}

		fetchResult, err := qshell.Fetch(&mac, remoteResUrl, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Fetch error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Fetch error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		} else {
			fmt.Println("Key:", fetchResult.Key)
			fmt.Println("Hash:", fetchResult.Hash)
			fmt.Printf("Fsize: %d (%s)\n", fetchResult.Fsize, FormatFsize(fetchResult.Fsize))
			fmt.Println("Mime:", fetchResult.MimeType)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Prefetch(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}

		if !IsHostFileSpecified {
			//get bucket zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Get bucket region info error,", gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set up host
			qshell.SetZone(bucketInfo.Region)
		}

		err := qshell.Prefetch(&mac, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Printf("Prefetch error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
			} else {
				fmt.Println("Prefetch error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}

func BatchStat(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		keyListFile := params[1]

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

		fp, err := os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := scanner.Text()
			items := strings.Split(line, "\t")
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
	} else {
		CmdHelp(cmd)
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

func BatchDelete(cmd string, params ...string) {
	var force bool
	var worker int
	flagSet := flag.NewFlagSet("batchdelete", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 {
		if !force {
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

		bucket := cmdParams[0]
		keyListFile := cmdParams[1]

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

		fp, err := os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := scanner.Text()
			items := strings.Split(line, "\t")
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
	} else {
		CmdHelp(cmd)
	}
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

func BatchChgm(cmd string, params ...string) {
	var force bool
	var worker int
	flagSet := flag.NewFlagSet("batchchgm", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 {
		if !force {
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

		bucket := cmdParams[0]
		keyMimeMapFile := cmdParams[1]

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

		fp, err := os.Open(keyMimeMapFile)
		if err != nil {
			fmt.Println("Open key mime map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := scanner.Text()
			items := strings.Split(line, "\t")
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
	} else {
		CmdHelp(cmd)
	}
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

func BatchChtype(cmd string, params ...string) {
	var force bool
	var worker int
	flagSet := flag.NewFlagSet("batchchtype", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 {
		if !force {
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

		bucket := cmdParams[0]
		keyTypeMapFile := cmdParams[1]

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

		fp, err := os.Open(keyTypeMapFile)
		if err != nil {
			fmt.Println("Open key file type map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.ChtypeEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := scanner.Text()
			items := strings.Split(line, "\t")
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
	} else {
		CmdHelp(cmd)
	}
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

func BatchDeleteAfterDays(cmd string, params ...string) {
	var force bool
	var worker int
	flagSet := flag.NewFlagSet("batchepxire", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 {
		if !force {
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

		bucket := cmdParams[0]
		keyExpireMapFile := cmdParams[1]

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

		fp, err := os.Open(keyExpireMapFile)
		if err != nil {
			fmt.Println("Open key expire map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
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
	} else {
		CmdHelp(cmd)
	}
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

func BatchRename(cmd string, params ...string) {
	var force bool
	var overwrite bool
	var worker int
	flagSet := flag.NewFlagSet("batchrename", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 2 {
		if !force {
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

		bucket := cmdParams[0]
		oldNewKeyMapFile := cmdParams[1]

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

		fp, err := os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Println("Open old new key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
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
					batchRename(client, toRenameEntries, overwrite)
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
				batchRename(client, toRenameEntries, overwrite)
			}
		}
		batchWaitGroup.Wait()
	} else {
		CmdHelp(cmd)
	}
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

func BatchMove(cmd string, params ...string) {
	var force bool
	var overwrite bool
	var worker int
	flagSet := flag.NewFlagSet("batchmove", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 3 {
		if !force {
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

		srcBucket := cmdParams[0]
		destBucket := cmdParams[1]
		srcDestKeyMapFile := cmdParams[2]

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

		fp, err := os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
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
					batchMove(client, toMoveEntries, overwrite)
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
				batchMove(client, toMoveEntries, overwrite)
			}
		}

		batchWaitGroup.Wait()
	} else {
		CmdHelp(cmd)
	}
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

func BatchCopy(cmd string, params ...string) {
	var force bool
	var overwrite bool
	var worker int

	flagSet := flag.NewFlagSet("batchcopy", flag.ExitOnError)
	flagSet.BoolVar(&force, "force", false, "force mode")
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite mode")
	flagSet.IntVar(&worker, "worker", 1, "worker count")
	flagSet.Parse(params)

	cmdParams := flagSet.Args()
	if len(cmdParams) == 3 {
		if !force {
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

		srcBucket := cmdParams[0]
		destBucket := cmdParams[1]
		srcDestKeyMapFile := cmdParams[2]

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

		fp, err := os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
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
					batchCopy(client, toCopyEntries, overwrite)
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
				batchCopy(client, toCopyEntries, overwrite)
			}
		}

		batchWaitGroup.Wait()
	} else {
		CmdHelp(cmd)
	}
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

func PrivateUrl(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		publicUrl := params[0]
		var deadline int64
		if len(params) == 2 {
			if val, err := strconv.ParseInt(params[1], 10, 64); err != nil {
				fmt.Println("Invalid <Deadline>")
				os.Exit(qshell.STATUS_HALT)
			} else {
				deadline = val
			}
		} else {
			deadline = time.Now().Add(time.Second * 3600).Unix()
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
		url, _ := qshell.PrivateUrl(&mac, publicUrl, deadline)
		fmt.Println(url)
	} else {
		CmdHelp(cmd)
	}
}

func BatchSign(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
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
	} else {
		CmdHelp(cmd)
	}
}

func Saveas(cmd string, params ...string) {
	if len(params) == 3 {
		publicUrl := params[0]
		saveBucket := params[1]
		saveKey := params[2]

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}
		url, err := qshell.Saveas(&mac, publicUrl, saveBucket, saveKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(qshell.STATUS_ERROR)
		} else {
			fmt.Println(url)
		}
	} else {
		CmdHelp(cmd)
	}
}

func M3u8Delete(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		m3u8Key := params[1]

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		mac := digest.Mac{
			account.AccessKey,
			[]byte(account.SecretKey),
		}

		//get bucket zone info
		bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
		if gErr != nil {
			fmt.Println("Get bucket region info error,", gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		//set up host
		qshell.SetZone(bucketInfo.Region)

		m3u8FileList, err := qshell.M3u8FileList(&mac, bucket, m3u8Key)
		if err != nil {
			fmt.Println(err)
			os.Exit(qshell.STATUS_ERROR)
		}
		client := rs.NewMacEx(&mac, &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   time.Duration(60) * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: time.Second * 60 * 10,
		}, "")
		entryCnt := len(m3u8FileList)
		if entryCnt == 0 {
			fmt.Println("no m3u8 slices found")
			os.Exit(qshell.STATUS_ERROR)
		}
		if entryCnt <= BATCH_ALLOW_MAX {
			batchDelete(client, m3u8FileList)
		} else {
			batchCnt := entryCnt / BATCH_ALLOW_MAX
			for i := 0; i < batchCnt; i++ {
				end := (i + 1) * BATCH_ALLOW_MAX
				if end > entryCnt {
					end = entryCnt
				}
				entriesToDelete := m3u8FileList[i*BATCH_ALLOW_MAX : end]
				batchDelete(client, entriesToDelete)
			}
		}
	} else {
		CmdHelp(cmd)
	}
}

func M3u8Replace(cmd string, params ...string) {
	if len(params) == 2 || len(params) == 3 {
		bucket := params[0]
		m3u8Key := params[1]
		var newDomain string
		if len(params) == 3 {
			newDomain = strings.TrimRight(params[2], "/")
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

		//get bucket zone info
		bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
		if gErr != nil {
			fmt.Println("Get bucket region info error,", gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		//set up host
		qshell.SetZone(bucketInfo.Region)

		err := qshell.M3u8ReplaceDomain(&mac, bucket, m3u8Key, newDomain)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("m3u8 replace domain error,", v.Err)
			} else {
				fmt.Println("m3u8 replace domain error,", err)
			}
			os.Exit(qshell.STATUS_ERROR)
		}
	} else {
		CmdHelp(cmd)
	}
}
