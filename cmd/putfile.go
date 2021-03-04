package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/spf13/cobra"
)

var upSettings = storage.Settings{
	Workers:   16,
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  3,
}

var (
	isResumeV2   bool
	pOverwrite   bool
	mimeType     string
	fileType     int
	workerCount  int
	rupHost      string
	fupHost      string
	callbackUrls string
	callbackHost string
)

var formPutCmd = &cobra.Command{
	Use:   "fput <Bucket> <Key> <LocalFile>",
	Short: "Form upload a local file",
	Args:  cobra.ExactArgs(3),
	Run:   FormPut,
}

var RePutCmd = &cobra.Command{
	Use:   "rput <Bucket> <Key> <LocalFile>",
	Short: "Resumable upload a local file",
	Args:  cobra.ExactArgs(3),
	Run:   ResumablePut,
}

func init() {
	formPutCmd.Flags().BoolVarP(&pOverwrite, "overwrite", "w", false, "overwrite mode")
	formPutCmd.Flags().StringVarP(&mimeType, "mimetype", "t", "", "file mime type")
	formPutCmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	formPutCmd.Flags().IntVarP(&workerCount, "worker", "c", 16, "worker count")
	formPutCmd.Flags().StringVarP(&fupHost, "up-host", "u", "", "uphost")
	formPutCmd.Flags().StringVarP(&callbackUrls, "callback-urls", "l", "", "upload callback urls, separated by comma")
	formPutCmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")

	RePutCmd.Flags().BoolVarP(&isResumeV2, "v2", "", false, "resume V2")
	RePutCmd.Flags().BoolVarP(&pOverwrite, "overwrite", "w", false, "overwrite mode")
	RePutCmd.Flags().StringVarP(&mimeType, "mimetype", "t", "", "file mime type")
	RePutCmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	RePutCmd.Flags().IntVarP(&workerCount, "worker", "c", 16, "worker count")
	RePutCmd.Flags().StringVarP(&rupHost, "up-host", "u", "", "uphost")
	RePutCmd.Flags().StringVarP(&callbackUrls, "callback-urls", "l", "", "upload callback urls, separated by comma")
	RePutCmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")

	RootCmd.AddCommand(formPutCmd, RePutCmd)
}

// 上传接口返回的文件信息
type PutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

// 【fput】使用表单上传本地文件到七牛存储空间
func FormPut(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	localFile := params[2]

	if fileType != 1 && fileType != 0 {
		fmt.Fprintln(os.Stderr, "Wrong Filetype, It should be 0 or 1")
		os.Exit(iqshell.STATUS_ERROR)
	}

	//create uptoken
	policy := storage.PutPolicy{}
	if pOverwrite {
		policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
	} else {
		policy.Scope = bucket
	}
	policy.FileType = fileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	if (callbackUrls == "" && callbackHost != "") || (callbackUrls != "" && callbackHost == "") {
		fmt.Fprintf(os.Stderr, "callbackUrls and callback must exist at the same time\n")
		os.Exit(1)
	}
	if callbackHost != "" && callbackUrls != "" {
		callbackUrls = strings.Replace(callbackUrls, ",", ";", -1)
		policy.CallbackHost = callbackHost
		policy.CallbackURL = callbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}

	var putExtra storage.PutExtra
	var upHost string

	if fupHost == "" {
		upHost = iqshell.UpHost()
	} else {
		upHost = fupHost
	}
	putExtra = storage.PutExtra{
		UpHost: upHost,
	}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}
	mac, err := iqshell.GetMac()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get Mac error: ", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()
	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Fprintf(os.Stderr, "Local file error: %v\n", statErr)
		os.Exit(iqshell.STATUS_ERROR)
	}
	fsize := fStat.Size()
	fmt.Printf("Uploading %s => %s : %s ...\n", localFile, bucket, key)

	formUploader := storage.NewFormUploader(nil)

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

	err = formUploader.PutFile(context.Background(), &putRet, uptoken, key, localFile, &putExtra)

	doneSignal <- true

	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			fmt.Fprintf(os.Stderr, "Put file error %d: %s, Reqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Fprintf(os.Stderr, "Put file error: %v\n", err)
		}
	} else {
		fmt.Print("\rProgress: 100%")
		os.Stdout.Sync()
		fmt.Println()

		fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "success!")
		fmt.Println("Hash:", putRet.Hash)
		fmt.Println("Fsize:", putRet.Fsize, "(", FormatFsize(fsize), ")")
		fmt.Println("MimeType:", putRet.MimeType)

		lastNano := time.Now().UnixNano() - startTime.UnixNano()
		lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
		avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
		fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
	}

	if err != nil {
		os.Exit(iqshell.STATUS_ERROR)
	}
}

// 使用分片上传本地文件到七牛存储空间, 一般用于较大文件的上传
// 文件会被分割成4M大小的块， 一块一块地上传文件
func ResumablePut(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	localFile := params[2]

	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println("Local file error", statErr)
		os.Exit(iqshell.STATUS_ERROR)
	}
	fsize := fStat.Size()

	upSettings.Workers = workerCount

	//create uptoken
	policy := storage.PutPolicy{}
	policy.Scope = fmt.Sprintf("%s:%s", bucket, key)

	if !pOverwrite {
		policy.InsertOnly = 1
	}
	policy.FileType = fileType
	policy.Expires = 7 * 24 * 3600
	policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
	if (callbackUrls == "" && callbackHost != "") || (callbackUrls != "" && callbackHost == "") {
		fmt.Fprintf(os.Stderr, "callbackUrls and callback must exist at the same time\n")
		os.Exit(1)
	}
	if callbackHost != "" && callbackUrls != "" {
		callbackUrls = strings.Replace(callbackUrls, ",", ";", -1)
		policy.CallbackHost = callbackHost
		policy.CallbackURL = callbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}

	var upHost string

	if rupHost == "" {
		upHost = iqshell.UpHost()
	} else {
		upHost = rupHost
	}

	mac, err := iqshell.GetMac()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get Mac error: ", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()

	fmt.Printf("Uploading %s => %s : %s ...\n", localFile, bucket, key)

	if isResumeV2 {

		resume_uploader := storage.NewResumeUploaderV2(nil)

		putExtra := storage.RputV2Extra{
			UpHost: upHost,
		}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}

		err = resume_uploader.PutFile(context.Background(), &putRet, uptoken, key, localFile, &putExtra)
	} else {

		resume_uploader := storage.NewResumeUploader(nil)

		putExtra := storage.RputExtra{
			UpHost: upHost,
		}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}

		err = resume_uploader.PutFile(context.Background(), &putRet, uptoken, key, localFile, &putExtra)
	}

	fmt.Println()
	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			fmt.Fprintf(os.Stderr, "Put file error %d: %s, Reqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Fprintf(os.Stderr, "Put file error: %v\n", err)
		}
	} else {
		fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "success!")
		fmt.Println("Hash:", putRet.Hash)
		fmt.Println("Fsize:", putRet.Fsize, "(", FormatFsize(fsize), ")")
		fmt.Println("MimeType:", putRet.MimeType)

		lastNano := time.Now().UnixNano() - startTime.UnixNano()
		lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
		avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
		fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
	}

	if err != nil {
		os.Exit(iqshell.STATUS_ERROR)
	}
}
