package cmd

import (
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"os"
	"strings"
	"time"
)

var upSettings = storage.Settings{
	Workers:   16,
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  3,
}

var (
	pOverwrite  bool
	mimeType    string
	upHost      string
	fileType    int
	workerCount int
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
	formPutCmd.Flags().StringVarP(&upHost, "uphost", "u", "", "upload host")
	formPutCmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	formPutCmd.Flags().IntVarP(&workerCount, "worker", "c", 16, "worker count")

	RePutCmd.Flags().BoolVarP(&pOverwrite, "overwrite", "w", false, "overwrite mode")
	RePutCmd.Flags().StringVarP(&mimeType, "mimetype", "t", "", "file mime type")
	RePutCmd.Flags().StringVarP(&upHost, "uphost", "u", "", "upload host")
	RePutCmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	RePutCmd.Flags().IntVarP(&workerCount, "worker", "c", 16, "worker count")

	RootCmd.AddCommand(formPutCmd, RePutCmd)
}

type PutRet struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	MimeType string `json:"mimeType"`
	Fsize    int64  `json:"fsize"`
}

func FormPut(cmd *cobra.Command, params []string) {
	bucket := params[0]
	key := params[1]
	localFile := params[2]

	if fileType != 1 && fileType != 0 {
		fmt.Fprintln(os.Stderr, "Wrong Filetype, It should be 0 or 1")
		os.Exit(qshell.STATUS_ERROR)
	}
	if strings.HasPrefix(upHost, "http://") || strings.HasPrefix(upHost, "https://") {
		upHost = strings.TrimSuffix(upHost, "/")
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
	putExtra := storage.PutExtra{}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}
	mac, err := qshell.GetMac()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get Mac error: ", err)
		os.Exit(qshell.STATUS_ERROR)
	}
	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()
	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Fprintf(os.Stderr, "Local file error: %v\n", statErr)
		os.Exit(qshell.STATUS_ERROR)
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

	err = formUploader.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)

	doneSignal <- true
	fmt.Print("\rProgress: 100%")
	os.Stdout.Sync()
	fmt.Println()

	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
			fmt.Fprintf(os.Stderr, "Put file error, %d %s, Reqid: %s\n", v.Code, v.Err, v.Reqid)
		} else {
			fmt.Fprintln(os.Stderr, "Put file error: %v\n", err)
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

	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println("Local file error", statErr)
		os.Exit(qshell.STATUS_ERROR)
	}
	fsize := fStat.Size()

	upSettings.Workers = workerCount

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

	putExtra := storage.RputExtra{}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}

	mac, err := qshell.GetMac()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get Mac error: ", err)
		os.Exit(qshell.STATUS_ERROR)
	}
	uptoken := policy.UploadToken(mac)

	//start to upload
	putRet := PutRet{}
	startTime := time.Now()

	fmt.Printf("Uploading %s => %s : %s ...\n", localFile, bucket, key)

	resume_uploader := storage.NewResumeUploader(nil)
	err = resume_uploader.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
	fmt.Println()
	if err != nil {
		if v, ok := err.(*storage.ErrorInfo); ok {
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
