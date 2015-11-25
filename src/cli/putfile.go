package cli

import (
	"fmt"
	"os"
	"qiniu/api.v6/auth/digest"
	"qiniu/api.v6/conf"
	fio "qiniu/api.v6/io"
	rio "qiniu/api.v6/resumable/io"
	"qiniu/api.v6/rs"
	"qiniu/rpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

var upSettings = rio.Settings{
	ChunkSize: 2 * 1024 * 1024,
	TryTimes:  7,
}

func FormPut(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 || len(params) == 5 || len(params) == 6 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		mimeType := ""
		upHost := ""
		overwrite := false

		optionalParams := params[3:]
		for _, param := range optionalParams {
			if val, pErr := strconv.ParseBool(param); pErr == nil {
				overwrite = val
				continue
			}

			if strings.HasPrefix(param, "http://") || strings.HasPrefix(param, "https://") {
				upHost = param
				continue
			}

			mimeType = param
		}

		gErr := accountS.Get()
		if gErr != nil {
			fmt.Println(gErr)
			return
		}

		mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
		policy := rs.PutPolicy{}
		if overwrite {
			policy.Scope = fmt.Sprintf("%s:%s", bucket, key)
		} else {
			policy.Scope = bucket
		}

		putExtra := fio.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}
		if upHost != "" {
			conf.UP_HOST = upHost
		}
		uptoken := policy.Token(&mac)
		putRet := fio.PutRet{}
		startTime := time.Now()
		fStat, statErr := os.Stat(localFile)
		if statErr != nil {
			fmt.Println("Local file error", statErr)
			return
		}
		fsize := fStat.Size()
		putClient := rpc.NewClient("")
		fmt.Println(fmt.Sprintf("Uploading %s => %s : %s ...", localFile, bucket, key))
		err := fio.PutFile(putClient, nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			fmt.Println("Put file error", err)
		} else {
			fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
		lastNano := time.Now().UnixNano() - startTime.UnixNano()
		lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
		avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
		fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
	} else {
		CmdHelp(cmd)
	}
}

func ResumablePut(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 || len(params) == 5 || len(params) == 6 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		mimeType := ""
		upHost := ""
		overwrite := false

		optionalParams := params[3:]
		for _, param := range optionalParams {
			if val, pErr := strconv.ParseBool(param); pErr == nil {
				overwrite = val
				continue
			}

			if strings.HasPrefix(param, "http://") || strings.HasPrefix(param, "https://") {
				upHost = param
				continue
			}

			mimeType = param
		}

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
			fmt.Println()
			fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
		lastNano := time.Now().UnixNano() - startTime.UnixNano()
		lastTime := fmt.Sprintf("%.2f", float32(lastNano)/1e9)
		avgSpeed := fmt.Sprintf("%.1f", float32(fsize)*1e6/float32(lastNano))
		fmt.Println("Last time:", lastTime, "s, Average Speed:", avgSpeed, "KB/s")
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
