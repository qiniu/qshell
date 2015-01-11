package cli

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	fio "github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"sync/atomic"
)

func FormPut(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		mimeType := ""
		if len(params) == 4 {
			mimeType = params[3]
		}
		accountS.Get()
		mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
		policy := rs.PutPolicy{}
		policy.Scope = bucket
		putExtra := fio.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}
		uptoken := policy.Token(&mac)
		putRet := fio.PutRet{}
		err := fio.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("Put file error", err)
		} else {
			fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
	} else {
		Help(cmd)
	}
}

func ResumablePut(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		mimeType := ""
		if len(params) == 4 {
			mimeType = params[3]
		}
		fi, err := os.Stat(localFile)
		if err != nil {
			log.Error("Stat local file error,", err)
			return
		}
		fileSize := fi.Size()
		accountS.Get()
		mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
		policy := rs.PutPolicy{}
		policy.Scope = bucket
		putExtra := rio.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}
		progressHandler := ProgressHandler{
			fileSize,
			0,
		}
		putExtra.Notify = progressHandler.Notify
		uptoken := policy.Token(&mac)
		putRet := rio.PutRet{}
		err = rio.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("\r\nPut file error", err)
		} else {
			fmt.Println("\r\nPut file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
	} else {
		Help(cmd)
	}
}

type ProgressHandler struct {
	FileSize int64
	Offset   int64
}

func (this *ProgressHandler) Percent() float32 {
	return float32(this.Offset) / float32(this.FileSize) * 100
}

func (this *ProgressHandler) Notify(blkIdx int, blkSize int, ret *rio.BlkputRet) {
	offset := ret.Offset
	atomic.AddInt64(&this.Offset, int64(offset))
	percent := this.Percent()
	output := fmt.Sprintf("Uploading %.2f%%", percent)
	for i := 0; i < len(output); i++ {
		fmt.Print("\b")
	}
	fmt.Print(output)
}
