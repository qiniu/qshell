package cli

import (
	"bufio"
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/rs"
	"qiniu/log"
	"qiniu/rpc"
	"qshell"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	BATCH_ALLOW_MAX = 1000
)

func printStat(bucket string, key string, entry rs.Entry) {
	statInfo := fmt.Sprintf("%-20s%s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "Fsize:", entry.Fsize, FormatFsize(entry.Fsize))

	putTime := time.Unix(0, entry.PutTime*100)
	statInfo += fmt.Sprintf("%-20s%d -> %s\r\n", "PutTime:", entry.PutTime, putTime.String())
	statInfo += fmt.Sprintf("%-20s%s\r\n", "MimeType:", entry.MimeType)
	fmt.Println(statInfo)
}

func DirCache(cmd string, params ...string) {
	if len(params) == 2 {
		cacheRootPath := params[0]
		cacheResultFile := params[1]
		dircacheS.Cache(cacheRootPath, cacheResultFile)
	} else {
		CmdHelp(cmd)
	}
}

func ListBucket(cmd string, params ...string) {
	if len(params) == 2 || len(params) == 3 {
		bucket := params[0]
		prefix := ""
		listResultFile := ""
		if len(params) == 2 {
			listResultFile = params[1]
		} else if len(params) == 3 {
			prefix = params[1]
			listResultFile = params[2]
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		if accountS.AccessKey != "" && accountS.SecretKey != "" {
			listbucketS.Account = accountS
			listbucketS.List(bucket, prefix, ListMarker, listResultFile)
		} else {
			fmt.Println("No AccessKey and SecretKey set error!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func Stat(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		entry, err := client.Stat(nil, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Stat error,", v.Code, v.Err)
			} else {
				fmt.Println("Stat error,", err.Error())
			}
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		err := client.Delete(nil, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Delete error,", v.Code, v.Err)
			} else {
				fmt.Println("Delete error,", err.Error())
			}
		} else {
			fmt.Println("Delete done!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func Move(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
		srcBucket := params[0]
		srcKey := params[1]
		destBucket := params[2]
		destKey := srcKey
		if len(params) == 4 {
			destKey = params[3]
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		err := client.Move(nil, srcBucket, srcKey, destBucket, destKey)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Move error,", v.Code, v.Err)
			} else {
				fmt.Println("Move error,", err.Error())
			}
		} else {
			fmt.Println("Move done!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func Copy(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
		srcBucket := params[0]
		srcKey := params[1]
		destBucket := params[2]
		destKey := srcKey
		if len(params) == 4 {
			destKey = params[3]
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		err := client.Copy(nil, srcBucket, srcKey, destBucket, destKey)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Copy error,", v.Code, v.Err)
			} else {
				fmt.Println("Copy error,", err.Error())
			}
		} else {
			fmt.Println("Copy done!")
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		err := client.ChangeMime(nil, bucket, key, newMimeType)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Change mimetype error,", v.Code, v.Err)
			} else {
				fmt.Println("Change mimetype error,", err.Error())
			}
		} else {
			fmt.Println("Change mimetype done!")
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		fetchResult, err := qshell.Fetch(&mac, remoteResUrl, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Fetch error,", v.Code, v.Err)
			} else {
				fmt.Println("Fetch error,", err.Error())
			}
		} else {
			fmt.Println("Key:", fetchResult.Key)
			fmt.Println("Hash:", fetchResult.Hash)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Prefetch(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		err := qshell.Prefetch(&mac, bucket, key)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("Prefetch error,", v.Code, v.Err)
			} else {
				fmt.Println("Prefetch error,", err.Error())
			}
		} else {
			fmt.Println("Prefetch done!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func BatchStat(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		keyListFile := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
			if item.Data.Error != "" {
				fmt.Println(entry.Key + "\t" + item.Data.Error)
			} else {
				fmt.Println(fmt.Sprintf("%s\t%d\t%s\t%s\t%d", entry.Key,
					item.Data.Fsize, item.Data.Hash, item.Data.MimeType, item.Data.PutTime))
			}
		}
	} else {
		if err != nil {
			fmt.Println("Batch stat error", err)
		}
	}
}

func BatchDelete(cmd string, params ...string) {
	if len(params) == 2 {
		if !ForceMode {
			//confirm
			rcode := CreateRandString(6)
			if rcode == "" {
				fmt.Println("Create confirm code failed")
				return
			}

			rcode2 := ""
			if runtime.GOOS == "windows" {
				fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
			} else {
				fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
			}
			fmt.Scanln(&rcode2)

			if rcode != rcode2 {
				fmt.Println("Task quit!")
				return
			}
		}

		bucket := params[0]
		keyListFile := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(keyListFile)
		if err != nil {
			fmt.Println("Open key list file error", err)
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
				batchDelete(client, entries)
				//reset slice
				entries = make([]rs.EntryPath, 0)
			}
		}
		//delete the last batch
		if len(entries) > 0 {
			batchDelete(client, entries)
		}
		fmt.Println("Batch delete done!")
	} else {
		CmdHelp(cmd)
	}
}

func batchDelete(client rs.Client, entries []rs.EntryPath) {
	ret, _ := qshell.BatchDelete(client, entries)

	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]

			if item.Data.Error != "" {
				log.Errorf("Delete '%s' => '%s' failed, Code: %d, Error: %s", entry.Bucket, entry.Key, item.Code, item.Data.Error)
			} else {
				log.Debugf("Delete '%s' => '%s' success", entry.Bucket, entry.Key)
			}
		}
	}
}

