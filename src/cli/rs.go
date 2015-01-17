package cli

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"qshell"
	"strings"
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
	if len(params) == 3 {
		remoteResUrl := params[0]
		bucket := params[1]
		key := params[2]
		accountS.Get()
		mac := digest.Mac{
			accountS.AccessKey,
			[]byte(accountS.SecretKey),
		}
		err := qshell.Fetch(&mac, remoteResUrl, bucket, key)
		if err != nil {
			log.Error("Fetch error,", err)
		} else {
			fmt.Println("Done!")
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
			log.Error("Open bucket key map file error", err)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		entries := make([]rs.EntryPath, 0)
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
		}
		ret, err := client.BatchDelete(nil, entries)
		if err != nil {
			log.Error("Batch delete error", err)
		} else {
			for i, entry := range entries {
				item := ret[i]
				if item.Error != "" {
					log.Debug(fmt.Sprintf("Delete %s=>%s Failed, Code: %d", entry.Bucket, entry.Key, item.Code))
				} else {
					log.Debug(fmt.Sprintf("Delete %s=>%s Success, Code: %d", entry.Bucket, entry.Key, item.Code))
				}
			}
			fmt.Println("All deleted!")
		}
	} else {
		CmdHelp(cmd)
	}
}
