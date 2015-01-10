package qshell

import (
	"fmt"
)

var version = "v1.0.0"

var optionDocs = map[string]string{
	"-d": "Show debug message",
}

var cmdDocs = map[string][]string{
	"account":    []string{"qshell [-d] account [<AccessKey> <SecretKey>]", "Get/Set AccessKey and SecretKey"},
	"dircache":   []string{"qshell [-d] dircache <CacheRootPath> <CacheResultFile>", "Cache the directory structure of a file path"},
	"listbucket": []string{"qshell [-d] listbucket <Bucket> <Prefix> <ListResultFile>", "List all the file in the bucket by prefix"},
	"prefop":     []string{"qshell [-d] prefop <PersistentId>", "Query the fop status"},
	"fput":       []string{"qshell [-d] fput <Bucket> <Key> <LocalSmallFile>", "Form upload a small file"},
	"rput":       []string{"qshell [-d] rput <Bucket> <Key> <LocalBigFile>", "Resumable upload a big file"},
	"sync":       []string{"qshell [-d] sync <Bucket> <LocalSyncConfig>", "Batch upload files to the bucket"},
	"stat":       []string{"qshell [-d] stat <Bucket> <Key>", "Get the basic info of a remote file"},
	"delete":     []string{"qshell [-d] delete <Bucket> <Key>", "Delete a remote file in the bucket"},
}

func CmdHelpList() string {
	helpAll := fmt.Sprintf("QShell %s\r\n\r\n", version)
	helpAll += "Options:\r\n"
	for k, v := range optionDocs {
		helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", k, v)
	}
	helpAll += "\r\n"
	helpAll += "Commands:\r\n"
	for _, v := range cmdDocs {
		cmd := v[0]
		desc := v[1]
		helpAll += fmt.Sprintf("\t%-60s%-100s\r\n", cmd, desc)
	}
	return helpAll
}

func CmdHelp(cmd string) (docStr string) {
	docStr = fmt.Sprintf("Unknow cmd `%s'", cmd)
	if cmdDoc, ok := cmdDocs[cmd]; ok {
		docStr = fmt.Sprintf("%-40s%-40s\r\n", cmdDoc[0], cmdDoc[1])
	}
	return docStr
}
