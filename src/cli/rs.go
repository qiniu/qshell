package cli

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"qshell"
	"strconv"
	"strings"
	"time"
)

const (
	BATCH_ALLOW_MAX = 1000
)

func printStat(bucket string, key string, entry rs.Entry) {
	statInfo := fmt.Sprintf("%-20s%-20s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "Fsize:", entry.Fsize)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "PutTime:", entry.PutTime)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "MimeType:", entry.MimeType)
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
		accountS.Get()
		if accountS.AccessKey != "" && accountS.SecretKey != "" {
			listbucketS.Account = accountS
			listbucketS.List(bucket, prefix, listResultFile)
		} else {
			log.Error("No AccessKey and SecretKey set error!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func Stat(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		key := params[1]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		entry, err := client.Stat(nil, bucket, key)
		if err != nil {
			log.Error("Stat error,", err)
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
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		err := client.Delete(nil, bucket, key)
		if err != nil {
			log.Error("Delete error,", err)
		} else {
			fmt.Println("Done!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func Move(cmd string, params ...string) {
	if len(params) == 4 {
		srcBucket := params[0]
		srcKey := params[1]
		destBucket := params[2]
		destKey := params[3]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		err := client.Move(nil, srcBucket, srcKey, destBucket, destKey)
		if err != nil {
			log.Error("Move error,", err)
		} else {
			fmt.Println("Done!")
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
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		err := client.Copy(nil, srcBucket, srcKey, destBucket, destKey)
		if err != nil {
			log.Error("Copy error,", err)
		} else {
			fmt.Println("Done!")
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
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		err := client.ChangeMime(nil, bucket, key, newMimeType)
		if err != nil {
			log.Error("Change mimeType error,", err)
		} else {
			fmt.Println("Done!")
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
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		fetchResult, err := qshell.Fetch(&mac, remoteResUrl, bucket, key)
		if err != nil {
			log.Error("Fetch error,", err)
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
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		err := qshell.Prefetch(&mac, bucket, key)
		if err != nil {
			log.Error("Prefetch error,", err)
		} else {
			fmt.Println("Done!")
		}
	} else {
		CmdHelp(cmd)
	}
}

func BatchDelete(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		keyListFile := params[1]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		fp, err := os.Open(keyListFile)
		if err != nil {
			log.Error("Open key list file error", err)
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			items := strings.Split(line, "\t")
			if len(items) > 0 {
				key := strings.TrimSpace(items[0])
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
		fmt.Println("All deleted!")
	} else {
		CmdHelp(cmd)
	}
}

func batchDelete(client rs.Client, entries []rs.EntryPath) {
	ret, err := qshell.BatchDelete(client, entries)
	if err != nil {
		log.Error("Batch delete error", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Error(fmt.Sprintf("Delete '%s' => '%s' Failed, Code: %d", entry.Bucket, entry.Key, item.Code))
			} else {
				log.Debug(fmt.Sprintf("Delete '%s' => '%s' Success, Code: %d", entry.Bucket, entry.Key, item.Code))
			}
		}
	}
}

func BatchChgm(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		keyMimeMapFile := params[1]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		fp, err := os.Open(keyMimeMapFile)
		if err != nil {
			log.Error("Open key mime map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.ChgmEntryPath, 0)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			items := strings.Split(line, "\t")
			if len(items) == 2 {
				key := strings.TrimSpace(items[0])
				mimeType := strings.TrimSpace(items[1])
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
		log.Error("Batch chgm error", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Error(fmt.Sprintf("Chgm '%s' => '%s' Failed, Code :%d", entry.Key, entry.MimeType, item.Code))
			} else {
				log.Debug(fmt.Sprintf("Chgm '%s' => '%s' Success, Code :%d", entry.Key, entry.MimeType, item.Code))
			}
		}
	}
}

func BatchRename(cmd string, params ...string) {
	if len(params) == 2 {
		bucket := params[0]
		oldNewKeyMapFile := params[1]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		fp, err := os.Open(oldNewKeyMapFile)
		if err != nil {
			log.Error("Open old new key map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.RenameEntryPath, 0)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			items := strings.Split(line, "\t")
			if len(items) == 2 {
				oldKey := strings.TrimSpace(items[0])
				newKey := strings.TrimSpace(items[1])
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
		log.Error("Batch rename error", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Error(fmt.Sprintf("Rename '%s' => '%s' Failed, Code :%d", entry.OldKey, entry.NewKey, item.Code))
			} else {
				log.Debug(fmt.Sprintf("Rename '%s' => '%s' Success, Code :%d", entry.OldKey, entry.NewKey, item.Code))
			}
		}
	}
}

func BatchMove(cmd string, params ...string) {
	if len(params) == 3 {
		srcBucket := params[0]
		destBucket := params[1]
		srcDestKeyMapFile := params[2]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		client := rs.New(&mac)
		fp, err := os.Open(srcDestKeyMapFile)
		if err != nil {
			log.Error("Open src dest key map file error")
			return
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]qshell.MoveEntryPath, 0)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			items := strings.Split(line, "\t")
			if len(items) == 1 || len(items) == 2 {
				srcKey := strings.TrimSpace(items[0])
				destKey := srcKey
				if len(items) == 2 {
					destKey = strings.TrimSpace(items[1])
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
		log.Error("Batch move error", err)
	}
	if len(ret) > 0 {
		for i, entry := range entries {
			item := ret[i]
			if item.Data.Error != "" {
				log.Error(fmt.Sprintf("Move '%s:%s' => '%s:%s' Failed, Code :%d",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code))
			} else {
				log.Debug(fmt.Sprintf("Move '%s:%s' => '%s:%s' Success, Code :%d",
					entry.SrcBucket, entry.SrcKey, entry.DestBucket, entry.DestKey, item.Code))
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
				log.Error("Invalid <Deadline>")
				return
			} else {
				deadline = val
			}
		} else {
			deadline = time.Now().Add(time.Second * 3600).Unix()
		}
		accountS.Get()
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

func Saveas(cmd string, params ...string) {
	if len(params) == 3 {
		publicUrl := params[0]
		saveBucket := params[1]
		saveKey := params[2]
		accountS.Get()
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
