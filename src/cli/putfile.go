package cli

import (
	"flag"
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	fio "qiniu/api.v6/io"
	rio "qiniu/api.v6/resumable/io"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
	"qshell"
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
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  3,
}

func FormPut(cmd string, params ...string) {
	if len(params) >= 3 && len(params) <= 7 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]

		var (
			mimeType  string
			upHost    string
			overwrite bool
			filetype  int
		)

		f := flag.NewFlagSet("fput", flag.ExitOnError)
		f.BoolVar(&overwrite, "overwrite", false, "Whether overwrite existing file")
		f.IntVar(&filetype, "filetype", 0, "Whether use cold storage")
		f.StringVar(&mimeType, "mimetype", "", "specify a mimetype for file")
		f.StringVar(&upHost, "uphost", "", "Specify a uphost")
		f.Parse(params[3:])
		if filetype != 1 && filetype != 0 {
			fmt.Println("Wrong filetype, It should be set to 0 or 1")
			os.Exit(qshell.STATUS_ERROR)
		}

		account, gErr := qshell.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(qshell.STATUS_ERROR)
		}

		//upload settings
		mac := digest.Mac{account.AccessKey, []byte(account.SecretKey)}
		if upHost == "" {
			//get bucket zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Get bucket region info error,", gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set up host
			qshell.SetZone(bucketInfo.Region)
		} else {
			conf.UP_HOST = upHost
		}

		//create uptoken
		policy := rs.PutPolicy{}
		if overwrite {
			policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
		} else {
			policy.Scope = bucket
		}
		if filetype == 1 {
			policy.FileType = 1
		}
		policy.Expires = 7 * 24 * 3600
		policy.ReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"mimeType":"$(mimeType)"}`
		putExtra := fio.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}

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
	} else {
		CmdHelp(cmd)
	}
}

func ResumablePut(cmd string, params ...string) {
	if len(params) >= 3 && len(params) <= 7 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		var (
			mimeType  string
			upHost    string
			overwrite bool
			filetype  int
		)

		f := flag.NewFlagSet("fput", flag.ExitOnError)
		f.BoolVar(&overwrite, "overwrite", false, "Whether overwrite existing file")
		f.IntVar(&filetype, "filetype", 0, "Whether use cold storage")
		f.StringVar(&mimeType, "mimetype", "", "specify a mimetype for file")
		f.StringVar(&upHost, "uphost", "", "Specify a uphost")
		f.Parse(params[3:])
		if filetype != 1 && filetype != 0 {
			fmt.Println("Wrong filetype, It should be set to 0 or 1")
			os.Exit(qshell.STATUS_ERROR)
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
			//get bucket zone info
			bucketInfo, gErr := qshell.GetBucketInfo(&mac, bucket)
			if gErr != nil {
				fmt.Println("Get bucket region info error,", gErr)
				os.Exit(qshell.STATUS_ERROR)
			}

			//set up host
			qshell.SetZone(bucketInfo.Region)
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
		if filetype == 1 {
			policy.FileType = 1
		}
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
	} else {
		CmdHelp(cmd)
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
