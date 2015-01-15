package cli

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	fio "github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
	"os"
	"sort"
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
		CmdHelp(cmd)
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
		accountS.Get()
		mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
		policy := rs.PutPolicy{}
		policy.Scope = bucket
		putExtra := rio.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}
		progressHandler := ProgressHandler{
			BlockIndices:    make([]int, 0),
			BlockProgresses: make(map[int]float32),
		}
		putExtra.Notify = progressHandler.Notify
		putExtra.NotifyErr = progressHandler.NotifyErr
		uptoken := policy.Token(&mac)
		putRet := rio.PutRet{}
		err := rio.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("Put file error", err)
		} else {
			fmt.Println("\r\nPut file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
	} else {
		CmdHelp(cmd)
	}
}

type ProgressHandler struct {
	BlockIndices    []int
	BlockProgresses map[int]float32
}

func (this *ProgressHandler) Notify(blkIdx int, blkSize int, ret *rio.BlkputRet) {
	offset := ret.Offset
	perent := float32(offset) * 100 / float32(blkSize)
	if _, ok := this.BlockProgresses[blkIdx]; !ok {
		this.BlockIndices = append(this.BlockIndices, blkIdx)
		sort.Ints(this.BlockIndices)
	}
	this.BlockProgresses[blkIdx] = perent
	output := fmt.Sprintf("\r")
	for _, blockIndex := range this.BlockIndices {
		blockProgress := this.BlockProgresses[blockIndex]
		output += fmt.Sprintf("[Block %d=>%.2f%%], ", blockIndex+1, blockProgress)
	}
	fmt.Print(output)
	os.Stdout.Sync()
}
func (this *ProgressHandler) NotifyErr(blkIdx int, blkSize int, err error) {

}
