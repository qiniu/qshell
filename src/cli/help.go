package cli

import (
	"fmt"
)

var version = "v1.0.0"

var optionDocs = map[string]string{
	"-d": "Show debug message",
}

var cmds = []string{
	"account",
	"dircache",
	"listbucket",
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
}
var cmdDocs = map[string][]string{
	"account":     []string{"qshell [-d] account [<AccessKey> <SecretKey>]", "Get/Set AccessKey and SecretKey"},
	"dircache":    []string{"qshell [-d] dircache <DirCacheRootPath> <DirCacheResultFile>", "Cache the directory structure of a file path"},
	"listbucket":  []string{"qshell [-d] listbucket <Bucket> [<Prefix>] <ListBucketResultFile>", "List all the file in the bucket by prefix"},
	"prefop":      []string{"qshell [-d] prefop <PersistentId>", "Query the fop status"},
	"fput":        []string{"qshell [-d] fput <Bucket> <Key> <LocalSmallFile>", "Form upload a small file"},
	"rput":        []string{"qshell [-d] rput <Bucket> <Key> <LocalBigFile>", "Resumable upload a big file"},
	"upload":      []string{"qshell [-d] upload <Bucket> <LocalUploadConfig>", "Batch upload files to the bucket"},
	"stat":        []string{"qshell [-d] stat <Bucket> <Key>", "Get the basic info of a remote file"},
	"delete":      []string{"qshell [-d] delete <Bucket> <Key>", "Delete a remote file in the bucket"},
	"move":        []string{"qshell [-d] move <SrcBucket> <SrcKey> <DestBucket> <DestKey>", "Move/Rename a file and save in bucket"},
	"copy":        []string{"qshell [-d] copy <SrcBucket> <SrcKey> <DestBucket> [<DestKey>]", "Make a copy of a file and save in bucket"},
	"chgm":        []string{"qshell [-d] chgm <Bucket> <Key> <NewMimeType>", "Change the mimeType of a file"},
	"fetch":       []string{"qshell [-d] fetch <RemoteResourceUrl> <Bucket> <Key>", "Fetch a remote resource by url and save in bucket"},
	"prefetch":    []string{"qshell [-d] prefetch <Bucket> <Key>", "Fetch and update the file in bucket using mirror storage"},
	"batchdelete": []string{"qshell [-d] batchdelete <Bucket> <KeyListFile>", "Batch delete files in bucket"},
	"checkqrsync": []string{"qshell [-d] checkqrsync <DirCacheResultFile> <ListBucketResultFile> <IgnoreLocalDir> [Prefix]", "Check the qrsync result"},
}

func CmdHelpList() string {
	helpAll := fmt.Sprintf("QShell %s\r\n\r\n", version)
	helpAll += "Options:\r\n"
	for k, v := range optionDocs {
		helpAll += fmt.Sprintf("  %-20s%-20s\r\n", k, v)
	}
	helpAll += "\r\n"
	helpAll += "Commands:\r\n"
	for _, cmd := range cmds {
		help := cmdDocs[cmd]
		fullCmd := help[0]
		fullDesc := help[1]
		helpAll += fmt.Sprintf("  %-100s%-80s\r\n", fullCmd, fullDesc)
	}
	return helpAll
}

func CmdHelp(cmd string) (docStr string) {
	docStr = fmt.Sprintf("Unknow cmd `%s'", cmd)
	if cmdDoc, ok := cmdDocs[cmd]; ok {
		docStr = fmt.Sprintf("%s\r\n  %s\r\n", cmdDoc[0], cmdDoc[1])
	}
	return docStr
}
