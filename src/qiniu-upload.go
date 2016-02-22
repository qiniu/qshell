package main

import (
	"flag"
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	rio "qiniu/api.v6/resumable/io"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
	"qshell"
	"sync"
	"time"
)

var upSettings = rio.Settings{
	Workers:   1,
	ChunkSize: 4 * 1024 * 1024,
	TryTimes:  7,
}

func main() {
	var bucket string
	var key string
	var localFile string
	var mimeType string
	var upHost string
	var overwrite bool
	var blocks int

	flag.StringVar(&bucket, "bucket", "", "bucket to save to")
	flag.StringVar(&key, "key", "", "key to save in bucket")
	flag.StringVar(&localFile, "file", "", "local file to upload")
	flag.StringVar(&mimeType, "mime", "", "mime type to set for file")
	flag.StringVar(&upHost, "host", "", "qiniu upload host")
	flag.StringVar(&upHost, "overwrite", "", "overwrite the existing file")
	flag.IntVar(&blocks, "block", 1, "block count to upload simultaneously")

	flag.Parse()

	if bucket == "" || key == "" || localFile == "" {
		fmt.Println("Err: please specify bucket, key and local file")
		return
	}

	if blocks > 0 {
		upSettings.Workers = blocks
	}

	ResumablePut(bucket, key, localFile, mimeType, upHost, overwrite)

}

func ResumablePut(bucket, key, localFile, mimeType, upHost string, overwrite bool) {
	accountS := qshell.Account{}
	gErr := accountS.Get()
	if gErr != nil {
		fmt.Println(gErr)
		return
	}

	fStat, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println("Local file error", statErr)
		return
	}
	fsize := fStat.Size()

	mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
	policy := rs.PutPolicy{}

	if overwrite {
		policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
	} else {
		policy.Scope = bucket
	}

	putExtra := rio.PutExtra{}
	if mimeType != "" {
		putExtra.MimeType = mimeType
	}

	if upHost != "" {
		conf.UP_HOST = upHost
	}

	progressHandler := ProgressHandler{
		rwLock:  &sync.RWMutex{},
		fsize:   fsize,
		offsets: make(map[int]int64, 0),
	}

	putExtra.Notify = progressHandler.Notify
	putExtra.NotifyErr = progressHandler.NotifyErr
	uptoken := policy.Token(&mac)
	putRet := rio.PutRet{}
	startTime := time.Now()

	rio.SetSettings(&upSettings)
	putClient := rio.NewClient(uptoken, "")
	fmt.Println(fmt.Sprintf("Uploading %s => %s : %s ...", localFile, bucket, key))
	err := rio.PutFile(putClient, nil, &putRet, key, localFile, &putExtra)
	fmt.Println()
	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Println(fmt.Sprintf("Put file error, %d %s, Reqid: %s", v.Code, v.Err, v.Reqid))
		} else {
			fmt.Println("Put file error,", err)
		}
	} else {
		fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
	}
	lastNano := time.Now().UnixNano() - startTime.UnixNano()
	lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
	avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
	fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
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
