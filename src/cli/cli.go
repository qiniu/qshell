package cli

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"qshell"
)

type CliFunc func(cmd string, params ...string)

var accountS = qshell.Account{}
var dircacheS = qshell.DirCache{}
var listbucketS = qshell.ListBucket{}
var rsfopS = qshell.RSFop{}

func printStat(bucket string, key string, entry rs.Entry) {
	statInfo := fmt.Sprintf("%-20s%-20s\r\n", "Bucket:", bucket)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Key:", key)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "Hash:", entry.Hash)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "Fsize:", entry.Fsize)
	statInfo += fmt.Sprintf("%-20s%-20d\r\n", "PutTime:", entry.PutTime)
	statInfo += fmt.Sprintf("%-20s%-20s\r\n", "MimeType:", entry.MimeType)
	fmt.Println(statInfo)
}

func Help(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(qshell.CmdHelpList())
	} else {
		for _, cmd := range cmds {
			fmt.Println(qshell.CmdHelp(cmd))
		}
	}
}

func Account(cmd string, params ...string) {
	if len(params) == 0 {
		accountS.Get()
		fmt.Println(accountS.String())
	} else if len(params) == 2 {
		accessKey := params[0]
		secretKey := params[1]
		accountS.Set(accessKey, secretKey)
	} else {
		Help(cmd)
	}
}

func DirCache(cmd string, params ...string) {
	if len(params) == 2 {
		cacheRootPath := params[0]
		cacheResultFile := params[1]
		dircacheS.Cache(cacheRootPath, cacheResultFile)
	} else {
		Help(cmd)
	}
}

func ListBucket(cmd string, params ...string) {
	if len(params) == 3 {
		bucket := params[0]
		prefix := params[1]
		listResultFile := params[2]
		//get ak,sk
		accountS.Get()
		if accountS.AccessKey != "" && accountS.SecretKey != "" {
			listbucketS.Account = accountS
			listbucketS.List(bucket, prefix, listResultFile)
		} else {
			log.Error("No AccessKey and SecretKey set error!")
		}
	} else {
		Help(cmd)
	}
}

func Prefop(cmd string, params ...string) {
	if len(params) == 1 {
		persistentId := params[0]
		accountS.Get()
		fopRet := qshell.FopRet{}
		err := rsfopS.Prefop(persistentId, &fopRet)
		if err != nil {
			log.Error(fmt.Sprintf("Can not get fop status for `%s'", persistentId), err)
		} else {
			fmt.Println(fopRet.String())
		}
	} else {
		Help(cmd)
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
		Help(cmd)
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
		Help(cmd)
	}
}