func BatchChgm(cmd string, params ...string) {
	if len(params) == 2 {
		if !ForceMode {
			//confirm
			rcode := CreateRandString(6)
			if rcode == "" {
				fmt.Println("Create confirm code failed")
				return
			}

			rcode2 := ""
			if runtime.GOOS == "windows" {
				fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
			} else {
				fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
			}
			fmt.Scanln(&rcode2)

			if rcode != rcode2 {
				fmt.Println("Task quit!")
				return
			}
		}

		bucket := params[0]
		keyMimeMapFile := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(keyMimeMapFile)
		if err != nil {
			fmt.Println("Open key mime map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.ChgmEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
				batchChgm(client, entries)
				entries = make([]qshell.ChgmEntryPath, 0)
			}
		}
		if len(entries) > 0 {
			batchChgm(client, entries)
		}
		fmt.Println("All Chgmed!")
	} else {
		CmdHelp(cmd)
	}
}

func batchChgm(client rs.Client, entries []qshell.ChgmEntryPath) {
	ret, err := qshell.BatchChgm(client, entries)
	if err != nil {
		fmt.Println("Batch chgm error")
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Errorf("Chgm '%s' => '%s' Failed, Code: %d, Error: %s", entry.Key, entry.MimeType, item.Code, item.Data.Error)
			} else {
				log.Debug(fmt.Sprintf("Chgm '%s' => '%s' success", entry.Key, entry.MimeType))
			}
		}
	}
}

func BatchRename(cmd string, params ...string) {
	if len(params) == 2 {
		if !ForceMode {
			//confirm
			rcode := CreateRandString(6)
			if rcode == "" {
				fmt.Println("Create confirm code failed")
				return
			}

			rcode2 := ""
			if runtime.GOOS == "windows" {
				fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
			} else {
				fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
			}
			fmt.Scanln(&rcode2)

			if rcode != rcode2 {
				fmt.Println("Task quit!")
				return
			}
		}

		bucket := params[0]
		oldNewKeyMapFile := params[1]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(oldNewKeyMapFile)
		if err != nil {
			fmt.Println("Open old new key map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.RenameEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
				batchRename(client, entries)
				entries = make([]qshell.RenameEntryPath, 0)
			}
		}
		if len(entries) > 0 {
			batchRename(client, entries)
		}
		fmt.Println("All Renamed!")
	} else {
		CmdHelp(cmd)
	}
}

func batchRename(client rs.Client, entries []qshell.RenameEntryPath) {
	ret, err := qshell.BatchRename(client, entries)
	if err != nil {
		fmt.Println("Batch rename error")
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Errorf("Rename '%s' => '%s' Failed, Code: %d, Error: %s", entry.OldKey, entry.NewKey, item.Code, item.Data.Error)
			} else {
				log.Debug(fmt.Sprintf("Rename '%s' => '%s' success", entry.OldKey, entry.NewKey))
			}
		}
	}
}

