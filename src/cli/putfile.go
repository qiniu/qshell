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
	"time"
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
		startTime := time.Now()
		fStat, statErr := os.Stat(localFile)
		if statErr != nil {
			log.Error("Local file error", statErr)
			return
		}
		fsize := fStat.Size()
		err := fio.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("Put file error", err)
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
		startTime := time.Now()
		fStat, statErr := os.Stat(localFile)
		if statErr != nil {
			log.Error("Local file error", statErr)
			return
		}
		fsize := fStat.Size()
		err := rio.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("Put file error", err)
		} else {
			fmt.Println("\r\nPut file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
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
	for i, blockIndex := range this.BlockIndices {
		blockProgress := this.BlockProgresses[blockIndex]
		if int(blockProgress) != 100 {
			output += fmt.Sprintf("[Block %d=>%.2f%%]", blockIndex+1, blockProgress)
			if i < len(this.BlockIndices)-1 {
				output += ", "
			}
		}
	}
	fmt.Print(output)
	os.Stdout.Sync()
}
func (this *ProgressHandler) NotifyErr(blkIdx int, blkSize int, err error) {

}
