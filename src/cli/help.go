package cli

import (
	"fmt"
	"os"
)

var version = "v1.0.0"

var optionDocs = map[string]string{
	"-d": "Show debug message",
}

var cmds = []string{
	"account",
	"dircache",
	"listbucket",
	"alilistbucket",
	"prefop",
	"fput",
	"rput",
	"upload",
	"stat",
	"delete",
	"move",
	"copy",
	"chgm",
	"fetch",
	"prefetch",
	"batchdelete",
	"checkqrsync",
	"b64encode",
	"b64decode",
}
var cmdDocs = map[string][]string{
	"account":       []string{"qshell [-d] account [<AccessKey> <SecretKey>]", "Get/Set AccessKey and SecretKey"},
	"dircache":      []string{"qshell [-d] dircache <DirCacheRootPath> <DirCacheResultFile>", "Cache the directory structure of a file path"},
	"listbucket":    []string{"qshell [-d] listbucket <Bucket> [<Prefix>] <ListBucketResultFile>", "List all the file in the bucket by prefix"},
	"alilistbucket": []string{"qshell [-d] alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccesskeySecret> [Prefix] <ListBucketResultFile>", "List all the file in the bucket of aliyun oss by prefix"},
	"prefop":        []string{"qshell [-d] prefop <PersistentId>", "Query the fop status"},
	"fput":          []string{"qshell [-d] fput <Bucket> <Key> <LocalFile> [MimeType]", "Form upload a local file"},
	"rput":          []string{"qshell [-d] rput <Bucket> <Key> <LocalFile> [MimeType]", "Resumable upload a local file"},
	"upload":        []string{"qshell [-d] upload <Bucket> <LocalUploadConfig>", "Batch upload files to the bucket"},
	"stat":          []string{"qshell [-d] stat <Bucket> <Key>", "Get the basic info of a remote file"},
	"delete":        []string{"qshell [-d] delete <Bucket> <Key>", "Delete a remote file in the bucket"},
	"move":          []string{"qshell [-d] move <SrcBucket> <SrcKey> <DestBucket> <DestKey>", "Move/Rename a file and save in bucket"},
	"copy":          []string{"qshell [-d] copy <SrcBucket> <SrcKey> <DestBucket> [<DestKey>]", "Make a copy of a file and save in bucket"},
	"chgm":          []string{"qshell [-d] chgm <Bucket> <Key> <NewMimeType>", "Change the mimeType of a file"},
	"fetch":         []string{"qshell [-d] fetch <RemoteResourceUrl> <Bucket> <Key>", "Fetch a remote resource by url and save in bucket"},
	"prefetch":      []string{"qshell [-d] prefetch <Bucket> <Key>", "Fetch and update the file in bucket using mirror storage"},
	"batchdelete":   []string{"qshell [-d] batchdelete <Bucket> <KeyListFile>", "Batch delete files in bucket"},
	"checkqrsync":   []string{"qshell [-d] checkqrsync <DirCacheResultFile> <ListBucketResultFile> <IgnoreLocalDir> [Prefix]", "Check the qrsync result"},
	"b64encode":     []string{"qshell [-d] b64encode <UrlSafe> <DataToEncode>", "Base64 Encode"},
	"b64decode":     []string{"qshell [-d] b64decode <UrlSafe> <DataToDecode>", "Base64 Decode"},
}

func Help(cmd string, params ...string) {
	if len(params) == 0 {
		fmt.Println(CmdList())
	} else {
		CmdHelps(params...)
	}
}

func CmdList() string {
	helpAll := fmt.Sprintf("QShell %s\r\n\r\n", version)
	helpAll += "Options:\r\n"
	for k, v := range optionDocs {
		helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", k, v)
	}
	helpAll += "\r\n"
	helpAll += "Commands:\r\n"
	for _, cmd := range cmds {
		help := cmdDocs[cmd]
		cmdDesc := help[1]
		helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", cmd, cmdDesc)
	}
	return helpAll
}

func CmdHelps(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(CmdList())
	} else {
		for _, cmd := range cmds {
			fmt.Println(CmdHelp(cmd))
		}
	}
}

func CmdHelp(cmd string) (docStr string) {
	docStr = fmt.Sprintf("Unknow cmd `%s'", cmd)
	if cmdDoc, ok := cmdDocs[cmd]; ok {
		docStr = fmt.Sprintf("Usage: %s\r\n  %s\r\n", cmdDoc[0], cmdDoc[1])
	}
	return docStr
}
