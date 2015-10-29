package cli

import (
	"fmt"
	"os"
	"runtime"
)

var version = "v1.5.7"

var optionDocs = map[string]string{
	"-d": "Show debug message",
	"-v": "Show version",
	"-h": "Show help",
}

var cmds = []string{
	"account",
	"zone",
	"dircache",
	"listbucket",
	"alilistbucket",
	"prefop",
	"fput",
	"rput",
	"qupload",
	"qdownload",
	"stat",
	"delete",
	"move",
	"copy",
	"chgm",
	"fetch",
	"prefetch",
	"batchstat",
	"batchdelete",
	"batchchgm",
	"batchcopy",
	"batchmove",
	"batchrename",
	"batchrefresh",
	"batchsign",
	"checkqrsync",
	"b64encode",
	"b64decode",
	"urlencode",
	"urldecode",
	"ts2d",
	"tms2d",
	"tns2d",
	"d2ts",
	"ip",
	"qetag",
	"unzip",
	"privateurl",
	"saveas",
	"reqid",
	"m3u8delete",
	"buckets",
	"domains",
}
var cmdDocs = map[string][]string{
	"account":       []string{"qshell [-d] account [<AccessKey> <SecretKey>]", "Get/Set AccessKey and SecretKey"},
	"zone":          []string{"qshell [-d] zone [Zone]", "Switch the zone, [nb,bc,aws]"},
	"dircache":      []string{"qshell [-d] dircache <DirCacheRootPath> <DirCacheResultFile>", "Cache the directory structure of a file path"},
	"listbucket":    []string{"qshell [-d] listbucket <Bucket> [<Prefix>] <ListBucketResultFile>", "List all the file in the bucket by prefix"},
	"alilistbucket": []string{"qshell [-d] alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccesskeySecret> [Prefix] <ListBucketResultFile>", "List all the file in the bucket of aliyun oss by prefix"},
	"prefop":        []string{"qshell [-d] prefop <PersistentId>", "Query the fop status"},
	"fput":          []string{"qshell [-d] fput <Bucket> <Key> <LocalFile> [Overwrite] [MimeType] [UpHost]", "Form upload a local file"},
	"rput":          []string{"qshell [-d] rput <Bucket> <Key> <LocalFile> [Overwrite] [MimeType] [UpHost]", "Resumable upload a local file"},
	"qupload":       []string{"qshell [-d] qupload [<ThreadCount>] <LocalUploadConfig>", "Batch upload files to the qiniu bucket"},
	"qdownload":     []string{"qshell [-d] qdownload [<ThreadCount>] <LocalDownloadConfig>", "Batch download files from the qiniu bucket"},
	"stat":          []string{"qshell [-d] stat <Bucket> <Key>", "Get the basic info of a remote file"},
	"delete":        []string{"qshell [-d] delete <Bucket> <Key>", "Delete a remote file in the bucket"},
	"move":          []string{"qshell [-d] move <SrcBucket> <SrcKey> <DestBucket> <DestKey>", "Move/Rename a file and save in bucket"},
	"copy":          []string{"qshell [-d] copy <SrcBucket> <SrcKey> <DestBucket> [<DestKey>]", "Make a copy of a file and save in bucket"},
	"chgm":          []string{"qshell [-d] chgm <Bucket> <Key> <NewMimeType>", "Change the mimeType of a file"},
	"fetch":         []string{"qshell [-d] fetch <RemoteResourceUrl> <Bucket> [<Key>]", "Fetch a remote resource by url and save in bucket"},
	"prefetch":      []string{"qshell [-d] prefetch <Bucket> <Key>", "Fetch and update the file in bucket using mirror storage"},
	"batchstat":     []string{"qshell [-d] batchstat <Bucket> <KeyListFile>", "Batch stat files in bucket"},
	"batchdelete":   []string{"qshell [-d] batchdelete <Bucket> <KeyListFile>", "Batch delete files in bucket"},
	"batchchgm":     []string{"qshell [-d] batchchgm <Bucket> <KeyMimeMapFile>", "Batch chgm files in bucket"},
	"batchcopy":     []string{"qshell [-d] batchcopy <SrcBucket> <DestBucket> <SrcDestKeyMapFile>", "Batch copy files from bucket to bucket"},
	"batchmove":     []string{"qshell [-d] batchmove <SrcBucket> <DestBucket> <SrcDestKeyMapFile>", "Batch move files from bucket to bucket"},
	"batchrename":   []string{"qshell [-d] batchrename <Bucket> <OldNewKeyMapFile>", "Batch rename files in the bucket"},
	"batchrefresh":  []string{"qshell [-d] batchrefresh <UrlListFile>", "Batch refresh the cdn cache by the url list file"},
	"batchsign":     []string{"qshell [-d] batchsign <UrlListFile> [<Deadline>]", "Batch create the private url from the public url list file"},
	"checkqrsync":   []string{"qshell [-d] checkqrsync <DirCacheResultFile> <ListBucketResultFile> <IgnoreLocalDir> [Prefix]", "Check the qrsync result"},
	"b64encode":     []string{"qshell [-d] b64encode [<UrlSafe>] <DataToEncode>", "Base64 Encode"},
	"b64decode":     []string{"qshell [-d] b64decode [<UrlSafe>] <DataToDecode>", "Base64 Decode"},
	"urlencode":     []string{"qshell [-d] urlencode <DataToEncode>", "Url encode"},
	"urldecode":     []string{"qshell [-d] urldecode <DataToDecode>", "Url decode"},
	"ts2d":          []string{"qshell [-d] ts2d <TimestampInSeconds>", "Convert timestamp in seconds to a date (TZ: Local)"},
	"tms2d":         []string{"qshell [-d] tms2d <TimestampInMilliSeconds>", "Convert timestamp in milli-seconds to a date (TZ: Local)"},
	"tns2d":         []string{"qshell [-d] tns2d <TimestampIn100NanoSeconds>", "Convert timestamp in 100 nano-seconds to a date (TZ: Local)"},
	"d2ts":          []string{"qshell [-d] d2ts <SecondsToNow>", "Create a timestamp in seconds using seconds to now"},
	"ip":            []string{"qshell [-d] ip <Ip1> [<Ip2> [<Ip3> ...]]]", "Query the ip information"},
	"qetag":         []string{"qshell [-d] qetag <LocalFilePath>", "Calculate the hash of local file using the algorithm of qiniu qetag"},
	"unzip":         []string{"qshell [-d] unzip <QiniuZipFilePath> [<UnzipToDir>]", "Unzip the archive file created by the qiniu mkzip API"},
	"privateurl":    []string{"qshell [-d] privateurl <PublicUrl> [<Deadline>]", "Create private resource access url"},
	"saveas":        []string{"qshell [-d] saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>", "Create a resource access url with fop and saveas"},
	"reqid":         []string{"qshell [-d] reqid <ReqIdToDecode>", "Decode a qiniu reqid"},
	"m3u8delete":    []string{"qshell [-d] m3u8delete <Bucket> <M3u8Key> [<IsPrivate>]", "Delete m3u8 playlist and the slices it references"},
	"buckets":       []string{"qshell [-d] buckets", "Get all buckets of the account"},
	"domains":       []string{"qshell [-d] domains <Bucket>", "Get all domains of the bucket"},
}

func Version() {
	fmt.Println("qshell", version)
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
		if help, ok := cmdDocs[cmd]; ok {
			cmdDesc := help[1]
			helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", cmd, cmdDesc)
		}
	}
	return helpAll
}

func CmdHelps(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(CmdList())
	} else {
		for _, cmd := range cmds {
			CmdHelp(cmd)
		}
	}
}

func CmdHelp(cmd string) {
	docStr := fmt.Sprintf("Unknow cmd `%s'", cmd)
	if cmdDoc, ok := cmdDocs[cmd]; ok {
		docStr = fmt.Sprintf("Usage: %s\r\n  %s\r\n", cmdDoc[0], cmdDoc[1])
	}
	fmt.Println(docStr)
}

func UserAgent() string {
	return fmt.Sprintf("QShell/%s (%s; %s; %s)", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
