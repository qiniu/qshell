package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"github.com/tonycai653/iqshell/qshell"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	dirCacheCmd = &cobra.Command{
		Use:   "dircache <DirCacheRootPath> [<DirCacheResultFile>]",
		Short: "Cache the directory structure of a file path",
		Long:  "Cache the directory structure of a file path to a file, \nif <DirCacheResultFile> not specified, cache to stdout",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				fmt.Errorf("accepts between 1 and 2 arg(s), received %d\n", len(args))
			}
			if len(args) == 2 {
				if args[1] == "" {
					fmt.Errorf("DirCacheResultFile cannot be be empty\n")
				}
			}
			return nil
		},
		Run: DirCache,
	}
	lsBucketCmd = &cobra.Command{
		Use:   "listbucket <Bucket> [<ListBucketResultFile>]",
		Short: "List all the files in the bucket by prefix",
		Long:  "List all the files in the bucket by prefix to stdout if ListBucketResultFile not specified",
		Args:  cobra.RangeArgs(1, 2),
		Run:   ListBucket,
	}
	statCmd = &cobra.Command{
		Use:   "stat <Bucket> <Key>",
		Short: "Get the basic info of a remote file",
		Args:  cobra.ExactArgs(2),
		Run:   Stat,
	}
	delCmd = &cobra.Command{
		Use:   "delete <Bucket> <Key>",
		Short: "Delete a remote file in the bucket",
		Args:  cobra.ExactArgs(2),
		Run:   Delete,
	}
	moveCmd = &cobra.Command{
		Use:   "move <SrcBucket> <SrcKey> <DestBucket> [<DestKey>]",
		Short: "Move/Rename a file and save in bucket",
		Args:  cobra.RangeArgs(3, 4),
		Run:   Move,
	}
	copyCmd = &cobra.Command{
		Use:   "copy <SrcBucket> <SrcKey> <DestBucket> [<DestKey>]",
		Short: "Make a copy of a file and save in bucket",
		Args:  cobra.RangeArgs(3, 4),
		Run:   Copy,
	}
	chgmCmd = &cobra.Command{
		Use:   "chgm <Bucket> <Key> <NewMimeType>",
		Short: "Change the mime type of a file",
		Args:  cobra.ExactArgs(3),
		Run:   Chgm,
	}
	chtypeCmd = &cobra.Command{
		Use:   "chtype <Bucket> <Key> <FileType>",
		Short: "Change the file type of a file",
		Args:  cobra.ExactArgs(3),
		Run:   Chtype,
	}
	delafterCmd = &cobra.Command{
		Use:   "expire <Bucket> <Key> <DeleteAfterDays>",
		Short: "Set the deleteAfterDays of a file",
		Args:  cobra.ExactArgs(3),
		Run:   DeleteAfterDays,
	}
	fetchCmd = &cobra.Command{
		Use:   "fetch <RemoteResourceUrl> <Bucket> [<Key>]",
		Short: "Fetch a remote resource by url and save in bucket",
		Args:  cobra.RangeArgs(2, 3),
		Run:   Fetch,
	}
	mirrorCmd = &cobra.Command{
		Use:   "mirrorupdate <Bucket> <Key>",
		Short: "Fetch and update the file in bucket using mirror storage",
		Args:  cobra.ExactArgs(2),
		Run:   Prefetch,
	}
	saveAsCmd = &cobra.Command{
		Use:   "saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>",
		Short: "Create a resource access url with fop and saveas",
		Args:  cobra.ExactArgs(3),
		Run:   Saveas,
	}
	m3u8DelCmd = &cobra.Command{
		Use:   "m3u8delete <Bucket> <M3u8Key>",
		Short: "Delete m3u8 playlist and the slices it references",
		Args:  cobra.ExactArgs(2),
		Run:   M3u8Delete,
	}
	m3u8RepCmd = &cobra.Command{
		Use:   "m3u8replace <Bucket> <M3u8Key> [<NewDomain>]",
		Short: "Replace m3u8 domain in the playlist",
		Args:  cobra.RangeArgs(2, 3),
		Run:   M3u8Replace,
	}
	privateUrlCmd = &cobra.Command{
		Use:   "privateurl <PublicUrl> [<Deadline>]",
		Short: "Create private resource access url",
		Args:  cobra.RangeArgs(1, 2),
		Run:   PrivateUrl,
	}
)

var (
	listMarker string
	prefix     string
	mOverwrite bool
	cOverwrite bool
)

func init() {
	lsBucketCmd.Flags().StringVarP(&listMarker, "marker", "m", "", "list marker")
	lsBucketCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "list by prefix")

	moveCmd.Flags().BoolVarP(&mOverwrite, "overwrite", "w", false, "overwrite mode")
	copyCmd.Flags().BoolVarP(&cOverwrite, "overwrite", "w", false, "overwrite mode")

	RootCmd.AddCommand(dirCacheCmd, lsBucketCmd, statCmd, delCmd, moveCmd,
		copyCmd, chgmCmd, chtypeCmd, delafterCmd, fetchCmd, mirrorCmd,
		saveAsCmd, m3u8DelCmd, m3u8RepCmd, privateUrlCmd)
}

func DirCache(cmd *cobra.Command, params []string) {
	var cacheResultFile string
	cacheRootPath := params[0]
	if len(params) == 2 {
		cacheResultFile = params[1]
	}
	if cacheResultFile == "" {
		cacheResultFile = "stdout"
	}
	_, retErr := qshell.DirCache(cacheRootPath, cacheResultFile)
	if retErr != nil {
		os.Exit(qshell.STATUS_ERROR)
	}
}

func ListBucket(cmd *cobra.Command, params []string) {
	bucket := params[0]
	listResultFile := ""
	if len(params) == 2 {
		listResultFile = params[1]
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}

	if HostFile == "" {
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
}

func Stat(cmd *cobra.Command, params []string) {
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
}

func Delete(cmd *cobra.Command, params []string) {
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
}

func Move(cmd *cobra.Command, params []string) {
	srcBucket := params[0]
	srcKey := params[1]
	destBucket := params[2]
	destKey := srcKey
	if len(params) == 4 {
		destKey = params[3]
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

	err := client.Move(nil, srcBucket, srcKey, destBucket, destKey, mOverwrite)
	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Printf("Move error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Println("Move error,", err)
		}
		os.Exit(qshell.STATUS_ERROR)
	}
}

func Copy(cmd *cobra.Command, params []string) {
	srcBucket := params[0]
	srcKey := params[1]
	destBucket := params[2]
	destKey := srcKey
	if len(params) == 4 {
		destKey = params[3]
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

	err := client.Copy(nil, srcBucket, srcKey, destBucket, destKey, cOverwrite)
	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Printf("Copy error, %d %s, xreqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Println("Copy error,", err)
		}
		os.Exit(qshell.STATUS_ERROR)
	}
}

func Chgm(cmd *cobra.Command, params []string) {
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
}

func Chtype(cmd *cobra.Command, params []string) {
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
}

func DeleteAfterDays(cmd *cobra.Command, params []string) {
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
}

func Fetch(cmd *cobra.Command, params []string) {
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

	if HostFile == "" {
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
}

func Prefetch(cmd *cobra.Command, params []string) {
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

	if HostFile == "" {
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
}

func Saveas(cmd *cobra.Command, params []string) {
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
}

func M3u8Delete(cmd *cobra.Command, params []string) {
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
}

func M3u8Replace(cmd *cobra.Command, params []string) {
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
}

func PrivateUrl(cmd *cobra.Command, params []string) {
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
