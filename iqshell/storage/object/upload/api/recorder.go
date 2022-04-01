package api

import (
	"encoding/json"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"time"
)

type ProgressRecorder struct {
	BlkCtxs      []storage.BlkputRet      `json:"blk_ctxs"`    // resume v1
	Parts        []storage.UploadPartInfo `json:"parts"`       // resume v2
	UploadId     string                   `json:"upload_id"`   // resume v2
	ExpireTime   int64                    `json:"expire_time"` // resume v2
	Offset       int64                    `json:"offset"`
	TotalSize    int64                    `json:"total_size"`
	LastModified int                      `json:"last_modified"` // 上传文件的modification time
	FilePath     string                   `json:"-"`             // 断点续传记录保存文件
}

func NewProgressRecorder(filePath string) *ProgressRecorder {
	p := new(ProgressRecorder)
	p.FilePath = filePath
	p.BlkCtxs = make([]storage.BlkputRet, 0)
	p.Parts = make([]storage.UploadPartInfo, 0)
	return p
}

func (p *ProgressRecorder) Recover() (err *data.CodeError) {
	if statInfo, statErr := os.Stat(p.FilePath); statErr == nil {
		//check file last modified time, if older than one week, ignore
		if statInfo.ModTime().Add(time.Hour * 24 * 5).After(time.Now()) {
			//try read old progress
			progressFh, openErr := os.Open(p.FilePath)
			if openErr != nil {
				err = data.NewEmptyError().AppendError(openErr)
				return
			}
			decoder := json.NewDecoder(progressFh)
			decoder.Decode(p)
			progressFh.Close()
		}
	}
	return
}

func (p *ProgressRecorder) Reset() {
	p.Offset = 0
	p.TotalSize = 0
	p.BlkCtxs = make([]storage.BlkputRet, 0)
	p.Parts = make([]storage.UploadPartInfo, 0)
}

func (p *ProgressRecorder) CheckValid(fileSize int64, lastModified int, isResumableV2 bool) {

	//check offset valid or not
	if p.Offset%data.BLOCK_SIZE != 0 {
		log.Info("Invalid offset from progress file,", p.Offset)
		p.Reset()
		return
	}

	// 分片 V1
	if !isResumableV2 {
		//check offset and blk ctxs
		if p.Offset != 0 && p.BlkCtxs != nil && int(p.Offset/data.BLOCK_SIZE) != len(p.BlkCtxs) {

			log.Info("Invalid offset and block info")
			p.Reset()
			return
		}

		//check blk ctxs, when no progress found
		if p.Offset == 0 || p.BlkCtxs == nil {
			p.Reset()
			return
		}

		if fileSize != p.TotalSize {
			if p.TotalSize != 0 {
				log.Warning("Remote file length changed, progress file out of date")
			}
			p.Offset = 0
			p.TotalSize = fileSize
			p.BlkCtxs = make([]storage.BlkputRet, 0)
			return
		}

		if len(p.BlkCtxs) > 0 {
			if lastModified != 0 && p.LastModified != lastModified {
				p.Reset()
			}
		}

		return
	}

	// 分片 V2
	//check offset and blk ctxs
	if p.Offset != 0 && p.Parts != nil && int(p.Offset/data.BLOCK_SIZE) != len(p.Parts) {

		log.Info("Invalid offset and block info")
		p.Reset()
		return
	}

	//check blk ctxs, when no progress found
	if p.Offset == 0 || p.Parts == nil {
		p.Reset()
		return
	}

	if fileSize != p.TotalSize {
		if p.TotalSize != 0 {
			log.Warning("Remote file length changed, progress file out of date")
		}
		p.Offset = 0
		p.TotalSize = fileSize
		p.Parts = make([]storage.UploadPartInfo, 0)
		return
	}

	if len(p.Parts) > 0 {
		if lastModified != 0 && p.LastModified != lastModified {
			p.Reset()
		}
	}
}

func (p *ProgressRecorder) RecordProgress() (err *data.CodeError) {
	fh, openErr := os.Create(p.FilePath)
	if openErr != nil {
		err = data.NewEmptyError().AppendDescF("Open progress file %s error, %s", p.FilePath, openErr.Error())
		return
	}
	defer fh.Close()

	jsonBytes, mErr := json.Marshal(p)
	if mErr != nil {
		err = data.NewEmptyError().AppendDescF("Marshal sync progress error, %s", mErr.Error())
		return
	}

	_, wErr := fh.Write(jsonBytes)
	if wErr != nil {
		err = data.NewEmptyError().AppendDescF("Write sync progress error, %s", wErr.Error())
	}

	return
}
