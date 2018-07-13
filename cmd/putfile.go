package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/api.v6/auth/digest"
	"github.com/tonycai653/iqshell/qiniu/api.v6/conf"
	fio "github.com/tonycai653/iqshell/qiniu/api.v6/io"
	rio "github.com/tonycai653/iqshell/qiniu/api.v6/resumable/io"
	"github.com/tonycai653/iqshell/qiniu/api.v6/rs"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"github.com/tonycai653/iqshell/qshell"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

var upSettings = rio.Settings{
	Workers:   16,
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  3,
}

var (
	pOverwrite bool
	mimeType   string
	upHost     string
	fileType   string
)

var formPutCmd = &cobra.Command{
	Use:   "fput <Bucket> <Key> <LocalFile> [<Overwrite>] [<MimeType>] [<UpHost>] [<FileType>]",
	Short: "Form upload a local file",
	Args:  cobra.ExactArgs(3),
	Run:   FormPut,
}

var RePutCmd = &cobra.Command{
	Use:   "rput <Bucket> <Key> <LocalFile> [<Overwrite>] [<MimeType>] [<UpHost>] [<FileType>]",
	Short: "Resumable upload a local file",
	Args:  cobra.ExactArgs(3),
	Run:   ResumablePut,
}

func init() {
	formPutCmd.Flags().BoolVarP(&pOverwrite, "overwrite", "w", false, "overwrite mode")
	formPutCmd.Flags().StringVarP(&mimeType, "mimetype", "t", "", "file mime type")
	formPutCmd.Flags().StringVarP(&upHost, "uphost", "u", "", "upload host")
	formPutCmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	RootCmd.AddCommand(formPutCmd, RePutCmd)
}

func FormPut(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	localFile := params[2]

	if fileType != 1 && fileType != 0 {
		fmt.Println("Wrong Filetype, It should be 0 or 1 ")
		os.Exit(qshell.STATUS_ERROR)
	}
	if strings.HasPrefix(upHost, "http://") || strings.HasPrefix(param, "https://") {
		upHost = strings.TrimSuffix(upHost, "/")
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	//upload settings
	mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}
	if upHost == "" {
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
	} else {
		conf.UP_HOST = upHost
	}

	//create uptoken
	policy := rs.PutPolicy{}
	if pOverwrite {
		policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
	} else {
		policy.Scope = bucket
	}
	policy.FileType = fileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	putExtra := fio.PutExtra{}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}
	putExtra.CheckCrc = 1

	uptoken := policy.Token(&mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()
	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println("Local file error", statErr)
		os.Exit(qshell.STATUS_ERROR)
	}
	fsize := fStat.Size()
	putClient := rpc.NewClient("")
	fmt.Printf("Uploading %s => %s : %s ...\n", localFile, bucket, key)
	doneSignal := make(chan bool)
	go func(ch chan bool) {
		progressSigns := []string{"|", "/", "-", "\\", "|"}
		for {
			for _, p := range progressSigns {
				fmt.Print("\rProgress: ", p)
				os.Stdout.Sync()
				select {
				case <-ch:
					return
				case <-time.After(time.Millisecond * 50):
					continue
				}
			}
		}
	}(doneSignal)

	err := fio.PutFile(putClient, nil, &putRet, uptoken, key, localFile, &putExtra)
	doneSignal <- true
	fmt.Print("\rProgress: 100%")
	os.Stdout.Sync()
	fmt.Println()

	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Printf("Put file error, %d %s, Reqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Println("Put file error,", err)
		}
	} else {
		fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "success!")
		fmt.Println("Hash:", putRet.Hash)
		fmt.Println("Fsize:", putRet.Fsize, "(", FormatFsize(fsize), ")")
		fmt.Println("MimeType:", putRet.MimeType)
	}
	lastNano := time.Now().UnixNano() - startTime.UnixNano()
	lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
	avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
	fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")

	if err != nil {
		os.Exit(qshell.STATUS_ERROR)
	}
}

func ResumablePut(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	localFile := params[2]
	mimeType := ""
	upHost := ""
	overwrite := false
	fileType := 0

	optionalParams := params[3:]
	for _, param := range optionalParams {

		if ft, err := strconv.Atoi(param); err == nil {
			if ft == 1 || ft == 0 {
				fileType = ft
				continue
			} else {
				fmt.Println("Wrong Filetype, It should be 0 or 1 ")
				os.Exit(qshell.STATUS_ERROR)
			}

		}
		if val, pErr := strconv.ParseBool(param); pErr == nil {
			overwrite = val
			continue
		}
		if strings.HasPrefix(param, "http://") || strings.HasPrefix(param, "https://") {
			upHost = strings.TrimSuffix(param, "/")
			continue
		}

		mimeType = param
	}

	account, gErr := qshell.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(qshell.STATUS_ERROR)
	}

	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println("Local file error", statErr)
		os.Exit(qshell.STATUS_ERROR)
	}
	fsize := fStat.Size()

	//upload settings
	mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}
	if upHost == "" {
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
	} else {
		conf.UP_HOST = upHost
	}
	rio.SetSettings(&upSettings)

	//create uptoken
	policy := rs.PutPolicy{}
	if overwrite {
		policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
	} else {
		policy.Scope = bucket
	}
	policy.FileType = fileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`

	putExtra := rio.PutExtra{}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}

	progressHandler := ProgressHandler{
		rwLock:  &sync.RWMutex{},
		fsize:   fsize,
		offsets: make(map[int]int64, 0),
	}

	putExtra.Notify = progressHandler.Notify
	putExtra.NotifyErr = progressHandler.NotifyErr
	uptoken := policy.Token(&mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()

	putClient := rio.NewClient(uptoken, "")
	fmt.Printf("Uploading %s => %s : %s ...\n", localFile, bucket, key)
	err := rio.PutFile(putClient, nil, &putRet, key, localFile, &putExtra)
	fmt.Println()
	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Printf("Put file error, %d %s, Reqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Println("Put file error,", err)
		}
	} else {
		fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "success!")
		fmt.Println("Hash:", putRet.Hash)
		fmt.Println("Fsize:", putRet.Fsize, "(", FormatFsize(fsize), ")")
		fmt.Println("MimeType:", putRet.MimeType)
	}
	lastNano := time.Now().UnixNano() - startTime.UnixNano()
	lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
	avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
	fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")

	if err != nil {
		os.Exit(qshell.STATUS_ERROR)
	}
}

type ProgressHandler struct {
	rwLock  *sync.RWMutex
	offsets map[int]int64
	fsize   int64
}

func (this *ProgressHandler) Notify(blkIdx int, blkSize int, ret *rio.BlkputRet) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()

	this.offsets[blkIdx] = int64(ret.Offset)
	var uploaded int64
	for _, offset := range this.offsets {
		uploaded += offset
	}

	percent := fmt.Sprintf("\rProgress: %.2f%%", float64(uploaded)/float64(this.fsize)*100)
	fmt.Print(percent)
	os.Stdout.Sync()
}

func (this *ProgressHandler) NotifyErr(blkIdx int, blkSize int, err error) {

}