func BatchMove(cmd string, params ...string) {
	if len(params) == 3 {
		if !ForceMode {
			//confirm
			rcode := CreateRandString(6)
			if rcode == "" {
				fmt.Println("Create confirm code failed")
				return
			}

			rcode2 := ""
			if runtime.GOOS == "windows" {
				fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
			} else {
				fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
			}
			fmt.Scanln(&rcode2)

			if rcode != rcode2 {
				fmt.Println("Task quit!")
				return
			}
		}

		srcBucket := params[0]
		destBucket := params[1]
		srcDestKeyMapFile := params[2]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.MoveEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
				batchMove(client, entries)
				entries = make([]qshell.MoveEntryPath, 0)
			}
		}
		if len(entries) > 0 {
			batchMove(client, entries)
		}
		fmt.Println("All Moved!")
	} else {
		CmdHelp(cmd)
	}
}

func batchMove(client rs.Client, entries []qshell.MoveEntryPath) {
	ret, err := qshell.BatchMove(client, entries)
	if err != nil {
		fmt.Println("Batch move error")
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Errorf("Move '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code, item.Data.Error)
			} else {
				log.Debug(fmt.Sprintf("Move '%s:%s' => '%s:%s' success",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey))
			}
		}
	}
}

func BatchCopy(cmd string, params ...string) {
	if len(params) == 3 {
		if !ForceMode {
			//confirm
			rcode := CreateRandString(6)
			if rcode == "" {
				fmt.Println("Create confirm code failed")
				return
			}

			rcode2 := ""
			if runtime.GOOS == "windows" {
				fmt.Print(fmt.Sprintf("<DANGER> Input %s to confirm operation: ", rcode))
			} else {
				fmt.Print(fmt.Sprintf("\033[31m<DANGER>\033[0m Input \033[32m%s\033[0m to confirm operation: ", rcode))
			}
			fmt.Scanln(&rcode2)

			if rcode != rcode2 {
				fmt.Println("Task quit!")
				return
			}
		}

		srcBucket := params[0]
		destBucket := params[1]
		srcDestKeyMapFile := params[2]

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.NewMac(&mac)
		fp, err := os.Open(srcDestKeyMapFile)
		if err != nil {
			fmt.Println("Open src dest key map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.CopyEntryPath, 0, BATCH_ALLOW_MAX)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
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
				batchCopy(client, entries)
				entries = make([]qshell.CopyEntryPath, 0)
			}
		}
		if len(entries) > 0 {
			batchCopy(client, entries)
		}
		fmt.Println("All Copyed!")
	} else {
		CmdHelp(cmd)
	}
}

func batchCopy(client rs.Client, entries []qshell.CopyEntryPath) {
	ret, err := qshell.BatchCopy(client, entries)
	if err != nil {
		fmt.Println("Batch copy error")
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Errorf("Copy '%s:%s' => '%s:%s' Failed, Code: %d, Error: %s",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code, item.Data.Error)
			} else {
				log.Debug(fmt.Sprintf("Copy '%s:%s' => '%s:%s' success",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey))
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
				return
			} else {
				deadline = val
			}
		} else {
			deadline = time.Now().Add(time.Second * 3600).Unix()
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		url := qshell.PrivateUrl(&mac, publicUrl, deadline)
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
				return
			} else {
				deadline = val
			}
		} else {
			deadline = time.Now().Add(time.Second * 3600 * 24 * 365).Unix()
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}

		fp, openErr := os.Open(urlListFile)
		if openErr != nil {
			fmt.Println("Open url list file error,", openErr)
			return
		}
		defer fp.Close()

		bReader := bufio.NewScanner(fp)
		bReader.Split(bufio.ScanLines)
		for bReader.Scan() {
			urlToSign := strings.TrimSpace(bReader.Text())
			if urlToSign == "" {
				continue
			}
			signedUrl := qshell.PrivateUrl(&mac, urlToSign, deadline)
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		url, err := qshell.Saveas(&mac, publicUrl, saveBucket, saveKey)
		if err != nil {
			fmt.Println(err)
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		m3u8FileList, err := qshell.M3u8FileList(&mac, bucket, m3u8Key)
		if err != nil {
			fmt.Println(err)
			return
		}
		client := rs.NewMac(&mac)
		entryCnt := len(m3u8FileList)
		if entryCnt == 0 {
			fmt.Println("no m3u8 slices found")
			return
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
		fmt.Println("All deleted!")
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

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		err := qshell.M3u8ReplaceDomain(&mac, bucket, m3u8Key, newDomain)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok {
				fmt.Println("m3u8 replace domain error,", v.Err)
			} else {
				fmt.Println("m3u8 replace domain error,", err)
			}
			return
		} else {
			fmt.Println("Done!")
		}
	} else {
		CmdHelp(cmd)
	}
}
